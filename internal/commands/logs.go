package commands

import (
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
	Config     *config.Project
	AppName    string
	EcsCluster string
	Task       string
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
	logGroup := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	taskID := o.Task
	if len(taskID) == 0 {
		lto, err := o.Config.AWSClient.ECSClient.ListTasks(&ecs.ListTasksInput{
			Cluster:       &o.EcsCluster,
			DesiredStatus: aws.String("RUNNING"),
			ServiceName:   &logGroup,
			MaxResults:    aws.Int64(1),
		})

		logrus.Infof("log group: %s, cluster name: %s", logGroup, o.EcsCluster)

		if err != nil {
			return fmt.Errorf("can't run logs: %w", err)
		}

		taskID = *lto.TaskArns[0]
		taskID = taskID[strings.LastIndex(taskID, "/")+1:]
	}

	var token *string
	logStreamName := fmt.Sprintf("main/%s/%s", o.AppName, taskID)

	GetLogs(o.Config.AWSClient.CloudWatchLogsClient, logGroup, logStreamName, token)

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
