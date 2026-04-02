package ecscron

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
)

// --- Mock ECS ---

type mockECS struct {
	ecsiface.ECSAPI
	listTDOut      *ecs.ListTaskDefinitionsOutput
	describeTDOut  *ecs.DescribeTaskDefinitionOutput
	registerTDOut  *ecs.RegisterTaskDefinitionOutput
	deregisterErr  error
	runTaskOut     *ecs.RunTaskOutput
	runTaskErr     error
}

func (m *mockECS) ListTaskDefinitions(*ecs.ListTaskDefinitionsInput) (*ecs.ListTaskDefinitionsOutput, error) {
	return m.listTDOut, nil
}

func (m *mockECS) DescribeTaskDefinition(*ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return m.describeTDOut, nil
}

func (m *mockECS) RegisterTaskDefinition(*ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return m.registerTDOut, nil
}

func (m *mockECS) DeregisterTaskDefinition(*ecs.DeregisterTaskDefinitionInput) (*ecs.DeregisterTaskDefinitionOutput, error) {
	return nil, m.deregisterErr
}

func (m *mockECS) RunTaskWithContext(ctx context.Context, input *ecs.RunTaskInput, opts ...request.Option) (*ecs.RunTaskOutput, error) {
	if m.runTaskErr != nil {
		return nil, m.runTaskErr
	}
	return m.runTaskOut, nil
}

// --- Mock EventBridge ---

type mockEventBridge struct {
	eventbridgeiface.EventBridgeAPI
	putRuleErr      error
	putTargetsErr   error
	describeRuleOut *eventbridge.DescribeRuleOutput
	describeRuleErr error
	listTargetsOut  *eventbridge.ListTargetsByRuleOutput
	listTargetsErr  error
	removeTargErr   error
	deleteRuleErr   error
}

func (m *mockEventBridge) PutRule(*eventbridge.PutRuleInput) (*eventbridge.PutRuleOutput, error) {
	return &eventbridge.PutRuleOutput{}, m.putRuleErr
}

func (m *mockEventBridge) PutTargets(*eventbridge.PutTargetsInput) (*eventbridge.PutTargetsOutput, error) {
	return &eventbridge.PutTargetsOutput{}, m.putTargetsErr
}

func (m *mockEventBridge) DescribeRule(*eventbridge.DescribeRuleInput) (*eventbridge.DescribeRuleOutput, error) {
	return m.describeRuleOut, m.describeRuleErr
}

func (m *mockEventBridge) ListTargetsByRule(*eventbridge.ListTargetsByRuleInput) (*eventbridge.ListTargetsByRuleOutput, error) {
	return m.listTargetsOut, m.listTargetsErr
}

func (m *mockEventBridge) RemoveTargets(*eventbridge.RemoveTargetsInput) (*eventbridge.RemoveTargetsOutput, error) {
	return &eventbridge.RemoveTargetsOutput{}, m.removeTargErr
}

func (m *mockEventBridge) DeleteRule(*eventbridge.DeleteRuleInput) (*eventbridge.DeleteRuleOutput, error) {
	return &eventbridge.DeleteRuleOutput{}, m.deleteRuleErr
}

// --- Mock SSM ---

type mockSSM struct {
	ssmiface.SSMAPI
}

func (m *mockSSM) GetParameter(*ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	nc := networkConfig{}
	nc.VpcPrivateSubnets.Value = []string{"subnet-abc123"}
	data, _ := json.Marshal(nc)
	encoded := base64.StdEncoding.EncodeToString(data)
	return &ssm.GetParameterOutput{
		Parameter: &ssm.Parameter{Value: &encoded},
	}, nil
}

// --- Mock STS ---

type mockSTS struct {
	stsiface.STSAPI
}

func (m *mockSTS) GetCallerIdentity(*sts.GetCallerIdentityInput) (*sts.GetCallerIdentityOutput, error) {
	return &sts.GetCallerIdentityOutput{
		Account: aws.String("123456789012"),
	}, nil
}

// --- Mock CloudWatch Logs ---

type mockCWL struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	streams []*cloudwatchlogs.LogStream
	events  []*cloudwatchlogs.OutputLogEvent
}

func (m *mockCWL) DescribeLogStreams(*cloudwatchlogs.DescribeLogStreamsInput) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	return &cloudwatchlogs.DescribeLogStreamsOutput{LogStreams: m.streams}, nil
}

func (m *mockCWL) GetLogEvents(*cloudwatchlogs.GetLogEventsInput) (*cloudwatchlogs.GetLogEventsOutput, error) {
	return &cloudwatchlogs.GetLogEventsOutput{Events: m.events}, nil
}

// --- Helper ---

func newTestManager(ecsMock *mockECS, ebMock *mockEventBridge, ssmMock *mockSSM, stsMock *mockSTS, cwlMock *mockCWL) *Manager {
	project := &config.Project{
		Env:            "test",
		Namespace:      "ns",
		AwsRegion:      "us-east-1",
		DockerRegistry: "123456789012.dkr.ecr.us-east-1.amazonaws.com",
		AWSClient: config.NewAWSClient(
			config.WithECSClient(ecsMock),
			config.WithEventBridgeClient(ebMock),
			config.WithSSMClient(ssmMock),
			config.WithSTSClient(stsMock),
			config.WithCloudWatchLogsClient(cwlMock),
		),
	}

	return &Manager{
		Project: project,
		App: &config.EcsCron{
			Name:     "my-job",
			Schedule: "rate(1 hour)",
		},
	}
}

