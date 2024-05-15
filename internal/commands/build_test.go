package commands

import (
	_ "embed"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/mocks"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed testdata/build_valid.toml
var buildToml string

func TestBuild(t *testing.T) {

	mockECS := func(m *mocks.MockECSAPI) {
		m.EXPECT().ListTasks(gomock.Any()).Return(&ecs.ListTasksOutput{
			NextToken: nil,
			TaskArns:  []*string{aws.String("test")},
		}, nil).AnyTimes()
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
			args:           []string{"build", "pecan"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
		},
		{
			name:           "success (flags and config file)",
			args:           []string{"-e=test", "-p=test", "build", "test"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
		},
		{
			name:           "success (flags, env and config file)",
			args:           []string{"-p=test", "build", "goblin"},
			env:            map[string]string{"ENV": "testnut"},
			withConfigFile: true,
			wantErr:        false,
			mockECSClient:  mockECS,
		},
		{
			name:          "success (flags and env)",
			args:          []string{"--aws-region", "us-east-1", "--namespace", "test-testnut", "build", "goblin"},
			env:           map[string]string{"ENV": "testnut", "AWS_PROFILE": "test"},
			wantErr:       false,
			mockECSClient: mockECS,
		},
		{
			name:          "success (only flags)",
			args:          []string{"-e=test", "-r=us-east-1", "-p=test", "-n=test", "build", "squibby"},
			wantErr:       false,
			mockECSClient: mockECS,
		},
		{
			name:          "success (only env)",
			args:          []string{"build", "goblin"},
			env:           map[string]string{"ENV": "test", "AWS_PROFILE": "test", "NAMESPACE": "dev-testnut", "AWS_REGION": "us-west-2"},
			wantErr:       false,
			mockECSClient: mockECS,
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

			err = initTestGitRepo(t, temp)
			if err != nil {
				t.Error(err)
				return
			}

			err = os.Chdir(temp)
			if err != nil {
				t.Error(err)
				return
			}

			//err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test"), 0777)
			//if err != nil {
			//	t.Error(err)
			//	return
			//}

			// Test with config file specified as a parameter
			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), buildToml, t)
			}

			t.Setenv("HOME", temp)

			_, err = generate.GenerateFiles("ecs-apps-monorepo", temp)
			if err != nil {
				t.Error(err)
				return
			}

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
			if (err != nil) != tt.wantErr {
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

func initTestGitRepo(t *testing.T, dir string) error {
	cmd := exec.Command("git", "init")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err != nil {
			t.Errorf("Can't find a directory: %s", err)
		}
	}
	cmd.Dir = dir
	err := cmd.Run()

	if err != nil {
		t.Errorf("Failed to initialize git repository: %s", err)
	}

	cmd = exec.Command("git", "commit -am 'test initial'")
	if err != nil {
		t.Errorf("Failed to commit to git repository: %s", err)
	}

	return err
}
