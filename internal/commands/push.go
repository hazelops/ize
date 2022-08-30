package commands

import (
	"context"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type PushOptions struct {
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

func NewPushFlags(project *config.Project) *PushOptions {
	return &PushOptions{
		Config: project,
	}
}

func NewCmdPush(project *config.Project) *cobra.Command {
	o := NewPushFlags(project)

	cmd := &cobra.Command{
		Use:               "push [flags] <app name>",
		Example:           pushExample,
		Short:             "push app's image",
		Long:              pushLongDesc,
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

	return cmd
}

func (o *PushOptions) Complete(cmd *cobra.Command) error {
	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *PushOptions) Validate() error {

	return nil
}

func (o *PushOptions) Run() error {
	ui := terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	var m manager.Manager

	if app, ok := o.Config.Serverless[o.AppName]; ok {
		app.Name = o.AppName
		m = &serverless.Manager{
			Project: o.Config,
			App:     app,
		}
	}
	if app, ok := o.Config.Alias[o.AppName]; ok {
		app.Name = o.AppName
		m = &alias.Manager{
			Project: o.Config,
			App:     app,
		}
	}
	if app, ok := o.Config.Ecs[o.AppName]; ok {
		app.Name = o.AppName
		m = &ecs.Manager{
			Project: o.Config,
			App:     app,
		}
	} else {
		m = &ecs.Manager{
			Project: o.Config,
			App:     &config.Ecs{Name: o.AppName},
		}
	}

	return m.Push(ui)
}
