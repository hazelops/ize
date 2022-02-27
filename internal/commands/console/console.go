package console

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ConsoleOptions struct {
	Config      *config.Config
	ServiceName string
	EcsCluster  string
}

func NewConsoleFlags() *ConsoleOptions {
	return &ConsoleOptions{}
}

func NewCmdConsole(ui terminal.UI) *cobra.Command {
	o := NewConsoleFlags()

	cmd := &cobra.Command{
		Use:   "console [service-name]",
		Short: "connect to a container in the ECS",
		Long:  "Connect to a container in the ECS service via AWS SSM.\nTakes ECS service name as an argument.",
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

			err = o.Run(ui)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")

	return cmd
}

func (o *ConsoleOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig(config.WithSSMPlugin())
	if err != nil {
		return err
	}

	o.Config = cfg

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	o.ServiceName = cmd.Flags().Args()[0]

	return nil
}

func (o *ConsoleOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate: env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate: namespace must be specified\n")
	}

	if len(o.ServiceName) == 0 {
		return fmt.Errorf("can't validate: service name must be specified\n")
	}
	return nil
}

func (o *ConsoleOptions) Run(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	serviceName := fmt.Sprintf("%s-%s", o.Config.Env, o.ServiceName)

	logrus.Infof("service name: %s, cluster name: %s", serviceName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	s := sg.Add("accessing container...")
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 50) }()

	ecsSvc := ecs.New(o.Config.Session)

	lto, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:       &o.EcsCluster,
		DesiredStatus: aws.String(ecs.DesiredStatusRunning),
		ServiceName:   &serviceName,
	})
	if err != nil {
		return err
	}

	logrus.Debugf("list task output: %s", lto)

	if len(lto.TaskArns) == 0 {
		return fmt.Errorf("running task not found\n")
	}

	s.Done()
	s = sg.Add("executing command...")

	out, err := ecsSvc.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.ServiceName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        lto.TaskArns[0],
		Command:     aws.String("/bin/sh"),
	})
	if err != nil {
		return err
	}

	s.Done()

	ssmCmd := ssmsession.NewSSMPluginCommand(o.Config.AwsRegion)
	ssmCmd.Start((out.Session))
	if err != nil {
		return err
	}

	return nil
}
