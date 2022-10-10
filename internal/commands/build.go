package commands

import (
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/manager/alias"
	"github.com/hazelops/ize/internal/manager/ecs"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/pkg/logs"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/spf13/cobra"
	"os"
)

type BuildOptions struct {
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

	# Build app for arm64
	ize build <app name> --prefer-runtime docker-arm64
`)

func NewBuildFlags(project *config.Project) *BuildOptions {
	return &BuildOptions{
		Config: project,
	}
}

func NewCmdBuild(project *config.Project) *cobra.Command {
	o := NewBuildFlags(project)

	cmd := &cobra.Command{
		Use:               "build [flags] <app name>",
		Example:           buildExample,
		Short:             "build apps",
		Long:              buildLongDesc,
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

func (o *BuildOptions) Complete(cmd *cobra.Command) error {
	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *BuildOptions) Validate() error {
	return nil
}

func (o *BuildOptions) Run() error {
	ui, cancel := logs.GetLogger(false, o.Config.PlainText, os.Stdout)
	defer cancel()

	var m manager.Manager

	m = &ecs.Manager{
		Project: o.Config,
		App:     &config.Ecs{Name: o.AppName},
	}

	fmt.Println(os.Getenv("IZE_CONFIG_FILE"))

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
	}

	return m.Build(ui)
}
