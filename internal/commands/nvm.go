package commands

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager/serverless"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type NvmOptions struct {
	Config  *config.Project
	AppName string
	Command []string
}

var nvmLongDesc = templates.LongDesc(`
	Run nvm with the specified command for app.
    Command must be specified for a command run. 
    App name must be specified for a command run. 
`)

var nvmExample = templates.Examples(`
	# Run nvm with command (config file required)
	ize nvm <app name> -- [command]

	# Run nvm with command via config file
	ize --config-file (or -c) /path/to/config nvm <app name> -- [command]

	# Run nvm with command via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize nvm <app name> -- [command]
`)

func NewNvmFlags(project *config.Project) *NvmOptions {
	return &NvmOptions{
		Config: project,
	}
}

func NewCmdNvm(project *config.Project) *cobra.Command {
	o := NewNvmFlags(project)

	cmd := &cobra.Command{
		Use:               "nvm [app-name] -- [commands]",
		Example:           nvmExample,
		Short:             "Run nvm with the specified command for app",
		Long:              nvmLongDesc,
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if len(cmd.Flags().Args()) == 0 {
				return fmt.Errorf("app name must be specified")
			}

			argsLenAtDash := cmd.ArgsLenAtDash()
			err := o.Complete(cmd, args, argsLenAtDash)
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

func (o *NvmOptions) Complete(cmd *cobra.Command, args []string, argsLenAtDash int) error {
	if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile(), requirements.WithNVM()); err != nil {
		return err
	}

	if argsLenAtDash > -1 {
		o.Command = args[argsLenAtDash:]
	}

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *NvmOptions) Validate() error {
	if len(o.Command) == 0 {
		return fmt.Errorf("can't validate: you must specify at least one command for the container")
	}

	return nil
}

func (o *NvmOptions) Run() error {
	ui := terminal.ConsoleUI(aws.BackgroundContext(), o.Config.PlainText)

	sg := ui.StepGroup()
	defer sg.Wait()

	var m *serverless.Manager

	if app, ok := o.Config.Serverless[o.AppName]; ok {
		app.Name = o.AppName
		m = &serverless.Manager{
			Project: o.Config,
			App:     app,
		}
	} else {
		return fmt.Errorf("%s not found in config file", o.AppName)
	}

	err := m.Nvm(ui, o.Command)
	if err != nil {
		return err
	}

	return nil
}
