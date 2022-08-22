package up

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/commands/gen"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type Options struct {
	Config           *config.Project
	AppName          string
	SkipBuildAndPush bool
	SkipGen          bool
	AutoApprove      bool
	UI               terminal.UI
}

type Apps map[string]*interface{}

var upLongDesc = templates.LongDesc(`
	Deploy infrastructure or service.
    App name must be specified for a bringing it up.  
	The infrastructure for the app must be ready to 
	receive the deployment (generally created via ize infra up in CI/CD).
`)

var upExample = templates.Examples(`
	# Deploy all (config file required)
	ize up

	# Deploy app (config file required)
	ize up <app name>

	# Deploy app with explicitly specified config file
	ize --config-file (or -c) /path/to/config up <app name>

	# Deploy app with explicitly specified config file passed via environment variable
	export IZE_CONFIG_FILE=/path/to/config
	ize up <app name>
`)

func NewUpFlags() *Options {
	return &Options{}
}

func NewCmdUp() *cobra.Command {
	o := NewUpFlags()

	cmd := &cobra.Command{
		Use:               "up [flags] <app name>",
		Example:           upExample,
		Short:             "Bring full application up from the bottom to the top.",
		Long:              upLongDesc,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if len(args) == 0 && !o.AutoApprove {
				pterm.Warning.Println("Please set flag --auto-approve")
				return nil
			}

			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.AutoApprove, "auto-approve", false, "approve deploy all")
	cmd.Flags().BoolVar(&o.SkipGen, "skip-gen", false, "skip generating terraform files")

	cmd.AddCommand(
		NewCmdUpInfra(),
	)

	return cmd
}

func (o *Options) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		o.Config, err = config.GetConfig()
		if err != nil {
			return fmt.Errorf("can't deploy your stack: %w", err)
		}

		if err := config.CheckRequirements(config.WithIzeStructure(), config.WithConfigFile()); err != nil {
			return err
		}

		if o.Config.Serverless != nil {
			if err = config.CheckRequirements(config.WithNVM()); err != nil {
				return err
			}
		}

		if o.Config.Terraform == nil {
			o.Config.Terraform = map[string]*config.Terraform{}
			o.Config.Terraform["infra"] = &config.Terraform{}
		}

		if len(o.Config.Terraform["infra"].AwsProfile) == 0 {
			o.Config.Terraform["infra"].AwsProfile = o.Config.AwsProfile
		}

		if len(o.Config.Terraform["infra"].AwsRegion) == 0 {
			o.Config.Terraform["infra"].AwsProfile = o.Config.AwsRegion
		}

		if len(o.Config.Terraform["infra"].Version) == 0 {
			o.Config.Terraform["infra"].Version = o.Config.TerraformVersion
		}
	} else {
		o.Config, err = config.GetConfig()
		if err != nil {
			return fmt.Errorf("can't deploy your stack: %w", err)
		}

		if err := config.CheckRequirements(config.WithIzeStructure(), config.WithConfigFile()); err != nil {
			return err
		}

		if o.Config.Serverless != nil {
			if err = config.CheckRequirements(config.WithNVM()); err != nil {
				return err
			}
		}

		o.AppName = cmd.Flags().Args()[0]
	}

	o.UI = terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	return nil
}

