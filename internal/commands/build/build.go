package build

import (
	"context"
	"fmt"

	"github.com/hazelops/ize/internal/apps"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BuildOptions struct {
	Config  *config.Config
	AppName string
	Tag     string
	App     interface{}
}

var buildLongDesc = templates.LongDesc(`
	Build sevice.
    App name must be specified for a app build. 
`)

var buildExample = templates.Examples(`
	# Build app (config file required)
	ize build <app name>

	# Build app via config file
	ize --config-file (or -c) /path/to/config build <app name>

	# Build app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize build <app name>
`)

func NewBuildFlags() *BuildOptions {
	return &BuildOptions{}
}

func NewCmdBuild() *cobra.Command {
	o := NewBuildFlags()

	cmd := &cobra.Command{
		Use:     "build [flags] <app name>",
		Example: buildExample,
		Short:   "manage builds",
		Long:    buildLongDesc,
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

func (o *BuildOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can`t complete options: %w", err)
	}

	viper.BindPFlags(cmd.Flags())
	o.AppName = cmd.Flags().Args()[0]
	viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)

	o.Tag = viper.GetString("tag")

	return nil
}

func (o *BuildOptions) Validate() error {

	return nil
}

func (o *BuildOptions) Run() error {
	ui := terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)

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

	return app.Push(ui)
}
