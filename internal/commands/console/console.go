package console

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ConsoleOptions struct {
	Config       *config.Config
	AppName      string
	EcsCluster   string
	Task         string
	CustomPrompt bool
}

func NewConsoleFlags() *ConsoleOptions {
	return &ConsoleOptions{}
}

func NewCmdConsole() *cobra.Command {
	o := NewConsoleFlags()

	cmd := &cobra.Command{
		Use:   "console [app-name]",
		Short: "Connect to a container in the ECS",
		Long:  "Connect to a container of the app via AWS SSM.\nTakes app name that is running on ECS as an argument",
		Args:  cobra.MinimumNArgs(1),
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

	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")
	cmd.Flags().StringVar(&o.Task, "task", "", "set task id")
	cmd.Flags().BoolVar(&o.CustomPrompt, "custom-prompt", false, "enable custom prompt in the console")

	return cmd
}

func (o *ConsoleOptions) Complete(cmd *cobra.Command, args []string) error {
	if err := config.CheckRequirements(config.WithSSMPlugin()); err != nil {
		return err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	if !o.CustomPrompt {
		o.CustomPrompt = viper.GetBool("CUSTOM_PROMPT")
	}

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *ConsoleOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate: namespace must be specified")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified")
	}

	return nil
}

func (o *ConsoleOptions) Run() error {
	appName := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	logrus.Infof("app name: %s, cluster name: %s", appName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	s, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start("Getting access to container...")

	ecsSvc := ecs.New(o.Config.Session)

	if o.Task == "" {
		lto, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String(ecs.DesiredStatusRunning),
			ServiceName:   &appName,
		})
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "ClusterNotFoundException":
				return fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
			default:
				return err
			}
		}

		logrus.Debugf("list task output: %s", lto)

		if len(lto.TaskArns) == 0 {
			return fmt.Errorf("running task not found")
		}

		o.Task = *lto.TaskArns[0]
	}

	s.UpdateText("Executing command...")
	consoleCommand := `/bin/sh`

	if o.CustomPrompt {
		// This is ASCII Prompt string with colors. See https://dev.to/ifenna__/adding-colors-to-bash-scripts-48g4 for reference
		// TODO: Make this customizable via a config
		promptString := fmt.Sprintf(`\e[1;35m★\e[0m $ENV-$APP_NAME\n\e[1;33m\e[0m \w \e[1;34m❯\e[0m `)
		consoleCommand = fmt.Sprintf(`/bin/sh -c '$(echo "export PS1=\"%s\"" > /etc/profile.d/ize.sh) /bin/bash --rcfile /etc/profile'`, promptString)
	}

	out, err := ecsSvc.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.AppName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        &o.Task,
		Command:     aws.String(consoleCommand),
	})
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "ClusterNotFoundException":
			return fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
		default:
			return err
		}
	}

	s.Success()

	ssmCmd := ssmsession.NewSSMPluginCommand(o.Config.AwsRegion)
	ssmCmd.Start((out.Session))
	if err != nil {
		return err
	}

	return nil
}
