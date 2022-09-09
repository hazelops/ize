package commands

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/mocks"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

//go:generate mockgen -package=mocks -destination ../../pkg/mocks/mock_ecs.go github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI
//go:generate mockgen -package=mocks -destination ../../pkg/mocks/mock_cwl.go github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface CloudWatchLogsAPI

//go:embed testdata/build_valid.toml
var logsToml string

func TestLogs(t *testing.T) {
	mockECS := func(m *mocks.MockECSAPI) {
		m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
			NextToken: nil,
			TaskArns:  []*string{aws.String("test")},
		}, nil).Times(1)
	}

	mockCWL := func(m *mocks.MockCloudWatchLogsAPI) {
		m.EXPECT().GetLogEvents(gomock.Any()).Return(&cloudwatchlogs.GetLogEventsOutput{
			Events: []*cloudwatchlogs.OutputLogEvent{
				{
					IngestionTime: nil,
					Message:       aws.String("test"),
					Timestamp:     aws.Int64(1),
				},
			},
			NextForwardToken: nil,
		}, nil).AnyTimes()
	}

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		withConfigFile bool
		env            map[string]string
		mockECSClient  func(m *mocks.MockECSAPI)
		mockCWLClient  func(m *mocks.MockCloudWatchLogsAPI)
	}{
		{
			name:           "success (only config file)",
			args:           []string{"logs", "squibby"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
			mockCWLClient:  mockCWL,
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"-e=test", "-p=test", "logs", "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
			mockCWLClient:  mockCWL,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"-p=test", "logs", "goblin"},
			env:            map[string]string{"ENV": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
			mockCWLClient:  mockCWL,
		},
		{
			name:          "success (flags and env)",
			args:          []string{"--aws-region", "us-east-1", "--namespace", "test-testnut", "logs", "goblin"},
			env:           map[string]string{"ENV": "testnut", "AWS_PROFILE": "test"},
			wantErr:       false,
			mockECSClient: mockECS,
			mockCWLClient: mockCWL,
		},
		{
			name:          "success (only flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "logs", "squibby"},
			wantErr:       false,
			mockECSClient: mockECS,
			mockCWLClient: mockCWL,
		},
		{
			name:          "success (only env)",
			args:          []string{"logs", "goblin"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr:       false,
			mockECSClient: mockECS,
			mockCWLClient: mockCWL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetEnv(os.Environ())
			viper.Reset()
			os.Unsetenv("IZE_CONFIG_FILE")
			// Set env
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
			err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test"), 0777)
			if err != nil {
				t.Error(err)
				return
			}

			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), buildToml, t)
			}

			err = os.WriteFile(filepath.Join(temp, "session-manager-plugin"), []byte("#!/bin/bash\necho \"session-manager-plugin\""), 0777)
			if err != nil {
				t.Error(err)
			}
			t.Setenv("PATH", fmt.Sprintf("%s:$PATH", temp))

			t.Setenv("HOME", temp)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockECSAPI := mocks.NewMockECSAPI(ctrl)
			tt.mockECSClient(mockECSAPI)

			mockCWLAPI := mocks.NewMockCloudWatchLogsAPI(ctrl)
			tt.mockCWLClient(mockCWLAPI)

			cfg := new(config.Project)
			cmd := newRootCmd(cfg)

			cmd.SetArgs(tt.args)
			cmd.PersistentFlags().ParseErrorsWhitelist.UnknownFlags = true
			err = cmd.PersistentFlags().Parse(tt.args)
			if err != nil {
				t.Error(err)
				return
			}

			cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
				if len(f.Value.String()) != 0 {
					_ = viper.BindPFlag(strings.ReplaceAll(f.Name, "-", "_"), cmd.PersistentFlags().Lookup(f.Name))
				}
			})

			config.InitConfig()

			err = cfg.GetTestConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("get config error = %v, wantErr %v", err, tt.wantErr)
				os.Exit(1)
			}

			cfg.AWSClient = config.NewAWSClient(
				config.WithECSClient(mockECSAPI),
				config.WithCloudWatchLogsClient(mockCWLAPI),
			)

			cfg.Session = getSession(false)

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			errC := make(chan error)
			go func() {
				errC <- cmd.ExecuteContext(ctx)
			}()
			time.Sleep(time.Second)
			signal.NotifyContext(ctx, os.Interrupt)
			cancel()

			time.Sleep(time.Second)

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
