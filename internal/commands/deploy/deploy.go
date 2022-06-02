package deploy

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/apps"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DeployOptions struct {
	Config  *config.Config
	AppName string
	Tag     string
	Image   string
	App     interface{}
}

var deployLongDesc = templates.LongDesc(`
	Deploy service.
    App name must be specified for a app deploy. 
`)

var deployExample = templates.Examples(`
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

	return cmd
}

func (o *DeployOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't deploy your stack: %w", err)
	}

	viper.BindPFlags(cmd.Flags())
	o.AppName = cmd.Flags().Args()[0]
	viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)

	o.Tag = viper.GetString("tag")

	return nil
}

func (o *DeployOptions) Validate() error {
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

func (o *DeployOptions) Run() error {
	ui := terminal.ConsoleUI(aws.BackgroundContext(), o.Config.IsPlainText)

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

	var deployment apps.App

	switch appType {
	case "ecs":
		deployment = apps.NewECSApp(o.AppName, o.App)
	case "serverless":
		deployment = apps.NewServerlessApp(o.AppName, o.App)
	case "alias":
		deployment = apps.NewAliasApp(o.AppName)
	default:
		return fmt.Errorf("apps type of %s not supported", appType)
	}

	err := deployment.Deploy(ui)
	if err != nil {
		return err
	}

	ui.Output("Deploy app %s completed\n", o.AppName, terminal.WithSuccessStyle())

	return nil
}
