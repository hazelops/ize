package up

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/apps"
	"github.com/hazelops/ize/internal/commands/gen"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type UpOptions struct {
	Config           *config.Config
	AppName          string
	Tag              string
	SkipBuildAndPush bool
	Apps             map[string]*interface{}
	Infra            Infra
	App              interface{}
	AutoApprove      bool
	UI               terminal.UI
}

type Apps map[string]*interface{}

type Infra struct {
	Version string `mapstructure:"terraform_version"`
	Region  string `mapstructure:"aws_region"`
	Profile string `mapstructure:"aws_profile"`
}

var upLongDesc = templates.LongDesc(`
	Deploy infrastructure or service.
    App name must be specified for a app deploy. 
	The infrastructure for the app must be prepared in advance.
`)

var upExample = templates.Examples(`
	# Deploy all (config file required)
	ize deploy

	# Deploy app (config file required)
	ize deploy <app name>

	# Deploy app via config file
	ize --config-file (or -c) /path/to/config deploy <app name>

	# Deploy app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize deploy <app name>
`)

func NewUpFlags() *UpOptions {
	return &UpOptions{}
}

func NewCmdUp() *cobra.Command {
	o := NewUpFlags()

	cmd := &cobra.Command{
		Use:     "up [flags] <app name>",
		Example: upExample,
		Short:   "Manage deployments",
		Long:    upLongDesc,
		Args:    cobra.MaximumNArgs(1),
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

	cmd.AddCommand(
		NewCmdUpInfra(),
	)

	return cmd
}

func (o *UpOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		if err := config.CheckRequirements(config.WithConfigFile()); err != nil {
			return err
		}
		o.Config, err = config.GetConfig()
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can't deploy your stack: %w", err)
		}

		viper.UnmarshalKey("app", &o.Apps)
		viper.UnmarshalKey("infra.terraform", &o.Infra)

		if len(o.Infra.Profile) == 0 {
			o.Infra.Profile = o.Config.AwsProfile
		}

		if len(o.Infra.Region) == 0 {
			o.Infra.Region = o.Config.AwsRegion
		}

		if len(o.Infra.Version) == 0 {
			o.Infra.Version = viper.GetString("terraform_version")
		}
	} else {
		o.Config, err = config.GetConfig()
		if err != nil {
			return fmt.Errorf("can't deploy your stack: %w", err)
		}

		viper.BindPFlags(cmd.Flags())
		o.AppName = cmd.Flags().Args()[0]
		viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)
	}

	o.Tag = viper.GetString("tag")
	o.UI = terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)

	return nil
}

func (o *UpOptions) Validate() error {
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

func (o *UpOptions) Run() error {
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

func validate(o *UpOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate options: app name must be specified")
	}

	return nil
}

func validateAll(o *UpOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified")
	}

	return nil
}

func deployAll(ui terminal.UI, o *UpOptions) error {
	err := deployInfra(ui, o.Infra, *o.Config)
	if err != nil {
		return err
	}

	logrus.Debug(o.Apps)

	ui.Output("Deploying apps...", terminal.WithHeaderStyle())

	err = apps.InDependencyOrder(aws.BackgroundContext(), o.Apps, func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Infra.Profile

		at := (*o.Apps[name]).(map[string]interface{})["type"].(string)

		var app apps.App

		switch at {
		case "ecs":
			app = apps.NewECSApp(name, *o.Apps[name])
		case "serverless":
			app = apps.NewServerlessApp(name, *o.Apps[name])
		case "alias":
			app = apps.NewAliasApp(name)
		default:
			return fmt.Errorf("apps type of %s not supported", at)
		}

		// build app container
		err := app.Build(ui)
		if err != nil {
			return fmt.Errorf("can't build app: %w", err)
		}

		// push app image
		err = app.Push(ui)
		if err != nil {
			return fmt.Errorf("can't push app: %w", err)
		}

		// deploy app image
		err = app.Deploy(ui)
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

func deployApp(ui terminal.UI, o *UpOptions) error {
	ui.Output("Deploying %s app...", o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	var appType string

	a, ok := o.App.(map[string]interface{})
	if !ok {
		appType = "ecs"
	} else {
		appType, ok = a["type"].(string)
		if !ok {
			appType = "ecs"
		}
	}

	var app apps.App

	switch appType {
	case "ecs":
		app = apps.NewECSApp(o.AppName, o.App)
	case "serverless":
		app = apps.NewServerlessApp(o.AppName, o.App)
	case "alias":
		app = apps.NewAliasApp(o.AppName)
	default:
		return fmt.Errorf("apps type of %s not supported", appType)
	}

	// build app container
	err := app.Build(ui)
	if err != nil {
		return fmt.Errorf("can't build app: %w", err)
	}

	// push app image
	err = app.Push(ui)
	if err != nil {
		return fmt.Errorf("can't push app: %w", err)
	}

	// deploy app image
	err = app.Deploy(ui)
	if err != nil {
		return fmt.Errorf("can't deploy app: %w", err)
	}

	ui.Output("Deploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}

func deployInfra(ui terminal.UI, infra Infra, config config.Config) error {
	var tf terraform.Terraform

	logrus.Infof("infra: %s", infra)

	v, err := config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	if !checkFileExists(filepath.Join(viper.GetString("ENV_DIR"), "backend.tf")) || !checkFileExists(filepath.Join(viper.GetString("ENV_DIR"), "terraform.tfvars")) {
		gen.NewCmdEnv().Execute()
	}

	env := []string{
		fmt.Sprintf("ENV=%v", config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", infra.Profile),
		fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
		fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	if config.IsDockerRuntime {
		tf = terraform.NewDockerTerraform(infra.Version, []string{"init", "-input=true"}, env, nil)
	} else {
		tf = terraform.NewLocalTerraform(infra.Version, []string{"init", "-input=true"}, env, nil)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't deploy all: %w", err)
		}
	}

	ui.Output(fmt.Sprintf("[%s] Running deploy infra...", viper.Get("ENV")), terminal.WithHeaderStyle())
	ui.Output("Execution terraform init...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	ui.Output("Execution terraform plan...", terminal.WithHeaderStyle())

	outPath := fmt.Sprintf("%s/.terraform/tfplan", viper.GetString("ENV_DIR"))

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
		Value:     aws.String(string(sDec)),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String("Intelligent-Tiering"),
		DataType:  aws.String("text"),
	})

	ui.Output("Deploy infra completed!\n", terminal.WithSuccessStyle())

	return nil
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}