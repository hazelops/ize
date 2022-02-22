package exec

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ExecOptions struct {
	Config      *config.Config
	ServiceName string
	EcsCluster  string
	Command     string
}

func NewExecFlags() *ExecOptions {
	return &ExecOptions{}
}

func NewCmdExec() *cobra.Command {
	o := NewExecFlags()

	cmd := &cobra.Command{
		Use:   "exec [service-name] -- [commands]",
		Short: "exec command in ECS container",
		Long:  "Connect to a container in the ECS service via AWS SSM and run command.\nTakes ECS service name as an argument.",
		Args:  cobra.MinimumNArgs(1),
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

			err = o.Run(cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")

	return cmd
}

func (o *ExecOptions) Complete(cmd *cobra.Command, args []string, argsLenAtDash int) error {
	cfg, err := config.InitializeConfig(config.WithSSMPlugin())
	if err != nil {
		return err
	}

	o.Config = cfg

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	o.ServiceName = cmd.Flags().Args()[0]

	o.Command = strings.Join(args[argsLenAtDash:], " ")

	return nil
}

func (o *ExecOptions) Validate() error {
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

func (o *ExecOptions) Run(cmd *cobra.Command) error {
	serviceName := fmt.Sprintf("%s-%s", o.Config.Env, o.ServiceName)

	logrus.Infof("service name: %s, cluster name: %s", serviceName, o.EcsCluster)
	logrus.Infof("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	ecsSvc := ecs.New(o.Config.Session)

	lto, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:       &o.EcsCluster,
		DesiredStatus: aws.String(ecs.DesiredStatusRunning),
		ServiceName:   &serviceName,
	})
	if err != nil {
		pterm.Error.Printfln("Accessing container")
		return err
	}

	logrus.Debugf("list task output: %s", lto)

	if len(lto.TaskArns) == 0 {
		return fmt.Errorf("running task not found")
	}

	pterm.Success.Printfln("Accessing container")

	out, err := ecsSvc.ExecuteCommand(&ecs.ExecuteCommandInput{
		Container:   &o.ServiceName,
		Interactive: aws.Bool(true),
		Cluster:     &o.EcsCluster,
		Task:        lto.TaskArns[0],
		Command:     aws.String(o.Command),
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
