package commands

import (
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_writeConfig(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		exists  map[string]string
		wantErr bool
	}{
		{name: "success", path: "/tmp/ize.toml", exists: map[string]string{"namespace": "test"}, wantErr: false},
		{name: "invalid path", path: "/invalid/path/ize.toml", exists: map[string]string{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeConfig(tt.path, tt.exists); (err != nil) != tt.wantErr {
				t.Errorf("writeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		env     map[string]string
	}{
		{
			name:    "success (only env)",
			args:    []string{"init"},
			env:     map[string]string{"ENV": "test"},
			wantErr: false,
		},
		{
			name:    "success (env and flag)",
			args:    []string{"init", "--skip-examples"},
			env:     map[string]string{"ENV": "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
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

			t.Setenv("HOME", temp)

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

			file, err := os.ReadFile(filepath.Join(".ize", "env", os.Getenv("ENV"), "ize.toml"))
			if err != nil {
				t.Error(err)
			}

			if !strings.Contains(string(file), filepath.Base(temp)) {
				t.Errorf("ize.toml = %v, want contains %s", string(file), filepath.Base(temp))
			}
			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}

func TestInitInternal(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		env     map[string]string
		want    string
	}{
		{
			name:    "success list",
			args:    []string{"init", "--list"},
			wantErr: false,
		},
		{
			name:    "success",
			want:    "\"examples.ize.sh\"",
			args:    []string{"init", "--template", "ecs-apps-monorepo"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
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

			t.Setenv("HOME", temp)

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

			if tt.want != "" {
				file, err := os.ReadFile(filepath.Join(".ize", "env", "testnut", "ize.toml"))
				if err != nil {
					t.Error(err)
				}

				if !strings.Contains(string(file), tt.want) {
					t.Errorf("ize.toml = %v, want contains %s", string(file), tt.want)
				}
			}

			// Unset env
			for k, _ := range tt.env {
				os.Unsetenv(k)
			}
		})
	}
}
