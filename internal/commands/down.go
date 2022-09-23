package commands

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type DownOptions struct {
	Config           *config.Project
	AppName          string
	SkipBuildAndPush bool
	AutoApprove      bool
	SkipGen          bool
	ui               terminal.UI
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

func NewDownFlags(project *config.Project) *DownOptions {
	return &DownOptions{
		Config: project,
	}
}

func NewCmdDown(project *config.Project) *cobra.Command {
	o := NewDownFlags(project)

	cmd := &cobra.Command{
		Use:               "down [flags] [app name]",
		Example:           downExample,
		Short:             "Destroy application",
		Long:              downLongDesc,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: config.GetApps,
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
	cmd.Flags().BoolVar(&o.SkipGen, "skip-gen", false, "skip generating terraform files")

	cmd.AddCommand(
		NewCmdDownInfra(project),
	)

	return cmd
}

func (o *DownOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	if len(args) == 0 {
		if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile()); err != nil {
			return err
		}

		if len(o.Config.Serverless) != 0 {
			if err = requirements.CheckRequirements(requirements.WithNVM()); err != nil {
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
			o.Config.Terraform["infra"].AwsRegion = o.Config.AwsRegion
		}

		if len(o.Config.Terraform["infra"].Version) == 0 {
			o.Config.Terraform["infra"].Version = o.Config.TerraformVersion
		}
	} else {
		if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile()); err != nil {
			return err
		}

		if len(o.Config.Serverless) != 0 {
			if err = requirements.CheckRequirements(requirements.WithNVM()); err != nil {
				return err
			}
		}

		o.AppName = cmd.Flags().Args()[0]
	}

	o.ui = terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	return nil
}

func (o *DownOptions) Validate() error {
	if o.AppName == "" {
		err := o.validateAll()
		if err != nil {
			return err
		}
	} else {
		err := o.validate()
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

func (o *DownOptions) validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate options: app name be specified")
	}

	return nil
}

func (o *DownOptions) validateAll() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	return nil
}

func destroyAll(ui terminal.UI, o *DownOptions) error {

	ui.Output("Destroying apps...", terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	err := manager.InReversDependencyOrder(aws.BackgroundContext(), o.Config.GetApps(), func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Config.Terraform["infra"].AwsProfile

		var m manager.Manager

		if app, ok := o.Config.Serverless[name]; ok {
			app.Name = name
			m = &serverless.Manager{
				Project: o.Config,
				App:     app,
			}
		}
		if app, ok := o.Config.Alias[name]; ok {
			app.Name = name
			m = &alias.Manager{
				Project: o.Config,
				App:     app,
			}
		}
		if app, ok := o.Config.Ecs[name]; ok {
			app.Name = name
			m = &ecs.Manager{
				Project: o.Config,
				App:     app,
			}
		}

		// destroy
		err := m.Destroy(ui)
		if err != nil {
			return fmt.Errorf("can't destroy app: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = destroyInfra(ui, o.Config, o.SkipGen)
	if err != nil {
		return err
	}

	ui.Output("Destroy all completed!\n", terminal.WithSuccessStyle())

	return nil
}

func destroyInfra(ui terminal.UI, config *config.Project, skipGen bool) error {
	if !skipGen {
		err := GenerateTerraformFiles(
			config,
			"",
		)
		if err != nil {
			return err
		}
	}

	var tf terraform.Terraform

	logrus.Infof("infra: %s", tf)

	v, err := config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
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
		tf = terraform.NewDockerTerraform(config.Terraform["infra"].Version, []string{"destroy", "-auto-approve"}, env, nil, config)
	case "native":
		tf = terraform.NewLocalTerraform(config.Terraform["infra"].Version, []string{"destroy", "-auto-approve"}, env, nil, config)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't destroy infra: %w", err)
		}
	default:
		return fmt.Errorf("can't supported %s runtime", config.PreferRuntime)
	}

	ui.Output("Running terraform destroy...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	ui.Output("Terraform destroy completed!\n", terminal.WithSuccessStyle())

	return nil
}

func destroyApp(ui terminal.UI, o *DownOptions) error {
	var m manager.Manager
	var icon string

	m = &ecs.Manager{
		Project: o.Config,
		App:     &config.Ecs{Name: o.AppName},
	}

	if app, ok := o.Config.Serverless[o.AppName]; ok {
		app.Name = o.AppName
		m = &serverless.Manager{
			Project: o.Config,
			App:     app,
		}
		icon = app.Icon
	}
	if app, ok := o.Config.Alias[o.AppName]; ok {
		app.Name = o.AppName
		m = &alias.Manager{
			Project: o.Config,
			App:     app,
		}
		icon = app.Icon
	}
	if app, ok := o.Config.Ecs[o.AppName]; ok {
		app.Name = o.AppName
		m = &ecs.Manager{
			Project: o.Config,
			App:     app,
		}
		icon = app.Icon
	}

	if len(icon) != 0 {
		icon += " "
	}

	ui.Output("Destroying %s%s app...\n", icon, o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	err := m.Destroy(ui)
	if err != nil {
		return fmt.Errorf("can't down: %w", err)
	}

	ui.Output("Destroy app %s%s completed\n", icon, o.AppName, terminal.WithSuccessStyle())

	return nil
}
