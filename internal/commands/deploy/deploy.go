package deploy

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/apps"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeployOptions struct {
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

var deployLongDesc = templates.LongDesc(`
	Deploy infrastructure or service.
    App name must be specified for a app deploy. 
	The infrastructure for the app must be prepared in advance.
`)

var deployExample = templates.Examples(`
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

func NewDeployFlags() *DeployOptions {
	return &DeployOptions{}
}

func NewCmdDeploy() *cobra.Command {
	o := NewDeployFlags()

	cmd := &cobra.Command{
		Use:     "deploy [flags] <app name>",
		Example: deployExample,
		Short:   "Manage deployments",
		Long:    deployLongDesc,
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
		NewCmdDeployInfra(),
	)

	return cmd
}

func (o *DeployOptions) Complete(cmd *cobra.Command, args []string) error {
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

func (o *DeployOptions) Validate() error {
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

func (o *DeployOptions) Run() error {
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

func validate(o *DeployOptions) error {
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

func validateAll(o *DeployOptions) error {
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

func deployAll(ui terminal.UI, o *DeployOptions) error {
	var tf terraform.Terraform

	logrus.Infof("infra: %s", o.Infra)

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", o.Infra.Profile),
		fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
		fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	if o.Config.IsDockerRuntime {
		tf = terraform.NewDockerTerraform(o.Infra.Version, []string{"init", "-input=true"}, env, nil)
	} else {
		tf = terraform.NewLocalTerraform(o.Infra.Version, []string{"init", "-input=true"}, env, nil)
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

	name := fmt.Sprintf("/%s/terraform-output", o.Config.Env)

	byteValue, _ := ioutil.ReadAll(&output)
	sDec := base64.StdEncoding.EncodeToString(byteValue)
	if err != nil {
		return err
	}

	ssm.New(o.Config.Session).PutParameter(&ssm.PutParameterInput{
		Name:      &name,
		Value:     aws.String(string(sDec)),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String("Intelligent-Tiering"),
		DataType:  aws.String("text"),
	})

	logrus.Debug(o.Apps)

	ui.Output("Deploying apps...", terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	err = apps.InDependencyOrder(aws.BackgroundContext(), o.Apps, func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Infra.Profile

		at := (*o.Apps[name]).(map[string]interface{})["type"].(string)

		var deployment apps.Deployment

		switch at {
		case "ecs":
			deployment = apps.NewECSDeployment(name, *o.Apps[name])
		case "serverless":
			deployment = apps.NewServerlessDeployment(name, *o.Apps[name])
		case "alias":
			deployment = apps.NewAliasDeployment(name)
		default:
			return fmt.Errorf("apps type of %s not supported", at)
		}

		err := deployment.Deploy(sg, ui)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	ui.Output("Deploy all completed!\n", terminal.WithSuccessStyle())

	return nil
}

func deployApp(ui terminal.UI, o *DeployOptions) error {
	ui.Output("Deploying %s app...", o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	var appType string

	app, ok := o.App.(map[string]interface{})
	if !ok {
		appType = "ecs"
	} else {
		appType, ok = app["type"].(string)
		if !ok {
			appType = "ecs"
		}
	}

	var deployment apps.Deployment

	switch appType {
	case "ecs":
		deployment = apps.NewECSDeployment(o.AppName, o.App)
	case "serverless":
		deployment = apps.NewServerlessDeployment(o.AppName, o.App)
	case "alias":
		deployment = apps.NewAliasDeployment(o.AppName)
	default:
		return fmt.Errorf("apps type of %s not supported", appType)
	}

	err := deployment.Deploy(sg, ui)
	if err != nil {
		return err
	}

	ui.Output("Deploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}
