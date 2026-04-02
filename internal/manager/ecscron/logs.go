package ecscron

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
)

// GetLastRunLogs fetches CloudWatch logs from the most recent task execution.
func (m *Manager) GetLastRunLogs(ui terminal.UI) error {
	m.prepare()

	if err := m.ensureSession(); err != nil {
		return err
	}

	cwl := m.Project.AWSClient.CloudWatchLogsClient
	logGroup := fmt.Sprintf("%s-%s", m.Project.Env, m.App.Name)

	out, err := cwl.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroup,
		Limit:        aws.Int64(1),
		Descending:   aws.Bool(true),
		OrderBy:      aws.String("LastEventTime"),
	})
	if err != nil {
		return fmt.Errorf("can't describe log streams for %s: %w", logGroup, err)
	}

	if len(out.LogStreams) == 0 {
		pterm.Info.Printfln("No log streams found for %s", logGroup)
		return nil
	}

	stream := out.LogStreams[0]
	pterm.DefaultSection.Printfln("Logs from %s (stream: %s):", logGroup, *stream.LogStreamName)

	events, err := cwl.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  &logGroup,
		LogStreamName: stream.LogStreamName,
	})
	if err != nil {
		return fmt.Errorf("can't get log events: %w", err)
	}

	for _, event := range events.Events {
		pterm.Println("| " + *event.Message)
	}

	return nil
}
