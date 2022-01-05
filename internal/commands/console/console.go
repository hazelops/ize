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
	"github.com/spf13/viper"
)

type ConsoleOptions struct {
	ServiceName string
	EcsCluster  string
	Env         string
	Namespace   string
	Profile     string
	Region      string
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
			err := o.Complete(cmd, args)
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
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Env = viper.GetString("env")
	o.Namespace = viper.GetString("namespace")

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Env, o.Namespace)
	}

	o.ServiceName = cmd.Flags().Args()[0]

	//TODO
	o.Profile = viper.GetString("aws_profile")
	o.Region = viper.GetString("aws_region")

	if o.Region == "" {
		o.Region = viper.GetString("aws-region")
	}

	if o.Profile == "" {
		o.Profile = viper.GetString("aws-profile")
	}

	return nil
}

func (o *ConsoleOptions) Validate() error {
	if len(o.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified")
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("AWS region must be specified")
	}

	if len(o.EcsCluster) == 0 {
		return fmt.Errorf("ECS cluster must be specified")
	}

	if len(o.ServiceName) == 0 {
		return fmt.Errorf("service name must be specified")
	}
	return nil
}

func (o *ConsoleOptions) Run() error {
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	serviceName := fmt.Sprintf("%s-%s", o.Env, o.ServiceName)
	clusterName := fmt.Sprintf("%s-%s", o.Env, o.Namespace)

	logrus.Infof("service name: %s, cluster name: %s", serviceName, clusterName)
	logrus.Infof("region: %s, profile: %s", o.Region, o.Profile)

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
	})
	if err != nil {
		pterm.Error.Printfln("Getting AWS session")
		return err
	}

	pterm.Success.Printfln("Getting AWS session")

	ecsSvc := ecs.New(sess)

	lto, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:       &clusterName,
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
		Cluster:     &clusterName,
		Task:        lto.TaskArns[0],
		Command:     aws.String("/bin/sh"),
	})
	if err != nil {
		pterm.Error.Printfln("Executing command")
		return err
	}

	pterm.Success.Printfln("Executing command")

	ssmCmd := ssmsession.NewSSMPluginCommand(o.Region)
	ssmCmd.Start((out.Session))

	if err != nil {
		return err
	}

	return nil
}
