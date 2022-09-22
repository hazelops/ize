package serverless

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/terminal"
	"os"
	"path/filepath"
	"testing"
)

func TestManager_Build(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Serverless
	}
	type args struct {
		ui terminal.UI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: new(config.Project),
			App: &config.Serverless{
				Name: "test",
			},
		}, args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sls := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := sls.Build(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Deploy(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Serverless
	}
	type args struct {
		ui terminal.UI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		env     map[string]string
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: new(config.Project),
			App: &config.Serverless{
				Name:         "pecan",
				CreateDomain: true,
			},
		},
			env:  map[string]string{"ENV": "test", "AWS_REGION": "test", "AWS_PROFILE": "test", "NAMESPACE": "test"},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
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

			_, err = generate.GenerateFiles("sls-apps-monorepo", temp)
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(filepath.Join(temp, "nvm.sh"), []byte("#!/bin/bash\nfunction nvm() {\n  echo \"nvm\"\n}\nfunction npm() {\n  echo \"npm\"\n}\nfunction npx() {\n  echo \"npx\"\n}\nexport -f npm\nexport -f nvm\nexport -f npx"), 0777)
			if err != nil {
				t.Error(err)
			}
			t.Setenv("NVM_DIR", temp)

			dir, err := os.ReadDir(filepath.Join(temp, "apps"))
			if err != nil {
				return
			}

			for i, entry := range dir {
				fmt.Println(i, entry.Name())
			}

			config.InitConfig()
			err = tt.fields.Project.GetTestConfig()
			if err != nil {
				t.Error(err)
			}

			sls := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := sls.Deploy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Deploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Destroy(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Serverless
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
		{name: "success", fields: fields{
			Project: new(config.Project),
			App: &config.Serverless{
				Name:         "pecan",
				CreateDomain: true,
			},
		},
			env:  map[string]string{"ENV": "test", "AWS_REGION": "test", "AWS_PROFILE": "test", "NAMESPACE": "test"},
			args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
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

			_, err = generate.GenerateFiles("sls-apps-monorepo", temp)
			if err != nil {
				t.Error(err)
				return
			}

			err = os.WriteFile(filepath.Join(temp, "nvm.sh"), []byte("#!/bin/bash\nfunction nvm() {\n  echo \"nvm\"\n}\nfunction npm() {\n  echo \"npm\"\n}\nfunction npx() {\n  echo \"npx\"\n}\nexport -f npm\nexport -f nvm\nexport -f npx"), 0777)
			if err != nil {
				t.Error(err)
			}
			t.Setenv("NVM_DIR", temp)

			dir, err := os.ReadDir(filepath.Join(temp, "apps"))
			if err != nil {
				return
			}

			for i, entry := range dir {
				fmt.Println(i, entry.Name())
			}

			config.InitConfig()
			err = tt.fields.Project.GetTestConfig()
			if err != nil {
				t.Error(err)
			}
			t.Run(tt.name, func(t *testing.T) {
				sls := &Manager{
					Project: tt.fields.Project,
					App:     tt.fields.App,
				}
				if err := sls.Destroy(tt.args.ui); (err != nil) != tt.wantErr {
					t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
		})
	}
}

func TestManager_Push(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Serverless
	}
	type args struct {
		ui terminal.UI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: new(config.Project),
			App: &config.Serverless{
				Name: "test",
			},
		}, args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sls := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := sls.Push(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Redeploy(t *testing.T) {
	type fields struct {
		Project *config.Project
		App     *config.Serverless
	}
	type args struct {
		ui terminal.UI
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: new(config.Project),
			App: &config.Serverless{
				Name: "test",
			},
		}, args: args{ui: terminal.ConsoleUI(context.TODO(), true)}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sls := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := sls.Redeploy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Redeploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
