package commands

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/cirruslabs/echelon"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/logs"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type DownOptions struct {
	Config           *config.Project
	AppName          string
	SkipBuildAndPush bool
	AutoApprove      bool
	SkipGen          bool
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
			return fmt.Errorf("you must specify at least one terraform stack in ize.toml")
		}

		if _, ok := o.Config.Terraform["infra"]; ok {
			if len(o.Config.Terraform["infra"].AwsProfile) == 0 {
				o.Config.Terraform["infra"].AwsProfile = o.Config.AwsProfile
			}

			if len(o.Config.Terraform["infra"].AwsRegion) == 0 {
				o.Config.Terraform["infra"].AwsRegion = o.Config.AwsRegion
			}

			if len(o.Config.Terraform["infra"].Version) == 0 {
				o.Config.Terraform["infra"].Version = o.Config.TerraformVersion
			}
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
	ui, cancel := logs.GetLogger(false, o.Config.PlainText, os.Stdout)
	defer cancel()

	if o.AppName == "" {
		err := destroyAll(ui, o)
		if err != nil {
			return err
		}
	} else {
		err := destroyApp(o.AppName, o.Config, ui)
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

func destroyAll(ui *echelon.Logger, o *DownOptions) error {
	s := ui.Scoped("Destroy apps")

	err := manager.InReversDependencyOrder(aws.BackgroundContext(), o.Config.GetApps(), func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Config.Terraform["infra"].AwsProfile
		return destroyApp(name, o.Config, s)
	})
	if err != nil {
		return err
	}
	s.Finish(true)

	s = ui.Scoped("Destroy infrastructure")
	if _, ok := o.Config.Terraform["infra"]; ok {
		err = destroyInfra("infra", o.Config, o.SkipGen, s)
		if err != nil {
			return err
		}
	}

	err = manager.InReversDependencyOrder(aws.BackgroundContext(), o.Config.GetStates(), func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Config.Terraform[name].AwsProfile

		return destroyInfra(name, o.Config, o.SkipGen, s)
	})

	s.Finish(true)

	return nil
}

func destroyInfra(state string, config *config.Project, skipGen bool, ui *echelon.Logger) error {
	var s *echelon.Logger
	if !skipGen {
		s = ui.Scoped(fmt.Sprintf("Generate terraform file for \"%s\"", state))
		err := GenerateTerraformFiles(state, "", config, s)
		if err != nil {
			return err
		}
	}

	var tf terraform.Terraform

	ui.Debugf("%v", tf)

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
		tf = terraform.NewDockerTerraform(state, []string{"init", "-input=true"}, env, nil, config)
	case "native":
		tf = terraform.NewLocalTerraform(state, []string{"init", "-input=true"}, env, nil, config)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't destroy infra: %w", err)
		}
	default:
		return fmt.Errorf("can't supported %s runtime", config.PreferRuntime)
	}

	sg := ui.Scoped(fmt.Sprintf("[%s][%s] Run destroy infra", config.Env, state))
	defer sg.Finish(false)

	s = sg.Scoped("Execute terraform init")
	err = tf.RunUI(s)
	if err != nil {
		return err
	}
	s.Finish(true)

	//terraform destroy run options
	tf.NewCmd([]string{"destroy", "-auto-approve"})

	s = sg.Scoped("Execute terraform destroy")
	err = tf.RunUI(s)
	if err != nil {
		return fmt.Errorf("can't deploy infra: %w", err)
	}
	s.Finish(true)
	sg.Finish(true)

	return nil
}

func destroyApp(name string, cfg *config.Project, ui *echelon.Logger) error {
	var m manager.Manager
	var icon string

	m = &ecs.Manager{
		Project: cfg,
		App:     &config.Ecs{Name: name},
	}

	if app, ok := cfg.Serverless[name]; ok {
		app.Name = name
		m = &serverless.Manager{
			Project: cfg,
			App:     app,
		}
		icon = app.Icon
	}
	if app, ok := cfg.Alias[name]; ok {
		app.Name = name
		m = &alias.Manager{
			Project: cfg,
			App:     app,
		}
		icon = app.Icon
	}
	if app, ok := cfg.Ecs[name]; ok {
		app.Name = name
		m = &ecs.Manager{
			Project: cfg,
			App:     app,
		}
		icon = app.Icon
	}

	if len(icon) != 0 {
		icon += " "
	}

	s := ui.Scoped(fmt.Sprintf("Destroy %s%s app", icon, name))

	err := m.Destroy(s)
	if err != nil {
		return fmt.Errorf("can't down: %w", err)
	}
	s.Finish(true)

	time.Sleep(time.Millisecond)

	return nil
}
