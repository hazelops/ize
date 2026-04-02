package commands

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	ecscron "github.com/hazelops/ize/internal/manager/ecscron"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type CronOptions struct {
	Config  *config.Project
	AppName string
}

func NewCmdCron(project *config.Project) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "Manage ECS cron tasks",
	}

	cmd.AddCommand(
		NewCmdCronRun(project),
		NewCmdCronLogs(project),
	)

	return cmd
}

func NewCmdCronRun(project *config.Project) *cobra.Command {
	o := &CronOptions{Config: project}

	cmd := &cobra.Command{
		Use:               "run [app-name]",
		Short:             "Manually run a cron task outside its schedule",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			o.AppName = args[0]

			app, ok := o.Config.EcsCron[o.AppName]
			if !ok {
				return fmt.Errorf("ecs-cron app %s not found in config", o.AppName)
			}

			app.Name = o.AppName
			m := &ecscron.Manager{Project: o.Config, App: app}
			ui := terminal.ConsoleUI(aws.BackgroundContext(), o.Config.PlainText)
			return m.RunNow(ui)
		},
	}

	return cmd
}

func NewCmdCronLogs(project *config.Project) *cobra.Command {
	o := &CronOptions{Config: project}

	cmd := &cobra.Command{
		Use:               "logs [app-name]",
		Short:             "Show logs from the last cron task execution",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			o.AppName = args[0]

			app, ok := o.Config.EcsCron[o.AppName]
			if !ok {
				return fmt.Errorf("ecs-cron app %s not found in config", o.AppName)
			}

			app.Name = o.AppName
			m := &ecscron.Manager{Project: o.Config, App: app}
			ui := terminal.ConsoleUI(aws.BackgroundContext(), o.Config.PlainText)
			return m.GetLastRunLogs(ui)
		},
	}

	return cmd
}
