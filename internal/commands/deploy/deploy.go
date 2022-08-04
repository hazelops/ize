package deploy

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Options struct {
	Config                 *config.Project
	AppName                string
	Image                  string
	App                    interface{}
	TaskDefinitionRevision string
	Unsafe                 bool
}

var deployLongDesc = templates.LongDesc(`
	Deploy service.
    App name must be specified for a app deploy. 

	If you install a revision of the task definition, the application will be redeployed (ECS only).
	Warning: Redeployment using the docker runtime, a new task definition will be deployed based on the specified revision.
`)

var deployExample = templates.Examples(`
	# Deploy app (config file required)
	ize deploy <app name>

	# Deploy app via config file
	ize --config-file (or -c) /path/to/config deploy <app name>

	# Deploy app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize deploy <app name>

	# Redeploy app (ECS only)
	ize deploy <app name> --task-definition-revision <task definition revision>
`)

func NewDeployFlags() *Options {
	return &Options{}
}

func NewCmdDeploy() *cobra.Command {
	o := NewDeployFlags()

	cmd := &cobra.Command{
		Use:               "deploy [flags] <app name>",
		Example:           deployExample,
		Short:             "Manage deployments",
		Long:              deployLongDesc,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			err := o.Complete(cmd)
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

	cmd.Flags().StringVar(&o.TaskDefinitionRevision, "task-definition-revision", "", "set task definition revision (ECS only)")
	cmd.Flags().BoolVar(&o.Unsafe, "unsafe", false, "set unsafe healtcheck options (accelerates deployment if possible)")

	return cmd
}

func (o *Options) Complete(cmd *cobra.Command) error {
	var err error

	if err = config.CheckRequirements(config.WithIzeStructure(), config.WithConfigFile()); err != nil {
		return err
	}

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't deploy your stack: %w", err)
	}

	if o.Config.Serverless != nil {
		if err = config.CheckRequirements(config.WithNVM()); err != nil {
			return err
		}
	}

	viper.BindPFlags(cmd.Flags())
	o.AppName = cmd.Flags().Args()[0]
	viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)

	return nil
}

func (o *Options) Validate() error {
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

func (o *Options) Run() error {
	ui := terminal.ConsoleUI(aws.BackgroundContext(), o.Config.PlainText)

	ui.Output("Deploying %s app...\n", o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	var manager manager.Manager

	manager = &ecs.Manager{
		Project: o.Config,
		App: &config.Ecs{
			Name:                   o.AppName,
			TaskDefinitionRevision: o.TaskDefinitionRevision,
			Unsafe:                 o.Unsafe,
		},
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
		app.TaskDefinitionRevision = o.TaskDefinitionRevision
		app.Unsafe = o.Unsafe
		manager = &ecs.Manager{
			Project: o.Config,
			App:     app,
		}
	}

	if len(o.TaskDefinitionRevision) != 0 {
		err := manager.Redeploy(ui)
		if err != nil {
			return err
		}

		ui.Output("Redeploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

		return nil
	}

	err := manager.Deploy(ui)
	if err != nil {
		return err
	}

	ui.Output("Deploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}
