package commands

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

type ExecOptions struct {
	Config        *config.Project
	AppName       string
	EcsCluster    string
	Command       []string
	Task          string
	ContainerName string
}

var execExample = templates.Examples(`
	# Connect to a container in the ECS via AWS SSM and run command.
	ize exec goblin ps aux
`)

func NewExecFlags(project *config.Project) *ExecOptions {
	return &ExecOptions{
		Config: project,
	}
}

func NewCmdExec(project *config.Project) *cobra.Command {
	o := NewExecFlags(project)

	cmd := &cobra.Command{
		Use:               "exec [app-name] -- [commands]",
		Example:           execExample,
		Short:             "Execute command in ECS container",
		Long:              "Connect to a container in the ECS via AWS SSM and run command.\nIt uses app name as an argument.",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: config.GetApps,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
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

	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")
	cmd.Flags().StringVar(&o.Task, "task", "", "set task id")
	cmd.Flags().StringVar(&o.ContainerName, "container-name", "", "set container name")

	return cmd
}

func (o *ExecOptions) Complete(cmd *cobra.Command, args []string, argsLenAtDash int) error {
	if err := requirements.CheckRequirements(requirements.WithSSMPlugin()); err != nil {
		return err
	}

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	o.AppName = cmd.Flags().Args()[0]

	if len(o.ContainerName) == 0 {
		o.ContainerName = o.AppName
	}

	if argsLenAtDash > -1 {
		o.Command = args[argsLenAtDash:]
	}

	return nil
}

func (o *ExecOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate: namespace must be specified")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified")
	}

	if len(o.Command) == 0 {
		return fmt.Errorf("can't validate: you must specify at least one command for the container")
	}

	return nil
}

func (o *ExecOptions) Run() error {
	appName := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	logrus.Infof("app name: %s, cluster name: %s", appName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	s, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start("Getting access to container...")

	if o.Task == "" {
		lto, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
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

	out, err := o.Config.AWSClient.ECSClient.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.AppName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        &o.Task,
		Command:     aws.String(strings.Join(o.Command, " ")),
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
	err = ssmCmd.Start(out.Session)
	if err != nil {
		return err
	}

	return nil
}
