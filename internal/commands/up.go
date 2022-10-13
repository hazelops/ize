package commands

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type UpOptions struct {
	Config           *config.Project
	SkipBuildAndPush bool
	SkipGen          bool
	AutoApprove      bool
	UI               terminal.UI
	Apps             []string
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

func NewUpFlags(project *config.Project) *UpOptions {
	return &UpOptions{
		Config: project,
	}
}

func NewCmdUp(project *config.Project) *cobra.Command {
	o := NewUpFlags(project)

	cmd := &cobra.Command{
		Use:               "up [flags] <app name>",
		Example:           upExample,
		Short:             "Bring full application up from the bottom to the top.",
		Long:              upLongDesc,
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
		NewCmdUpInfra(project),
		NewCmdUpApps(project),
	)

	return cmd
}

func (o *UpOptions) Complete(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile()); err != nil {
			return err
		}

		if len(o.Config.Serverless) != 0 {
			if err := requirements.CheckRequirements(requirements.WithNVM()); err != nil {
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
			if err := requirements.CheckRequirements(requirements.WithNVM()); err != nil {
				return err
			}
		}

		o.Apps = cmd.Flags().Args()
	}

	o.UI = terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	return nil
}

func (o *UpOptions) Validate() error {
	if len(o.Apps) > 0 {
		err := o.validate()
		if err != nil {
			return err
		}
	} else {
		err := o.validateAll()
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *UpOptions) Run() error {
	ui := o.UI
	if len(o.Apps) > 0 {
		err := manager.InDependencyOrder(aws.BackgroundContext(), o.Config.GetStates(o.Apps...), func(c context.Context, name string) error {
			return deployInfra(name, ui, o.Config, o.SkipGen)
		})
		if err != nil {
			return err
		}
		err = manager.InDependencyOrder(aws.BackgroundContext(), o.Config.GetApps(o.Apps...), func(c context.Context, name string) error {
			return deployApp(name, ui, o.Config)
		})
		if err != nil {
			return err
		}

	} else {
		err := deployAll(ui, o)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *UpOptions) validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	return nil
}

func (o *UpOptions) validateAll() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate options: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate options: namespace must be specified")
	}

	return nil
}

func deployAll(ui terminal.UI, o *UpOptions) error {
	if _, ok := o.Config.Terraform["infra"]; ok {
		err := deployInfra("infra", ui, o.Config, o.SkipGen)
		if err != nil {
			return err
		}
	}

	err := manager.InDependencyOrder(aws.BackgroundContext(), o.Config.GetStates(), func(c context.Context, name string) error {
		return deployInfra(name, ui, o.Config, o.SkipGen)
	})
	if err != nil {
		return err
	}

	ui.Output("Deploying apps...", terminal.WithHeaderStyle())

	err = manager.InDependencyOrder(aws.BackgroundContext(), o.Config.GetApps(), func(c context.Context, name string) error {
		o.Config.AwsProfile = o.Config.Terraform["infra"].AwsProfile

		err := deployApp(name, ui, o.Config)
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
