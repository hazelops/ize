package destroy

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/apps"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DestroyOptions struct {
	Config           *config.Config
	AppName          string
	Tag              string
	SkipBuildAndPush bool
	Apps             Apps
	Infra            Infra
	App              apps.App
	AutoApprove      bool
}

type Apps map[string]*apps.App

type Infra struct {
	Version string `mapstructure:"terraform_version"`
	Region  string `mapstructure:"aws_region"`
	Profile string `mapstructure:"aws_profile"`
}

var destoyLongDesc = templates.LongDesc(`
	Destroy infraftructure or app.
	For destoy app the app name must be specimfied. 
`)

var destoyExample = templates.Examples(`
	# Destroy all (config file required)
	ize destroy

	# Destroy app (config file required)
	ize destroy <app name>

	# Destroy app via config file
	ize --config-file (or -c) /path/to/config destroy <app name>

	# Destroy app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize destroy <app name>
`)

func NewDestroyFlags() *DestroyOptions {
	return &DestroyOptions{}
}

func NewCmdDestroy(ui terminal.UI) *cobra.Command {
	o := NewDestroyFlags()

	cmd := &cobra.Command{
		Use:     "destroy [flags] [app name]",
		Example: destoyExample,
		Short:   "destoy anything",
		Long:    destoyLongDesc,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			if len(args) == 0 && !o.AutoApprove {
				pterm.Warning.Println("please set flag --auto-approve")
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

			err = o.Run(ui)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&o.AutoApprove, "auto-approve", false, "approve deploy all")

	cmd.AddCommand(
		NewCmdDestroyInfra(ui),
	)

	return cmd
}

func (o *DestroyOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		o.Config, err = config.InitializeConfig(config.WithConfigFile())
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
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
		o.Config, err = config.InitializeConfig(config.WithDocker(), config.WithConfigFile())
		viper.BindPFlags(cmd.Flags())
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
		}
		o.AppName = cmd.Flags().Args()[0]
		viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)
	}

	o.Tag = viper.GetString("tag")

	return nil
}

func (o *DestroyOptions) Validate() error {
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

func (o *DestroyOptions) Run(ui terminal.UI) error {
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

func validate(o *DestroyOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified\n")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified\n")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate options: app name be specified\n")
	}

	return nil
}

func validateAll(o *DestroyOptions) error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified\n")
	}

	if len(o.Tag) == 0 {
		return fmt.Errorf("can't validate options: tag must be specified\n")
	}

	for sname, svc := range o.Apps {
		if len(svc.Type) == 0 {
			return fmt.Errorf("can't validate options: type for app %s must be specified\n", sname)
		}
	}

	return nil
}

func destroyAll(ui terminal.UI, o *DestroyOptions) error {

	logrus.Debug(o.Apps)

	ui.Output("Destroying apps...", terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	err := apps.InReversDependencyOrder(aws.BackgroundContext(), o.Apps, func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Infra.Profile

		o.Apps[name].Name = name
		err := o.Apps[name].Destroy(sg, ui)
		if err != nil {
			return fmt.Errorf("can't destroy all: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	ui.Output("Running destroy infra...", terminal.WithHeaderStyle())

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

	//terraform destroy run options
	opts := terraform.Options{
		ContainerName:    "terraform",
		Cmd:              []string{"destroy", "-auto-approve"},
		Env:              env,
		TerraformVersion: o.Infra.Version,
	}

	ui.Output("Execution terraform destroy...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return fmt.Errorf("can't destroy all: %w", err)
	}

	ui.Output("Destroy all completed!\n", terminal.WithSuccessStyle())

	return nil
}

func destroyApp(ui terminal.UI, o *DestroyOptions) error {
	ui.Output("Destroying %s app...", o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	o.App.Name = o.AppName
	err := o.App.Destroy(sg, ui)
	if err != nil {
		return fmt.Errorf("can't destroy: %w", err)
	}

	ui.Output("destroy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}
