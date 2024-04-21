package commands

import (
	_ "github.com/golang/mock/mockgen/model"
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAWSProfile(t *testing.T) {
	tests := []struct {
		name               string
		args               []string
		wantErr            bool
		wantAWSCredentials string
		env                map[string]string
	}{
		{
			name:               "success (only env)",
			args:               []string{"gen", "aws-profile"},
			env:                map[string]string{"AWS_ACCESS_KEY_ID": "test", "AWS_SECRET_ACCESS_KEY": "test", "AWS_REGION": "us-east-1", "AWS_PROFILE": "test"},
			wantAWSCredentials: "[test]\naws_access_key_id = test\naws_secret_access_key = test\nregion = us-east-1\n\n",
			wantErr:            false,
		},
		{
			name:               "success (only flags)",
			args:               []string{"-r=us-east-2", "-p=test", "gen", "aws-profile"},
			env:                map[string]string{"AWS_ACCESS_KEY_ID": "test", "AWS_SECRET_ACCESS_KEY": "test"},
			wantAWSCredentials: "[test]\naws_access_key_id = test\naws_secret_access_key = test\nregion = us-east-2\n\n",
			wantErr:            false,
		},
		{
			name:               "success (flags and env)",
			args:               []string{"-r=us-west-2", "gen", "aws-profile"},
			env:                map[string]string{"AWS_ACCESS_KEY_ID": "test", "AWS_SECRET_ACCESS_KEY": "test", "AWS_PROFILE": "testnut"},
			wantAWSCredentials: "[testnut]\naws_access_key_id = test\naws_secret_access_key = test\nregion = us-west-2\n\n",
			wantErr:            false,
		},
		{
			name:    "failed (missing aws creds)",
			args:    []string{"gen", "aws-profile"},
			env:     map[string]string{"ENV": "test", "AWS_PROFILE": "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			// Set env
			for k, v := range tt.env {
				os.Setenv(k, v)
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

			os.Setenv("HOME", temp)

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

			err = cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("ize gen tfenv error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			file, err := os.ReadFile(filepath.Join(temp, ".aws", "credentials"))
			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(string(file), tt.wantAWSCredentials) {
				t.Errorf("aws credentials = %v, want %v", string(file), tt.wantAWSCredentials)
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