func (o *Options) Validate() error {
	if o.AppName == "" {
		err := validateAll(o)
		if err != nil {
			return err
		}
	} else {
		err := validate(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *Options) Run() error {
	ui := o.UI
	if o.AppName == "" {
		err := deployAll(ui, o)
		if err != nil {
			return err
		}
	} else {
		err := deployApp(ui, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func validate(o *Options) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate options: app name must be specified")
	}

	return nil
}

func validateAll(o *Options) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	return nil
}

func deployAll(ui terminal.UI, o *Options) error {
	err := deployInfra(ui, o.Config, o.SkipGen)
	if err != nil {
		return err
	}

	ui.Output("Deploying apps...", terminal.WithHeaderStyle())

	err = manager.InDependencyOrder(aws.BackgroundContext(), o.Config.GetApps(), func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Config.Terraform["infra"].AwsProfile

		var manager manager.Manager

		if app, ok := o.Config.Serverless[name]; ok {
			app.Name = name
			manager = &serverless.Manager{
				Project: o.Config,
				App:     app,
			}
		}
		if app, ok := o.Config.Alias[name]; ok {
			app.Name = name
			manager = &alias.Manager{
				Project: o.Config,
				App:     app,
			}
		}
		if app, ok := o.Config.Ecs[name]; ok {
			app.Name = name
			manager = &ecs.Manager{
				Project: o.Config,
				App:     app,
			}
		}

		// build app container
		err := manager.Build(ui)
		if err != nil {
			return fmt.Errorf("can't build app: %w", err)
		}

		// push app image
		err = manager.Push(ui)
		if err != nil {
			return fmt.Errorf("can't push app: %w", err)
		}

		// deploy app image
		err = manager.Deploy(ui)
		if err != nil {
			return fmt.Errorf("can't deploy app: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	ui.Output("Deploy all completed!\n", terminal.WithSuccessStyle())

	return nil
}

func deployApp(ui terminal.UI, o *Options) error {
	var manager manager.Manager
	var icon string

	manager = &ecs.Manager{
		Project: o.Config,
		App:     &config.Ecs{Name: o.AppName},
	}

	if app, ok := o.Config.Serverless[o.AppName]; ok {
		app.Name = o.AppName
		manager = &serverless.Manager{
			Project: o.Config,
			App:     app,
		}
	}
	if app, ok := o.Config.Alias[o.AppName]; ok {
		app.Name = o.AppName
		manager = &alias.Manager{
			Project: o.Config,
			App:     app,
		}
	}
	if app, ok := o.Config.Ecs[o.AppName]; ok {
		app.Name = o.AppName
		manager = &ecs.Manager{
			Project: o.Config,
			App:     app,
		}
	}

	if len(icon) != 0 {
	}

	ui.Output("Deploying %s%s app...", icon, o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	// build app container
	err := manager.Build(ui)
	if err != nil {
		return fmt.Errorf("can't build app: %w", err)
	}

	// push app image
	err = manager.Push(ui)
	if err != nil {
		return fmt.Errorf("can't push app: %w", err)
	}

	// deploy app image
	err = manager.Deploy(ui)
	if err != nil {
		return fmt.Errorf("can't deploy app: %w", err)
	}

	ui.Output("Deploy app %s%s completed\n", icon, o.AppName, terminal.WithSuccessStyle())

	return nil
}

func deployInfra(ui terminal.UI, config *config.Project, skipGen bool) error {
	if !skipGen {
		err := gen.GenerateTerraformFiles(
			config,
			"",
		)
		if err != nil {
			return err
		}
	}

	var tf terraform.Terraform

	logrus.Infof("infra: %s", config.Terraform["infra"])

	v, err := config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't get AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", config.Terraform["infra"].AwsProfile),
		fmt.Sprintf("TF_LOG=%v", config.TFLog),
		fmt.Sprintf("TF_LOG_PATH=%v", config.TFLogPath),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	switch config.PreferRuntime {
	case "docker":
		tf = terraform.NewDockerTerraform(config.Terraform["infra"].Version, []string{"init", "-input=true"}, env, nil, config)
	case "native":
		tf = terraform.NewLocalTerraform(config.Terraform["infra"].Version, []string{"init", "-input=true"}, env, nil, config)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't deploy all: %w", err)
		}
	default:
		return fmt.Errorf("can't supported %s runtime", config.PreferRuntime)
	}

	ui.Output(fmt.Sprintf("[%s] Running deploy infra...", config.Env), terminal.WithHeaderStyle())
	ui.Output("Execution terraform init...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	ui.Output("Execution terraform plan...", terminal.WithHeaderStyle())

	outPath := fmt.Sprintf("%s/.terraform/tfplan", config.EnvDir)

	//terraform plan run options
	tf.NewCmd([]string{"plan", fmt.Sprintf("-out=%s", outPath)})

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	//terraform apply run options
	tf.NewCmd([]string{"apply", "-auto-approve", outPath})

	ui.Output("Execution terraform apply...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	//terraform output run options

	tf.NewCmd([]string{"output", "-json"})

	var output bytes.Buffer

	tf.SetOut(&output)

	ui.Output("Execution terraform output...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("/%s/terraform-output", config.Env)

	byteValue, _ := ioutil.ReadAll(&output)
	sDec := base64.StdEncoding.EncodeToString(byteValue)
	if err != nil {
		return err
	}

	ssm.New(config.Session).PutParameter(&ssm.PutParameterInput{
		Name:      &name,
		Value:     aws.String(sDec),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String("Intelligent-Tiering"),
		DataType:  aws.String("text"),
	})

	ui.Output("Deploy infra completed!\n", terminal.WithSuccessStyle())

	return nil
}