func baseTD() (*ecs.ListTaskDefinitionsOutput, *ecs.DescribeTaskDefinitionOutput, *ecs.RegisterTaskDefinitionOutput) {
	arn := "arn:aws:ecs:us-east-1:123456789012:task-definition/test-my-job:1"
	family := "test-my-job"
	execRole := "arn:aws:iam::123456789012:role/exec-role"
	containerName := "my-job"
	image := "old-image:latest"

	listOut := &ecs.ListTaskDefinitionsOutput{
		TaskDefinitionArns: []*string{&arn},
	}
	descOut := &ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecs.TaskDefinition{
			TaskDefinitionArn: &arn,
			Family:            &family,
			Revision:          aws.Int64(1),
			ExecutionRoleArn:  &execRole,
			ContainerDefinitions: []*ecs.ContainerDefinition{
				{Name: &containerName, Image: &image},
			},
		},
	}
	regOut := &ecs.RegisterTaskDefinitionOutput{
		TaskDefinition: &ecs.TaskDefinition{
			TaskDefinitionArn: aws.String(arn + "2"),
			Family:            &family,
			Revision:          aws.Int64(2),
		},
	}
	return listOut, descOut, regOut
}

// --- Tests ---

func TestDeploy_WithSchedule(t *testing.T) {
	listOut, descOut, regOut := baseTD()

	m := newTestManager(
		&mockECS{listTDOut: listOut, describeTDOut: descOut, registerTDOut: regOut},
		&mockEventBridge{},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.Deploy(ui)
	if err != nil {
		t.Fatalf("Deploy() error = %v", err)
	}
}

func TestDeploy_WithoutSchedule_RuleExists(t *testing.T) {
	listOut, descOut, regOut := baseTD()

	m := newTestManager(
		&mockECS{listTDOut: listOut, describeTDOut: descOut, registerTDOut: regOut},
		&mockEventBridge{
			describeRuleOut: &eventbridge.DescribeRuleOutput{
				ScheduleExpression: aws.String("rate(1 hour)"),
			},
		},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)
	m.App.Schedule = "" // no schedule

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.Deploy(ui)
	if err != nil {
		t.Fatalf("Deploy() without schedule (rule exists) error = %v", err)
	}
}

func TestDeploy_WithoutSchedule_RuleNotFound(t *testing.T) {
	listOut, descOut, regOut := baseTD()

	m := newTestManager(
		&mockECS{listTDOut: listOut, describeTDOut: descOut, registerTDOut: regOut},
		&mockEventBridge{
			describeRuleErr: fmt.Errorf("ResourceNotFoundException"),
		},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)
	m.App.Schedule = ""

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.Deploy(ui)
	if err == nil {
		t.Fatal("Deploy() without schedule and no rule should fail")
	}
}

func TestDestroy(t *testing.T) {
	m := newTestManager(
		&mockECS{
			listTDOut: &ecs.ListTaskDefinitionsOutput{
				TaskDefinitionArns: []*string{aws.String("arn:1")},
			},
		},
		&mockEventBridge{
			listTargetsOut: &eventbridge.ListTargetsByRuleOutput{
				Targets: []*eventbridge.Target{
					{Id: aws.String("my-job")},
				},
			},
		},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.Destroy(ui, true)
	if err != nil {
		t.Fatalf("Destroy() error = %v", err)
	}
}

func TestDestroy_RuleNotFound(t *testing.T) {
	m := newTestManager(
		&mockECS{
			listTDOut: &ecs.ListTaskDefinitionsOutput{},
		},
		&mockEventBridge{
			listTargetsErr: fmt.Errorf("ResourceNotFoundException"),
			deleteRuleErr:  fmt.Errorf("ResourceNotFoundException"),
		},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.Destroy(ui, true)
	if err != nil {
		t.Fatalf("Destroy() with missing rule should succeed, got error = %v", err)
	}
}

func TestRunNow(t *testing.T) {
	m := newTestManager(
		&mockECS{
			runTaskOut: &ecs.RunTaskOutput{
				Tasks: []*ecs.Task{
					{TaskArn: aws.String("arn:aws:ecs:us-east-1:123456789012:task/test-ns/abc123")},
				},
			},
		},
		&mockEventBridge{},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.RunNow(ui)
	if err != nil {
		t.Fatalf("RunNow() error = %v", err)
	}
}

func TestGetLastRunLogs_WithLogs(t *testing.T) {
	streamName := "main/my-job/abc123"
	m := newTestManager(
		&mockECS{},
		&mockEventBridge{},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{
			streams: []*cloudwatchlogs.LogStream{
				{LogStreamName: &streamName},
			},
			events: []*cloudwatchlogs.OutputLogEvent{
				{Message: aws.String("hello from cron")},
			},
		},
	)

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.GetLastRunLogs(ui)
	if err != nil {
		t.Fatalf("GetLastRunLogs() error = %v", err)
	}
}

func TestGetLastRunLogs_NoStreams(t *testing.T) {
	m := newTestManager(
		&mockECS{},
		&mockEventBridge{},
		&mockSSM{},
		&mockSTS{},
		&mockCWL{},
	)

	ui := terminal.ConsoleUI(context.Background(), true)
	err := m.GetLastRunLogs(ui)
	if err != nil {
		t.Fatalf("GetLastRunLogs() with no streams should succeed, got error = %v", err)
	}
}
