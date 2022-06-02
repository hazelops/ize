package push

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

type PushOptions struct {
	Config  *config.Config
	AppName string
	Tag     string
	App     interface{}
}

var pushLongDesc = templates.LongDesc(`
	Push app image (so far only ECR).
    App name must be specified for a app image push. 
`)

var pushExample = templates.Examples(`
	# Push image app (config file required)
	ize push <app name>

	# Push image app via config file
	ize --config-file (or -c) /path/to/config push <app name>

	# Push image app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize push <app name>
`)

func NewPushFlags() *PushOptions {
	return &PushOptions{}
}

func NewCmdPush() *cobra.Command {
	o := NewPushFlags()

	cmd := &cobra.Command{
		Use:     "push [flags] <app name>",
		Example: pushExample,
		Short:   "push app image",
		Long:    pushLongDesc,
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

func (o *PushOptions) Complete(cmd *cobra.Command, args []string) error {
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

func (o *PushOptions) Validate() error {

	return nil
}

func (o *PushOptions) Run() error {
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
