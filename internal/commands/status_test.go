package commands

import (
	"bytes"
	_ "embed"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/mocks"
	"github.com/pterm/pterm"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:generate mockgen -package=mocks -destination ../mocks/mock_aim.go github.com/aws/aws-sdk-go/service/iam/iamiface IAMAPI
//go:generate mockgen -package=mocks -destination ../mocks/mock_sts.go github.com/aws/aws-sdk-go/service/sts/stsiface STSAPI

//go:embed testdata/status_valid.toml
var statusToml string

func TestStatus(t *testing.T) {
	mockIAMClient := func(m *mocks.MockIAMAPI) {
		m.EXPECT().ListUserTags(gomock.Any()).Return(&iam.ListUserTagsOutput{
			IsTruncated: nil,
			Marker:      nil,
			Tags:        []*iam.Tag{{Key: aws.String("devEnvironmentName"), Value: aws.String("test")}},
		}, nil).Times(1)
		m.EXPECT().GetUser(gomock.Any()).Return(&iam.GetUserOutput{
			User: &iam.User{
				UserName: aws.String("test"),
			},
		}, nil).Times(1)
	}
	mockSTSClient := func(m *mocks.MockSTSAPI) {
		m.EXPECT().GetCallerIdentity(gomock.Any()).Return(&sts.GetCallerIdentityOutput{
			Account: aws.String("0"),
		}, nil).Times(1)
	}

	pterm.DisableStyling()
	tests := []struct {
		name           string
		args           []string
		wantErr        bool
		env            map[string]string
		withConfigFile bool
		contains       []string
		mockIAMClient  func(m *mocks.MockIAMAPI)
		mockSTSClient  func(m *mocks.MockSTSAPI)
	}{
		{
			name:           "success (only config)",
			args:           []string{"status"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockIAMClient:  mockIAMClient,
			mockSTSClient:  mockSTSClient,
			contains: []string{
				"ENV            | test",
				"NAMESPACE      | testnut",
				"TAG            | test",
				"TERRAFORM_VERSION | 1.2.6",
				"AWS PROFILE | test",
				"AWS USER    | test",
				"AWS ACCOUNT | 0",
				"AWS_DEV_ENV_NAME | test",
			},
		},
		{
			name:           "success (config and flags)",
			args:           []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "--terraform-version=1.1.3", "status"},
			withConfigFile: true,
			wantErr:        false,
			mockIAMClient:  mockIAMClient,
			mockSTSClient:  mockSTSClient,
			contains: []string{
				"ENV            | test",
				"NAMESPACE      | test",
				"TAG            | test",
				"TERRAFORM_VERSION | 1.1.3",
				"AWS PROFILE | test",
				"AWS USER    | test",
				"AWS ACCOUNT | 0",
				"AWS_DEV_ENV_NAME | test",
			},
		},
		{
			name:           "success (config and env and flags)",
			args:           []string{"-p=testnut", "-r=us-east-1", "--terraform-version=1.1.5", "status"},
			env:            map[string]string{"IZE_TERRAFORM_VERSION": "1.1.7", "ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockIAMClient:  mockIAMClient,
			mockSTSClient:  mockSTSClient,
			contains: []string{
				"ENV            | test",
				"NAMESPACE      | test",
				"TAG            | test",
				"TERRAFORM_VERSION | 1.1.5",
				"AWS PROFILE | testnut",
				"AWS USER    | test",
				"AWS ACCOUNT | 0",
				"AWS_DEV_ENV_NAME | test",
			},
		},
		{
			name:          "success (env and flags)",
			args:          []string{"--namespace", "test", "-r=us-west-2", "--terraform-version=1.1.7", "status"},
			env:           map[string]string{"IZE_TERRAFORM_VERSION": "1.1.3", "ENV": "test", "AWS_PROFILE": "test"},
			wantErr:       false,
			mockIAMClient: mockIAMClient,
			mockSTSClient: mockSTSClient,
			contains: []string{
				"ENV            | test",
				"NAMESPACE      | test",
				"TAG            | test",
				"TERRAFORM_VERSION | 1.1.7",
				"AWS PROFILE | test",
				"AWS USER    | test",
				"AWS ACCOUNT | 0",
				"AWS_DEV_ENV_NAME | test",
			},
		},
		{
			name:          "success (flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=testnut", "--terraform-version=1.1.7", "status"},
			wantErr:       false,
			mockIAMClient: mockIAMClient,
			mockSTSClient: mockSTSClient,
			contains: []string{
				"ENV            | test",
				"NAMESPACE      | testnut",
				"TAG            | test",
				"TERRAFORM_VERSION | 1.1.7",
				"AWS PROFILE | test",
				"AWS USER    | test",
				"AWS ACCOUNT | 0",
				"AWS_DEV_ENV_NAME | test",
			},
		},
		{
			name:          "success (env)",
			wantErr:       false,
			args:          []string{"status"},
			env:           map[string]string{"IZE_TERRAFORM_VERSION": "1.1.5", "ENV": "test", "AWS_PROFILE": "test", "AWS_REGION": "us-east-1", "NAMESPACE": "dev-test"},
			mockIAMClient: mockIAMClient,
			mockSTSClient: mockSTSClient,
			contains: []string{
				"ENV            | test",
				"NAMESPACE      | dev-test",
				"TAG            | test",
				"TERRAFORM_VERSION | 1.1.5",
				"AWS PROFILE | test",
				"AWS USER    | test",
				"AWS ACCOUNT | 0",
				"AWS_DEV_ENV_NAME | test",
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
				setConfigFile(filepath.Join(temp, "ize.toml"), terraformToml, t)
			}
			t.Setenv("HOME", temp)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockIAMAPI := mocks.NewMockIAMAPI(ctrl)
			tt.mockIAMClient(mockIAMAPI)

			mockSTSAPI := mocks.NewMockSTSAPI(ctrl)
			tt.mockSTSClient(mockSTSAPI)

			cfg := new(config.Project)
			cmd := newRootCmd(cfg)

			cmd.SetArgs(tt.args)
			buf := new(bytes.Buffer)
			pterm.SetDefaultOutput(buf)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
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
				config.WithIAMClient(mockIAMAPI),
				config.WithSTSClient(mockSTSAPI),
			)

			cfg.Session = getSession(false)
			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ize terraform error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, contain := range tt.contains {
				if !strings.Contains(buf.String(), contain) {
					t.Errorf("output = %v, want %v", buf.String(), contain)
					return
				}
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
