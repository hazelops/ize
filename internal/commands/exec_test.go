package commands

import (
	_ "embed"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/mocks"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:generate mockgen -package=mocks -destination ../../pkg/mocks/mock_ecs.go github.com/aws/aws-sdk-go/service/ecs/ecsiface ECSAPI

//go:embed testdata/build_valid.toml
var execToml string

func TestExec(t *testing.T) {
	mockECS := func(m *mocks.MockECSAPI) {
		m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
			NextToken: nil,
			TaskArns:  []*string{aws.String("test")},
		}, nil).Times(1)
		m.EXPECT().ExecuteCommand(gomock.Any()).Return(&ecs.ExecuteCommandOutput{
			Session: &ecs.Session{
				SessionId:  aws.String("test"),
				StreamUrl:  aws.String("test"),
				TokenValue: aws.String("test"),
			},
		}, nil).Times(1)
	}

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		withConfigFile bool
		env            map[string]string
		mockECSClient  func(m *mocks.MockECSAPI)
	}{
		{
			name:           "success (only config file)",
			args:           []string{"exec", "squibby", "--", "ls"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"-e=test", "-p=test", "exec", "test", "--", "ls"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"-p=test", "exec", "goblin", "--", "ls"},
			env:            map[string]string{"ENV": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
		},
		{
			name:          "success (flags and env)",
			args:          []string{"--aws-region", "us-east-1", "--namespace", "test-testnut", "exec", "goblin", "--", "ls"},
			env:           map[string]string{"ENV": "testnut", "AWS_PROFILE": "test"},
			wantErr:       false,
			mockECSClient: mockECS,
		},
		{
			name:          "success (only flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "exec", "squibby", "--", "ls"},
			wantErr:       false,
			mockECSClient: mockECS,
		},
		{
			name:          "success (only env)",
			args:          []string{"exec", "goblin", "--", "ls"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr:       false,
			mockECSClient: mockECS,
		},
		{
			name:    "failed (list tasks cluster not found)",
			args:    []string{"console", "goblin", "--plain-text"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			mockECSClient: func(m *mocks.MockECSAPI) {
				m.EXPECT().ListTasks(gomock.Any()).Return(nil, awserr.New(ecs.ErrCodeClusterNotFoundException, "", nil)).Times(1)
			},
		},
		{
			name:    "failed (list tasks any err)",
			args:    []string{"console", "goblin", "--plain-text"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			mockECSClient: func(m *mocks.MockECSAPI) {
				m.EXPECT().ListTasks(gomock.Any()).Return(nil, awserr.New("error", "", nil)).Times(1)
			},
		},
		{
			name:    "failed (execute command cluster not found)",
			args:    []string{"console", "goblin", "--plain-text"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			mockECSClient: func(m *mocks.MockECSAPI) {
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().ExecuteCommand(gomock.Any()).Return(nil, awserr.New(ecs.ErrCodeClusterNotFoundException, "", nil)).Times(1)
			},
		},
		{
			name:    "failed (execute command any err)",
			args:    []string{"console", "goblin", "--plain-text"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			mockECSClient: func(m *mocks.MockECSAPI) {
				m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
					NextToken: nil,
					TaskArns:  []*string{aws.String("test")},
				}, nil).Times(1)
				m.EXPECT().ExecuteCommand(gomock.Any()).Return(nil, awserr.New("", "", nil)).Times(1)
			},
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
			if err != nil {
				t.Errorf("get config error = %v, wantErr %v", err, tt.wantErr)
				os.Exit(1)
			}

			cfg.AWSClient = config.NewAWSClient(
				config.WithECSClient(mockECSAPI),
			)

			cfg.Session = getSession(false)

			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ize build error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
