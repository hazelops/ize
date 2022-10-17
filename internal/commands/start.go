package commands

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type StartOptions struct {
	Config     *config.Project
	AppName    string
	EcsCluster string
}

type NetworkConfiguration struct {
	SecurityGroups struct {
		Value string `json:"value"`
	} `json:"security_groups"`
	Subnets struct {
		Value [][]string `json:"value"`
	} `json:"subnets"`
	VpcPrivateSubnets struct {
		Value []string `json:"value"`
	} `json:"vpc_private_subnets"`
	VpcPublicSubnets struct {
		Value []string `json:"value"`
	} `json:"vpc_public_subnets"`
}

var startExample = templates.Examples(`
	# Connect to a container in the ECS via AWS SSM and run command.
	ize start goblin
`)

func NewStartFlags(project *config.Project) *StartOptions {
	return &StartOptions{
		Config: project,
	}
}

func NewCmdStart(project *config.Project) *cobra.Command {
	o := NewStartFlags(project)

	cmd := &cobra.Command{
		Use:               "start [app-name]",
		Example:           startExample,
		Short:             "Start ECS task",
		Long:              "Start ECS task and stream logs until it dies or canceled.\nIt uses app name as an argument.",
		Args:              cobra.MinimumNArgs(1),
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

	cmd.Flags().StringVar(&o.EcsCluster, "ecs-cluster", "", "set ECS cluster name")

	return cmd
}

func (o *StartOptions) Complete(cmd *cobra.Command) error {
	if err := requirements.CheckRequirements(requirements.WithSSMPlugin()); err != nil {
		return err
	}

	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *StartOptions) Validate() error {
	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified")
	}

	return nil
}

func getNetworkConfiguration(svc ssmiface.SSMAPI, env string) (NetworkConfiguration, error) {
	resp, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return NetworkConfiguration{}, fmt.Errorf("can't get terraform output: %w", err)
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		return NetworkConfiguration{}, fmt.Errorf("can't get terraform output: %w", err)
	}

	var output NetworkConfiguration

	err = json.Unmarshal(value, &output)
	if err != nil {
		return NetworkConfiguration{}, fmt.Errorf("can't get network configuration: %w", err)
	}

	return output, nil
}

func (o *StartOptions) Run() error {
	ctx := context.Background()

	appName := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)
	logGroup := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	logrus.Debugf("app name: %s, cluster name: %s", appName, o.EcsCluster)
	logrus.Debugf("region: %s, profile: %s", o.Config.AwsProfile, o.Config.AwsRegion)

	configuration, err := getNetworkConfiguration(o.Config.AWSClient.SSMClient, o.Config.Env)
	if err != nil {
		return err
	}

	logrus.Debugf("network configuration: %+v", configuration)

	if len(configuration.VpcPrivateSubnets.Value) == 0 {
		return fmt.Errorf("output private_subnets is missing. Please add it to your Terraform")
	}

	if len(configuration.SecurityGroups.Value) == 0 {
		return fmt.Errorf("output security_groups is missing. Please add it to your Terraform")
	}

	out, err := o.Config.AWSClient.ECSClient.RunTaskWithContext(ctx, &ecs.RunTaskInput{
		TaskDefinition: &appName,
		StartedBy:      aws.String("IZE"),
		Cluster:        &o.EcsCluster,
		NetworkConfiguration: &ecs.NetworkConfiguration{AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
			Subnets: aws.StringSlice(configuration.VpcPrivateSubnets.Value),
		}},
		LaunchType: aws.String(ecs.LaunchTypeFargate),
	})
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case "ClusterNotFoundException":
			return fmt.Errorf("ECS cluster %s not found", o.EcsCluster)
		default:
			return err
		}
	}

	taskID := getTaskID(*out.Tasks[0].TaskArn)

	c := make(chan os.Signal)
	ch := make(chan bool)
	errorChannel := make(chan error)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		s, _ := pterm.DefaultSpinner.WithRemoveWhenDone().Start(fmt.Sprintf("Please wait until task %s running...", appName))

		err := o.Config.AWSClient.ECSClient.WaitUntilTasksRunningWithContext(ctx, &ecs.DescribeTasksInput{
			Cluster: &o.EcsCluster,
			Tasks:   aws.StringSlice([]string{*out.Tasks[0].TaskArn}),
		})
		if err != nil {
			errorChannel <- err
		}

		s.Success()
		pterm.DefaultSection.Println("Logs:")

		var token *string
		go GetLogs(o.Config.AWSClient.CloudWatchLogsClient, logGroup, fmt.Sprintf("main/%s/%s", o.AppName, taskID), token)
		err = o.Config.AWSClient.ECSClient.WaitUntilTasksStoppedWithContext(ctx, &ecs.DescribeTasksInput{
			Cluster: &o.EcsCluster,
			Tasks:   aws.StringSlice([]string{*out.Tasks[0].TaskArn}),
		})
		if err != nil {
			errorChannel <- err
		}
		ch <- true
	}()

	select {
	case <-c:
		fmt.Print("\r")
		_, err := o.Config.AWSClient.ECSClient.StopTask(&ecs.StopTaskInput{
			Cluster: &o.EcsCluster,
			Reason:  aws.String("Task stopped by IZE"),
			Task:    out.Tasks[0].TaskArn,
		})
		if err != nil {
			return err
		}
		pterm.Success.Printfln("Stop task %s by interrupt", appName)
	case <-ch:
		tasks, err := o.Config.AWSClient.ECSClient.DescribeTasks(&ecs.DescribeTasksInput{
			Cluster: &o.EcsCluster,
			Tasks:   aws.StringSlice([]string{taskID}),
		})
		if err != nil {
			return err
		}

		sr := *tasks.Tasks[0].StoppedReason
		st := *tasks.Tasks[0].StopCode
		logrus.Debugf("stop code: %s", st)
		pterm.Success.Printfln("%s was stopped with reason: %s\n", appName, sr)
		return nil
	case err := <-errorChannel:
		return err
	}

	return nil
}

func getTaskID(taskArn string) string {
	split := strings.Split(taskArn, "/")
	return split[len(split)-1]
}
