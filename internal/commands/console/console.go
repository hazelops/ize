package console

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/pterm/pterm"
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

func NewCmdConsole() *cobra.Command {
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

			err = o.Run()
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
	cfg, err := config.InitializeConfig()
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
		return fmt.Errorf("can't validate: env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate: namespace must be specified")
	}

	if len(o.ServiceName) == 0 {
		return fmt.Errorf("can't validate: service name must be specified")
	}
	return nil
}

func (o *ConsoleOptions) Run() error {
	serviceName := fmt.Sprintf("%s-%s", o.Config.Env, o.ServiceName)

	logrus.Infof("service name: %s, cluster name: %s", serviceName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		return fmt.Errorf("can't run console: failed to create aws session")
	}

	ecsSvc := ecs.New(sess)

	lto, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:       &o.EcsCluster,
		DesiredStatus: aws.String(ecs.DesiredStatusRunning),
		ServiceName:   &serviceName,
	})
	if err != nil {
		pterm.Error.Printfln("Getting running task")
		return err
	}

	logrus.Debugf("list task output: %s", lto)

	if len(lto.TaskArns) == 0 {
		return fmt.Errorf("running task not found")
	}

	pterm.Success.Printfln("Getting running task")

	out, err := ecsSvc.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.ServiceName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        lto.TaskArns[0],
		Command:     aws.String("/bin/sh"),
	})
	if err != nil {
		pterm.Error.Printfln("Executing command")
		return err
	}

	pterm.Success.Printfln("Executing command")

	ssmCmd := ssmsession.NewSSMPluginCommand(o.Config.AwsRegion)
	ssmCmd.Start((out.Session))
	if err != nil {
		return err
	}

	return nil
}
