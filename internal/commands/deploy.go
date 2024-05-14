package commands

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/helm"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type DeployOptions struct {
	Config                 *config.Project
	AppName                string
	Image                  string
	TaskDefinitionRevision string
	Unsafe                 bool
	Force                  bool
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

func NewDeployFlags(project *config.Project) *DeployOptions {
	return &DeployOptions{
		Config: project,
	}
}

func NewCmdDeploy(project *config.Project) *cobra.Command {
	o := NewDeployFlags(project)

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
	cmd.Flags().BoolVar(&o.Unsafe, "unsafe", false, "set unsafe healthcheck options (accelerates deployment if possible)")
	cmd.Flags().BoolVar(&o.Force, "force", false, "forces a deployment to take place (only serverless)")

	return cmd
}

func (o *DeployOptions) Complete(cmd *cobra.Command) error {
	if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile()); err != nil {
		return err
	}

	if len(o.Config.Serverless) != 0 {
		if err := requirements.CheckRequirements(requirements.WithNVM()); err != nil {
			return err
		}
	}

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *DeployOptions) Validate() error {
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

func (o *DeployOptions) Run() error {
	ui := terminal.ConsoleUI(aws.BackgroundContext(), o.Config.PlainText)

	ui.Output("Deploying %s app...\n", o.AppName, terminal.WithHeaderStyle())
	sg := ui.StepGroup()
	defer sg.Wait()

	var m manager.Manager
	var providerUsed string
	// Note, Viper doesn't read empty TOML sections (https://github.com/spf13/viper/issues/1131 so if there are no app sections, we'll use apps_provider
	logrus.Debugf("FYI, Viper can't read/see empty TOML sections. If they are empty, we'll try to use `apps_provider` config if it's set in ize.toml. See more here https://github.com/spf13/viper/issues/1131")
	//if o.Config.AppsProvider == "helm" {
	//	logrus.Debugf("Found helm app")
	//

	if app, ok := o.Config.Ecs[o.AppName]; o.Config.AppsProvider == "ecs" || ok {
		providerUsed = "ecs"
		app.Name = o.AppName
		m = &ecs.Manager{
			Project: o.Config,
			App:     app,
		}
	}

	if app, ok := o.Config.Helm[o.AppName]; o.Config.AppsProvider == "helm" || ok {
		providerUsed = "helm"
		app.Name = o.AppName
		app.Force = o.Force
		m = &helm.Manager{
			Project: o.Config,
			App:     app,
		}
	}

	if app, ok := o.Config.Serverless[o.AppName]; o.Config.AppsProvider == "serverless" || ok {
		providerUsed = "serverless"
		app.Name = o.AppName
		app.Force = o.Force

		m = &serverless.Manager{
			Project: o.Config,
			App:     app,
		}
	}

	if _, ok := o.Config.Alias[o.AppName]; o.Config.AppsProvider == "alias" || ok {
		providerUsed = "alias"
		m = &alias.Manager{
			Project: o.Config,
			App: &config.Alias{
				Name: o.AppName,
			},
		}
	}

	if len(o.TaskDefinitionRevision) != 0 {
		err := m.Redeploy(ui)
		if err != nil {
			return err
		}

		ui.Output("Redeploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

		return nil
	}

	logrus.Debugf("Deploying using %s. (default_app_provier=%s)", providerUsed, o.Config.AppsProvider)
	err := m.Deploy(ui)
	if err != nil {
		return err
	}

	ui.Output("Deploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}
