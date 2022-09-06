package commands

import (
	_ "embed"
	"github.com/golang/mock/gomock"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed testdata/tfenv_valid.toml
var terraformToml string

func TestTerraform(t *testing.T) {
	tests := []struct {
		name             string
		args             []string
		wantErr          bool
		env              map[string]string
		withConfigFile   bool
		withoutIzeStruct bool
	}{
		{
			name:           "native success (only config)",
			args:           []string{"terraform", "version"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
		},
		{
			name:           "native success (config and flags)",
			args:           []string{"--terraform-version=1.1.3", "terraform", "version"},
			env:            map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
		},
		{
			name:           "native success (config and env and flags)",
			args:           []string{"terraform", "version"},
			env:            map[string]string{"IZE_TERRAFORM_VERSION": "1.1.7", "ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
		},
		{
			name:           "native success (env and flags)",
			args:           []string{"--terraform-version=1.1.7", "terraform", "version"},
			env:            map[string]string{"IZE_TERRAFORM_VERSION": "1.1.3", "ENV": "test", "AWS_PROFILE": "test"},
			withConfigFile: true,
			wantErr:        false,
		},
		{
			name:           "native success (flags)",
			args:           []string{"-e=test", "-r=us-east-1", "-p=test", "-n=testnut", "--terraform-version=1.1.7", "terraform", "version"},
			withConfigFile: true,
			wantErr:        false,
		},
		{
			name:             "failed (without ize struct)",
			args:             []string{"-e=test", "-r=us-east-1", "-p=test", "-n=testnut", "--terraform-version=1.1.7", "terraform", "version"},
			withoutIzeStruct: true,
			withConfigFile:   true,
			wantErr:          true,
		},
		{
			name:    "failed (without ize config)",
			args:    []string{"-e=test", "-r=us-east-1", "-p=test", "-n=testnut", "--terraform-version=1.1.7", "terraform", "version"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			if !tt.withoutIzeStruct {
				err = os.MkdirAll(filepath.Join(temp, ".ize", "env", "test"), 0777)
				if err != nil {
					t.Error(err)
					return
				}
			}

			if tt.withConfigFile {
				setConfigFile(filepath.Join(temp, "ize.toml"), terraformToml, t)
			}
			t.Setenv("HOME", temp)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

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

			cfg.Session = getSession(false)
			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ize terraform error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
