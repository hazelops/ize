package push

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/serverless"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type Options struct {
	Config  *config.Project
	AppName string
	App     interface{}
}

var pushLongDesc = templates.LongDesc(`
	Push app image (so far only ECR).
    App name must be specified for a app image push. 
`)

var pushExample = templates.Examples(`
	# Push app's artifact (Docker image to ECR, for example).
	ize push <app name>

	# Push app's artifact with explicitly specified config file
	ize --config-file (or -c) /path/to/config push <app name>

	# Push app's artifact with explicitly specified config file passed via environment variable.
	export IZE_CONFIG_FILE=/path/to/config
	ize push <app name>
`)

func NewPushFlags() *Options {
	return &Options{}
}

func NewCmdPush() *cobra.Command {
	o := NewPushFlags()

	cmd := &cobra.Command{
		Use:     "push [flags] <app name>",
		Example: pushExample,
		Short:   "push app's image",
		Long:    pushLongDesc,
		Args:    cobra.MaximumNArgs(1),
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

	return cmd
}

func (o *Options) Complete(cmd *cobra.Command) error {
	var err error
	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can`t complete options: %w", err)
	}

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *Options) Validate() error {

	return nil
}

func (o *Options) Run() error {
	ui := terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	var manager manager.Manager

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
		manager = &ecs.Manager{
			Project: o.Config,
			App:     app,
		}
	} else {
		manager = &ecs.Manager{
			Project: o.Config,
			App:     &config.Ecs{Name: o.AppName},
		}
	}

	return manager.Push(ui)
}
