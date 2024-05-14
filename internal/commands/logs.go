package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LogsOptions struct {
	Config        *config.Project
	AppName       string
	EcsCluster    string
	Task          string
	LogGroupName  string
	LogStreamName string
}

func NewLogsFlags(project *config.Project) *LogsOptions {
	return &LogsOptions{
		Config: project,
	}
}

func NewCmdLogs(project *config.Project) *cobra.Command {
	o := NewLogsFlags(project)

	cmd := &cobra.Command{
		Use:               "logs [app-name]",
		Short:             "Stream logs of container in the ECS",
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
	cmd.Flags().StringVar(&o.Task, "task", "", "set ECS task id")

	return cmd
}

func (o *LogsOptions) Complete(cmd *cobra.Command) error {
	if o.EcsCluster == "" {
		o.EcsCluster = fmt.Sprintf("%s-%s", o.Config.Env, o.Config.Namespace)
	}

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *LogsOptions) Validate() error {
	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified\n")
	}
	return nil
}

func (o *LogsOptions) Run() error {
	var err error

	if len(o.LogGroupName) == 0 {
		o.LogGroupName, err = getEcsServiceLogGroupName(o)
		if err != nil {
			return err
		}
	}

	taskID := o.Task
	if len(taskID) == 0 {
		lto, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String("RUNNING"),
			ServiceName:   &o.LogGroupName,
			MaxResults:    aws.Int64(1),
		})

		logrus.Infof("log group: %s, cluster name: %s", o.LogGroupName, o.EcsCluster)

		if err != nil {
			return fmt.Errorf("can't get logs: %w", err)
		}

		taskID = *lto.TaskArns[0]
		taskID = taskID[strings.LastIndex(taskID, "/")+1:]
	}

	var token *string
	if len(o.LogStreamName) == 0 {
		var logStreamPrefix string
		logStreamPrefix, err = getEcsServiceLogStreamPrefix(o)
		if err != nil {
			logrus.Errorf("can't get log stream prefix: %s", err)
		}

		o.LogStreamName = fmt.Sprintf("%s/%s", logStreamPrefix, taskID)
	}

	GetLogs(o.Config.AWSClient.CloudWatchLogsClient, o.LogGroupName, o.LogStreamName, token)

	return nil
}

func GetLogs(clw cloudwatchlogsiface.CloudWatchLogsAPI, logGroup string, logStreamName string, token *string) {
	for {
		logEvents, err := clw.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &logGroup,
			LogStreamName: &logStreamName,
			NextToken:     token,
		})
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", logStreamName, err)
			continue
		}

		for _, e := range logEvents.Events {
			_, m := formatMessage(e)
			fmt.Println(m)
		}

		token = logEvents.NextForwardToken

		time.Sleep(time.Second * 5)
	}
}

func formatMessage(e *cloudwatchlogs.OutputLogEvent) (t time.Time, m string) {
	m = *e.Message

	if len(m) > 16 {
		if _, err := time.Parse("Jan  2 15:04:05 ", m[:16]); err == nil {
			m = m[16:]
		}
	}

	t = time.Unix(0, *e.Timestamp*1000000)
	m = t.Format("2006-01-02 15:04:05 ") + m
	return
}

func getEcsServiceLogGroupName(o *LogsOptions) (string, error) {
	// TODO: Move core logic to a shared function (since it's used in deploy too)
	ecsServiceLogGroupCandidates := []string{
		o.AppName,
		fmt.Sprintf("%s-%s", o.Config.Env, o.AppName),
		fmt.Sprintf("%s-%s-%s", o.Config.Env, o.Config.Namespace, o.AppName),
	}

	for _, v := range ecsServiceLogGroupCandidates {
		logrus.Debugf("Checking if Log Group %s exists.", v)

		resp, err := o.Config.AWSClient.CloudWatchLogsClient.DescribeLogGroups(&cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(v),
		})

		if len(resp.LogGroups) == 0 {
			logrus.Debug("No log groups with prefix %s. Trying other options", v)
			continue
		}

		for _, logGroup := range resp.LogGroups {
			if aws.StringValue(logGroup.LogGroupName) == v {
				logrus.Debugf("Found Log Group %s", v)
				return v, err
			}
		}

		return v, err
	}
	err := errors.New("Log group not found")
	return "", err
}

func getEcsServiceLogStreamPrefix(o *LogsOptions) (string, error) {
	ecsServiceLogStreamNameCandidates := []string{
		o.AppName,
		fmt.Sprintf("main/%s-%s", o.Config.Namespace, o.AppName),
		fmt.Sprintf("main/%s-%s", o.Config.Env, o.AppName),
		fmt.Sprintf("main/%s-%s-%s", o.Config.Env, o.Config.Namespace, o.AppName),
	}

	for _, v := range ecsServiceLogStreamNameCandidates {
		logrus.Debugf("Checking if logStream %s/* exists", v)
		resp, err := o.Config.AWSClient.CloudWatchLogsClient.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(o.LogGroupName),
			LogStreamNamePrefix: aws.String(v),
		})

		if len(resp.LogStreams) == 0 {
			logrus.Debugf("No log streams with prefix %s. Trying other options", v)
			continue
		}

		for _, logStream := range resp.LogStreams {
			if strings.Contains(aws.StringValue(logStream.LogStreamName), v) {
				logrus.Debugf("Found Log Stream %s", v)
				return v, err
			}
		}

		return v, err
	}
	err := errors.New(fmt.Sprintf("ECS Container for %s not found", o.AppName))
	return "", err
}
