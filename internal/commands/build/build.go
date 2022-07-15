package build

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
}

var buildLongDesc = templates.LongDesc(`
	Build app.
    App name must be specified for a app build. 
`)

var buildExample = templates.Examples(`
	# Build app (config file required)
	ize build <app name>

	# Build app with explicitly specified config file
	ize --config-file (or -c) /path/to/config build <app name>

	#  Build app with explicitly specified config file passed via environment variable.
	export IZE_CONFIG_FILE=/path/to/config
	ize build <app name>
`)

func NewBuildFlags() *Options {
	return &Options{}
}

func NewCmdBuild() *cobra.Command {
	o := NewBuildFlags()

	cmd := &cobra.Command{
		Use:     "build [flags] <app name>",
		Example: buildExample,
		Short:   "build apps",
		Long:    buildLongDesc,
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
		return fmt.Errorf("can't load options for a command: %w", err)
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

	return manager.Build(ui)
}
