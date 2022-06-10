package down

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/apps"
	"github.com/hazelops/ize/internal/apps/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DownOptions struct {
	Config           *config.Config
	AppName          string
	Tag              string
	SkipBuildAndPush bool
	Apps             Apps
	Infra            Infra
	App              interface{}
	AutoApprove      bool
	ui               terminal.UI
}

type Apps map[string]*interface{}

type Infra struct {
	Version string `mapstructure:"terraform_version"`
	Region  string `mapstructure:"aws_region"`
	Profile string `mapstructure:"aws_profile"`
}

var downLongDesc = templates.LongDesc(`
	Destroy infrastructure or application.
	For app destroy the app name must be specified.
`)

var downExample = templates.Examples(`
	# Destroy all (config file required)
	ize down

	# Destroy app (config file required)
	ize down <app name>

	# Destroy app with explicitly specified config file
	ize --config-file (or -c) /path/to/config down <app name>

	# Destroy app with explicitly specified config file passed via environment variable.
	export IZE_CONFIG_FILE=/path/to/config
	ize down <app name>
`)

func NewDownFlags() *DownOptions {
	return &DownOptions{}
}

func NewCmdDown() *cobra.Command {
	o := NewDownFlags()

	cmd := &cobra.Command{
		Use:     "down [flags] [app name]",
		Example: downExample,
		Short:   "Destroy application",
		Long:    downLongDesc,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if len(args) == 0 && !o.AutoApprove {
				pterm.Warning.Println("Please set the flag: --auto-approve")
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
		NewCmdDownInfra(),
	)

	return cmd
}

func (o *DownOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		if err := config.CheckRequirements(config.WithConfigFile()); err != nil {
			return err
		}
		o.Config, err = config.GetConfig()
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can't load options for a command: %w", err)
		}

		viper.UnmarshalKey("app", &o.Apps)
		viper.UnmarshalKey("infra.terraform", &o.Infra)

		for _, v := range o.Apps {
			fmt.Println(*v)
		}

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
			return fmt.Errorf("can't load options for a command: %w", err)
		}

		viper.BindPFlags(cmd.Flags())
		o.AppName = cmd.Flags().Args()[0]
		viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)
	}

	o.Tag = viper.GetString("tag")
	o.ui = terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)

	return nil
}

func (o *DownOptions) Validate() error {
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

func (o *DownOptions) Run() error {
	ui := o.ui
	if o.AppName == "" {
		err := destroyAll(ui, o)
		if err != nil {
			return err
		}
	} else {
		err := destroyApp(ui, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func validate(o *DownOptions) error {
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
		return fmt.Errorf("can't validate options: app name be specified")
	}

	return nil
}

func validateAll(o *DownOptions) error {
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

func destroyAll(ui terminal.UI, o *DownOptions) error {

	logrus.Debug(o.Apps)

	ui.Output("Destroying apps...", terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	err := apps.InReversDependencyOrder(aws.BackgroundContext(), o.Apps, func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Infra.Profile

		at := (*o.Apps[name]).(map[string]interface{})["type"].(string)

		var deployment apps.App

		switch at {
		case "ecs":
			deployment = ecs.NewECSApp(name, *o.Apps[name])
		case "serverless":
			deployment = apps.NewServerlessApp(name, *o.Apps[name])
		case "alias":
			deployment = apps.NewAliasApp(name)
		default:
			return fmt.Errorf("%s apps are not supported in this command", at)
		}

		err := deployment.Destroy(ui)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

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
		tf = terraform.NewDockerTerraform(o.Infra.Version, []string{"destroy", "-auto-approve"}, env, nil)
	} else {
		tf = terraform.NewLocalTerraform(o.Infra.Version, []string{"destroy", "-auto-approve"}, env, nil)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't destroy all: %w", err)
		}
	}

	ui.Output("Running destroy infra...", terminal.WithHeaderStyle())
	ui.Output("Execution terraform destroy...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't destroy all: %w", err)
	}

	ui.Output("Destroy all completed!\n", terminal.WithSuccessStyle())

	return nil
}

func destroyApp(ui terminal.UI, o *DownOptions) error {
	ui.Output("Destroying %s app...", o.AppName, terminal.WithHeaderStyle())
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

	var deployment apps.App

	switch appType {
	case "ecs":
		deployment = ecs.NewECSApp(o.AppName, o.App)
	case "serverless":
		deployment = apps.NewServerlessApp(o.AppName, o.App)
	case "alias":
		deployment = apps.NewAliasApp(o.AppName)
	default:
		return fmt.Errorf("%s apps are not supported in this command", appType)
	}

	err := deployment.Destroy(ui)
	if err != nil {
		return fmt.Errorf("can't down: %w", err)
	}

	ui.Output("Destroy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}
