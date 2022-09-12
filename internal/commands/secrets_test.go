package commands

import (
	_ "embed"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
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

//go:generate mockgen -package=mocks -destination ../../pkg/mocks/mock_ssm.go github.com/aws/aws-sdk-go/service/ssm/ssmiface SSMAPI

//go:embed testdata/build_valid.toml
var secretsToml string

func TestSecretsPull(t *testing.T) {
	mockSSM := func(m *mocks.MockSSMAPI) {
		m.EXPECT().GetParametersByPath(gomock.Any()).Return(&ssm.GetParametersByPathOutput{
			Parameters: []*ssm.Parameter{
				{
					Name:  aws.String("test"),
					Value: aws.String("test"),
				},
			},
		}, nil).Times(1)
	}

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		withConfigFile bool
		env            map[string]string
		mockSSMClient  func(m *mocks.MockSSMAPI)
		want           string
	}{
		{
			name:           "success (only config file)",
			args:           []string{"secrets", "pull", "squibby"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			want:           "\"test\": \"test\"",
			mockSSMClient:  mockSSM,
		},
		{
			name:    "failed",
			args:    []string{"secrets", "pull", "squibby"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			want:    "",
			mockSSMClient: func(m *mocks.MockSSMAPI) {
				m.EXPECT().GetParametersByPath(gomock.Any()).Return(nil, awserr.New("error", "", nil)).Times(1)
			},
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"-e=test", "-p=test", "secrets", "pull", "squibby"},
			withConfigFile: true,
			wantErr:        false,
			want:           "\"test\": \"test\"",
			mockSSMClient:  mockSSM,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"-p=test", "secrets", "pull", "squibby"},
			env:            map[string]string{"ENV": "test"},
			withConfigFile: true,
			wantErr:        false,
			want:           "\"test\": \"test\"",
			mockSSMClient:  mockSSM,
		},
		{
			name:          "success (flags and env)",
			args:          []string{"--aws-region", "us-east-1", "--namespace", "test-testnut", "secrets", "pull", "squibby"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			wantErr:       false,
			want:          "\"test\": \"test\"",
			mockSSMClient: mockSSM,
		},
		{
			name:          "success (only flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "secrets", "pull", "squibby"},
			wantErr:       false,
			want:          "\"test\": \"test\"",
			mockSSMClient: mockSSM,
		},
		{
			name:          "success (only env)",
			args:          []string{"secrets", "pull", "squibby"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr:       false,
			want:          "\"test\": \"test\"",
			mockSSMClient: mockSSM,
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
			err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test", "secrets"), 0777)
			if err != nil {
				t.Error(err)
				return
			}

			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), buildToml, t)
			}

			t.Setenv("HOME", temp)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSSMAPI := mocks.NewMockSSMAPI(ctrl)
			tt.mockSSMClient(mockSSMAPI)

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
				config.WithSSMClient(mockSSMAPI),
			)

			cfg.Session = getSession(false)

			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ize build error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file, err := os.ReadFile(filepath.Join(temp, ".ize", "env", "test", "secrets", "squibby.json"))
			if (err != nil) != tt.wantErr {
				t.Error(err)
				return
			}

			if !strings.Contains(string(file), tt.want) {
				t.Errorf("output = %v, want %v", string(file), tt.want)
				return
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}

func TestSecretsPush(t *testing.T) {
	mockSSM := func(m *mocks.MockSSMAPI) {
		m.EXPECT().PutParameter(gomock.Any()).Return(nil, nil).Times(2)
		m.EXPECT().AddTagsToResource(gomock.Any()).Return(nil, nil).Times(2)
	}

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		withConfigFile bool
		env            map[string]string
		mockSSMClient  func(m *mocks.MockSSMAPI)
	}{
		{
			name:           "success (only config file)",
			args:           []string{"secrets", "push", "squibby"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockSSMClient:  mockSSM,
		},
		{
			name:    "failed",
			args:    []string{"secrets", "push", "squibby"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			mockSSMClient: func(m *mocks.MockSSMAPI) {
				m.EXPECT().PutParameter(gomock.Any()).Return(nil, awserr.New("error", "", nil)).Times(1)
			},
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"-e=test", "-p=test", "secrets", "push", "squibby"},
			withConfigFile: true,
			wantErr:        false,
			mockSSMClient:  mockSSM,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"-p=test", "secrets", "push", "squibby"},
			env:            map[string]string{"ENV": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockSSMClient:  mockSSM,
		},
		{
			name:          "success (flags and env)",
			args:          []string{"--aws-region", "us-east-1", "--namespace", "test-testnut", "secrets", "push", "squibby"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			wantErr:       false,
			mockSSMClient: mockSSM,
		},
		{
			name:          "success (only flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "secrets", "push", "squibby"},
			wantErr:       false,
			mockSSMClient: mockSSM,
		},
		{
			name:          "success (only env)",
			args:          []string{"secrets", "push", "squibby"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr:       false,
			mockSSMClient: mockSSM,
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
			err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test", "secrets"), 0777)
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(filepath.Join(temp, ".ize", "env", "test", "secrets", "squibby.json"), []byte("{\n  \"service__key__one\": \"value one\",\n  \"service__key__two\": \"test value two\"\n}"), 0666)
			if err != nil {
				t.Error(err)
			}

			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), buildToml, t)
			}

			t.Setenv("HOME", temp)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSSMAPI := mocks.NewMockSSMAPI(ctrl)
			tt.mockSSMClient(mockSSMAPI)

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
				config.WithSSMClient(mockSSMAPI),
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

func TestSecretsRm(t *testing.T) {
	mockSSM := func(m *mocks.MockSSMAPI) {
		m.EXPECT().GetParametersByPath(gomock.Any()).Return(&ssm.GetParametersByPathOutput{
			Parameters: []*ssm.Parameter{
				{
					Name:  aws.String("test"),
					Value: aws.String("test"),
				},
			},
		}, nil).Times(1)
		m.EXPECT().DeleteParameters(gomock.Any()).Return(nil, nil).Times(1)
	}

	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		withConfigFile bool
		env            map[string]string
		mockSSMClient  func(m *mocks.MockSSMAPI)
	}{
		{
			name:           "success (only config file)",
			args:           []string{"secrets", "rm", "squibby"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockSSMClient:  mockSSM,
		},
		{
			name:    "failed",
			args:    []string{"secrets", "rm", "squibby"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr: true,
			mockSSMClient: func(m *mocks.MockSSMAPI) {
				m.EXPECT().GetParametersByPath(gomock.Any()).Return(nil, awserr.New("error", "", nil)).Times(1)
			},
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"-e=test", "-p=test", "secrets", "rm", "squibby"},
			withConfigFile: true,
			wantErr:        false,
			mockSSMClient:  mockSSM,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"-p=test", "secrets", "rm", "squibby"},
			env:            map[string]string{"ENV": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockSSMClient:  mockSSM,
		},
		{
			name:          "success (flags and env)",
			args:          []string{"--aws-region", "us-east-1", "--namespace", "test-testnut", "secrets", "rm", "squibby"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			wantErr:       false,
			mockSSMClient: mockSSM,
		},
		{
			name:          "success (only flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "secrets", "rm", "squibby"},
			wantErr:       false,
			mockSSMClient: mockSSM,
		},
		{
			name:          "success (only env)",
			args:          []string{"secrets", "rm", "squibby"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr:       false,
			mockSSMClient: mockSSM,
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
			err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test", "secrets"), 0777)
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(filepath.Join(temp, ".ize", "env", "test", "secrets", "squibby.json"), []byte("{\n  \"service__key__one\": \"value one\",\n  \"service__key__two\": \"test value two\"\n}"), 0666)
			if err != nil {
				t.Error(err)
			}

			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), buildToml, t)
			}

			t.Setenv("HOME", temp)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSSMAPI := mocks.NewMockSSMAPI(ctrl)
			tt.mockSSMClient(mockSSMAPI)

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
				config.WithSSMClient(mockSSMAPI),
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
