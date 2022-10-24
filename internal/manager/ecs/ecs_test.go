package ecs

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/golang/mock/gomock"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/mocks"
	"github.com/hazelops/ize/pkg/terminal"
)

//go:generate mockgen -package=mocks -destination ../../../pkg/mocks/mock_ecs.go github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI
//go:generate mockgen -package=mocks -destination ../../../pkg/mocks/mock_ecr.go github.com/aws/aws-sdk-go/service/ecr/ecriface ECRAPI
//go:generate mockgen -package=mocks -destination ../../../pkg/mocks/mock_elb.go github.com/aws/aws-sdk-go/service/elbv2/elbv2iface ELBV2API

func TestManager_Build(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Ecs
	}
	type args struct {
		ui terminal.UI
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		env     map[string]string
	}{
		{
			name: "success",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name: "goblin",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			wantErr: false,
		},
		{
			name: "success skip build",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:  "goblin",
					Image: "test",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			wantErr: false,
		},
		{
			name: "invalid path",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name: "goblin",
					Path: "invalid",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			temp, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Error(err)
				return
			}

			err = os.Chdir(temp)
			if err != nil {
				t.Error(err)
				return
			}

			_, err = generate.GenerateFiles("ecs-apps-monorepo", temp)
			if err != nil {
				t.Error(err)
				return
			}

			config.InitConfig()
			err = tt.fields.Project.GetTestConfig()
			if err != nil {
				t.Error(err)
				return
			}

			e := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			e.prepare()
			if err := e.Build(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Deploy(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Ecs
	}
	type args struct {
		ui terminal.UI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		env     map[string]string
		mockCWL func(m *mocks.MockCloudWatchLogsAPI)
		mockECS func(m *mocks.MockECSAPI)
		mockELB func(m *mocks.MockELBV2API)
	}{
		{
			name: "success",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name: "goblin",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().ListTaskDefinitions(gomock.Any()).Return(&ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(2)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).AnyTimes()
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
				m.EXPECT().DeregisterTaskDefinition(gomock.Any()).Return(&ecs.DeregisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{},
				}, nil).Times(1)
				m.EXPECT().RegisterTaskDefinition(gomock.Any()).Return(&ecs.RegisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						Family:            aws.String("test"),
						Revision:          aws.Int64(1),
						TaskDefinitionArn: aws.String("test"),
					},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {},
			wantErr: false,
		},
		{
			name: "success",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:      "goblin",
					DependsOn: nil,
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().ListTaskDefinitions(gomock.Any()).Return(&ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []*string{aws.String("test-arn")},
				}, nil).Times(1)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(1)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
				m.EXPECT().DeregisterTaskDefinition(gomock.Any()).Return(&ecs.DeregisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{},
				}, nil).Times(1)
				m.EXPECT().RegisterTaskDefinition(gomock.Any()).Return(&ecs.RegisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						Family:            aws.String("test"),
						Revision:          aws.Int64(1),
						TaskDefinitionArn: aws.String("test"),
					},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {},
			wantErr: false,
		},
		{
			name: "success unsafe",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:   "goblin",
					Unsafe: true,
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().RegisterTaskDefinition(gomock.Any()).Return(&ecs.RegisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						Family:            aws.String("test"),
						Revision:          aws.Int64(1),
						TaskDefinitionArn: aws.String("test"),
					},
				}, nil).Times(1)
				m.EXPECT().ListTaskDefinitions(gomock.Any()).Return(&ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(2)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
				m.EXPECT().DeregisterTaskDefinition(gomock.Any()).Return(&ecs.DeregisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {
				m.EXPECT().DescribeTargetGroups(gomock.Any()).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{},
					},
				}, nil).Times(1)
				m.EXPECT().ModifyTargetGroup(gomock.Any()).Return(&elbv2.ModifyTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{},
					},
				}, nil).Times(1)
				m.EXPECT().ModifyTargetGroup(gomock.Any()).Return(&elbv2.ModifyTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{},
					},
				}, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failed",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:    "goblin",
					Timeout: 15,
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {
				m.EXPECT().DescribeLogStreams(gomock.Any()).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
					LogStreams: []*cloudwatchlogs.LogStream{
						{
							LogStreamName: aws.String("test"),
						},
					},
				}, nil).Times(1)
				m.EXPECT().GetLogEvents(gomock.Any()).Return(&cloudwatchlogs.GetLogEventsOutput{
					Events: []*cloudwatchlogs.OutputLogEvent{
						{
							Message: aws.String("test"),
						},
					},
				}, nil).Times(1)
			},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().RegisterTaskDefinition(gomock.Any()).Return(&ecs.RegisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						Family:            aws.String("test"),
						Revision:          aws.Int64(1),
						TaskDefinitionArn: aws.String("test"),
					},
				}, nil).Times(1)
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().ListTaskDefinitions(gomock.Any()).Return(&ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(2)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).Times(2)
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(2)
				m.EXPECT().UpdateService(gomock.Any()).Return(nil, awserr.New(ecs.ErrCodeTaskSetNotFoundException, "", nil)).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
				m.EXPECT().DeregisterTaskDefinition(gomock.Any()).Return(&ecs.DeregisterTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			temp, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Error(err)
				return
			}

			err = os.Chdir(temp)
			if err != nil {
				t.Error(err)
				return
			}

			_, err = generate.GenerateFiles("ecs-apps-monorepo", temp)
			if err != nil {
				t.Error(err)
				return
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockECSAPI := mocks.NewMockECSAPI(ctrl)
			mockCWLAPI := mocks.NewMockCloudWatchLogsAPI(ctrl)
			mockELBAPI := mocks.NewMockELBV2API(ctrl)
			tt.mockECS(mockECSAPI)
			tt.mockCWL(mockCWLAPI)
			tt.mockELB(mockELBAPI)

			config.InitConfig()
			tt.fields.Project.AWSClient = config.NewAWSClient(config.WithECSClient(mockECSAPI), config.WithCloudWatchLogsClient(mockCWLAPI), config.WithELBV2Client(mockELBAPI))
			err = tt.fields.Project.GetTestConfig()
			if err != nil {
				t.Error(err)
				return
			}

			e := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			e.prepare()
			if err := e.Deploy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Deploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Redeploy(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Ecs
	}
	type args struct {
		ui terminal.UI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		env     map[string]string
		mockCWL func(m *mocks.MockCloudWatchLogsAPI)
		mockECS func(m *mocks.MockECSAPI)
		mockELB func(m *mocks.MockELBV2API)
	}{
		{
			name: "success",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:                   "goblin",
					TaskDefinitionRevision: "1",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(1)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).AnyTimes()
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {},
			wantErr: false,
		},
		{
			name: "success current",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:                   "goblin",
					TaskDefinitionRevision: "current",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(1)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {},
			wantErr: false,
		},
		{
			name: "success latest",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:                   "goblin",
					TaskDefinitionRevision: "latest",
					Unsafe:                 true,
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(2)
				m.EXPECT().ListTaskDefinitions(gomock.Any()).Return(&ecs.ListTaskDefinitionsOutput{
					TaskDefinitionArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(1)
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					TaskArns: []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().DescribeTasks(gomock.Any()).Return(&ecs.DescribeTasksOutput{
					Tasks: []*ecs.Task{
						{
							LastStatus:        aws.String("RUNNING"),
							TaskDefinitionArn: aws.String("test-arn"),
							Version:           nil,
						},
					},
				}, nil).Times(1)
				m.EXPECT().UpdateService(gomock.Any()).Return(&ecs.UpdateServiceOutput{
					Service: &ecs.Service{
						LoadBalancers: []*ecs.LoadBalancer{
							{
								TargetGroupArn: aws.String("test"),
							},
						},
					},
				}, nil).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {
				m.EXPECT().DescribeTargetGroups(gomock.Any()).Return(&elbv2.DescribeTargetGroupsOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{},
					},
				}, nil).Times(1)
				m.EXPECT().ModifyTargetGroup(gomock.Any()).Return(&elbv2.ModifyTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{},
					},
				}, nil).Times(1)
				m.EXPECT().ModifyTargetGroup(gomock.Any()).Return(&elbv2.ModifyTargetGroupOutput{
					TargetGroups: []*elbv2.TargetGroup{
						{},
					},
				}, nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failed",
			fields: fields{
				Project: new(config.Project),
				App: &config.Ecs{
					Name:                   "goblin",
					TaskDefinitionRevision: "1",
				},
			},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)},
			env: map[string]string{
				"ENV":         "test",
				"AWS_PROFILE": "test",
				"AWS_REGION":  "test",
				"NAMESPACE":   "test",
			},
			mockCWL: func(m *mocks.MockCloudWatchLogsAPI) {
				m.EXPECT().DescribeLogStreams(gomock.Any()).Return(&cloudwatchlogs.DescribeLogStreamsOutput{
					LogStreams: []*cloudwatchlogs.LogStream{
						{
							LogStreamName: aws.String("test"),
						},
					},
				}, nil).Times(1)
				m.EXPECT().GetLogEvents(gomock.Any()).Return(&cloudwatchlogs.GetLogEventsOutput{
					Events: []*cloudwatchlogs.OutputLogEvent{
						{
							Message: aws.String("test"),
						},
					},
				}, nil).Times(1)
			},
			mockECS: func(m *mocks.MockECSAPI) {
				m.EXPECT().DescribeServices(gomock.Any()).Return(&ecs.DescribeServicesOutput{
					Services: []*ecs.Service{
						{
							TaskDefinition: aws.String("test-arn"),
							DesiredCount:   aws.Int64(1),
							Deployments: []*ecs.Deployment{
								{},
							},
						},
					},
				}, nil).Times(1)
				m.EXPECT().DescribeTaskDefinition(gomock.Any()).Return(&ecs.DescribeTaskDefinitionOutput{
					TaskDefinition: &ecs.TaskDefinition{
						ContainerDefinitions: []*ecs.ContainerDefinition{{
							Image: aws.String("test"),
							Name:  aws.String("test"),
						}},
						Family:            aws.String("test-family"),
						Revision:          aws.Int64(1),
						Status:            nil,
						TaskDefinitionArn: aws.String("test-arn"),
						TaskRoleArn:       nil,
						Volumes:           nil,
					},
				}, nil).Times(1)

				m.EXPECT().UpdateService(gomock.Any()).Return(nil, awserr.New(ecs.ErrCodeTaskSetNotFoundException, "", nil)).Times(1)
			},
			mockELB: func(m *mocks.MockELBV2API) {},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			temp, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Error(err)
				return
			}

			err = os.Chdir(temp)
			if err != nil {
				t.Error(err)
				return
			}

			_, err = generate.GenerateFiles("ecs-apps-monorepo", temp)
			if err != nil {
				t.Error(err)
				return
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockECSAPI := mocks.NewMockECSAPI(ctrl)
			mockCWLAPI := mocks.NewMockCloudWatchLogsAPI(ctrl)
			mockELBAPI := mocks.NewMockELBV2API(ctrl)
			tt.mockECS(mockECSAPI)
			tt.mockCWL(mockCWLAPI)
			tt.mockELB(mockELBAPI)

			config.InitConfig()
			tt.fields.Project.AWSClient = config.NewAWSClient(config.WithECSClient(mockECSAPI), config.WithCloudWatchLogsClient(mockCWLAPI), config.WithELBV2Client(mockELBAPI))
			err = tt.fields.Project.GetTestConfig()
			if err != nil {
				t.Error(err)
				return
			}

			e := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			e.prepare()
			if err := e.Redeploy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Deploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Destroy(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Ecs
	}
	type args struct {
		ui terminal.UI
	}

	env := map[string]string{
		"ENV":         "test",
		"AWS_PROFILE": "test",
		"AWS_REGION":  "test",
		"NAMESPACE":   "test",
	}

	for k, v := range env {
		t.Setenv(k, v)
	}

	config.InitConfig()
	project := config.Project{}
	err := project.GetTestConfig()
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: &project,
			App: &config.Ecs{
				Name: "goblin",
			},
		}, args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := e.Destroy(tt.args.ui, true); (err != nil) != tt.wantErr {
				t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func TestManager_Push(t *testing.T) {
//	type fields struct {
//		Project *config.Project
//		App     *config.Ecs
//	}
//	type args struct {
//		ui terminal.UI
//	}
//
//	env := map[string]string{
//		"ENV":         "test",
//		"AWS_PROFILE": "test",
//		"AWS_REGION":  "test",
//		"NAMESPACE":   "test",
//		"IZE_TAG":     "test",
//	}
//
//	for k, v := range env {
//		t.Setenv(k, v)
//	}
//
//	mockECR := func(m *mocks.MockECRAPI) {
//		m.EXPECT().DescribeRepositories(gomock.Any()).Return(&ecr.DescribeRepositoriesOutput{
//			Repositories: []*ecr.Repository{
//				{
//					RepositoryUri: aws.String("test"),
//				},
//			},
//		}, nil).AnyTimes()
//		m.EXPECT().GetAuthorizationToken(gomock.Any()).Return(&ecr.GetAuthorizationTokenOutput{
//			AuthorizationData: []*ecr.AuthorizationData{
//				{
//					AuthorizationToken: aws.String("dG9rZW4="),
//					ExpiresAt:          nil,
//					ProxyEndpoint:      nil,
//				},
//			},
//		}, nil).AnyTimes()
//	}
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockECRAPI := mocks.NewMockECRAPI(ctrl)
//	mockECR(mockECRAPI)
//
//	config.InitConfig()
//	project := config.Project{}
//	project.AWSClient = config.NewAWSClient(config.WithECRClient(mockECRAPI))
//	err := project.GetTestConfig()
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		wantErr bool
//	}{
//		{name: "success", fields: fields{
//			Project: &project,
//			App: &config.Ecs{
//				Name: "goblin",
//			},
//		}, args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			e := &Manager{
//				Project: tt.fields.Project,
//				App:     tt.fields.App,
//			}
//			if err := e.Push(tt.args.ui); (err != nil) != tt.wantErr {
//				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
