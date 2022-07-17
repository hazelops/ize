package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

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
}

func NewLogsFlags() *LogsOptions {
	return &LogsOptions{}
}

func NewCmdLogs() *cobra.Command {
	o := NewLogsFlags()

	cmd := &cobra.Command{
		Use:   "logs [app-name]",
		Short: "Stream logs of container in the ECS",
		Args:  cobra.MinimumNArgs(1),
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

func (o *LogsOptions) Complete(cmd *cobra.Command) error {
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

	o.AppName = cmd.Flags().Args()[0]

	return nil
}

func (o *LogsOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("can't validate: env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("can't validate: namespace must be specified\n")
	}

	if len(o.AppName) == 0 {
		return fmt.Errorf("can't validate: app name must be specified\n")
	}
	return nil
}

func (o *LogsOptions) Run() error {
	logGroup := fmt.Sprintf("%s-%s", o.Config.Env, o.AppName)

	lto, err := ecs.New(o.Config.Session).ListTasks(&ecs.ListTasksInput{
		Cluster:       &o.EcsCluster,
		DesiredStatus: aws.String("RUNNING"),
		ServiceName:   &logGroup,
		MaxResults:    aws.Int64(1),
	})

	logrus.Infof("log group: %s, cluster name: %s", logGroup, o.EcsCluster)

	svc := cloudwatchlogs.New(o.Config.Session)

	if err != nil {
		return fmt.Errorf("can't run logs: %w", err)
	}

	taskID := *lto.TaskArns[0]
	taskID = taskID[strings.LastIndex(taskID, "/")+1:]

	var token *string
	logStreamName := fmt.Sprintf("main/%s/%s", o.AppName, taskID)

	for {
		logEvents, err := svc.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  &logGroup,
			LogStreamName: &logStreamName,
			NextToken:     token,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", logStreamName, err)
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
