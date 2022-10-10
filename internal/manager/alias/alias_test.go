package alias

import (
	"github.com/cirruslabs/echelon"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/logs"
	"os"
	"testing"
)

func TestManager_Build(t *testing.T) {
	ui, c := logs.GetLogger(false, false, os.Stdout)
	defer c()
	type fields struct {
		Project *config.Project
		App     *config.Alias
	}
	type args struct {
		ui *echelon.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: &config.Project{},
			App: &config.Alias{
				Name:      "test",
				Icon:      "!",
				DependsOn: nil,
			},
		}, args: args{ui: ui}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := a.Build(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Deploy(t *testing.T) {
	ui, c := logs.GetLogger(false, false, os.Stdout)
	defer c()
	type fields struct {
		Project *config.Project
		App     *config.Alias
	}
	type args struct {
		ui *echelon.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: &config.Project{},
			App: &config.Alias{
				Name:      "test",
				Icon:      "!",
				DependsOn: nil,
			},
		}, args: args{ui: ui}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := a.Deploy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Deploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Destroy(t *testing.T) {
	ui, c := logs.GetLogger(false, false, os.Stdout)
	defer c()
	type fields struct {
		Project *config.Project
		App     *config.Alias
	}
	type args struct {
		ui *echelon.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: &config.Project{},
			App: &config.Alias{
				Name:      "test",
				Icon:      "!",
				DependsOn: nil,
			},
		}, args: args{ui: ui}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := a.Destroy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Destroy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Push(t *testing.T) {
	ui, c := logs.GetLogger(false, false, os.Stdout)
	defer c()
	type fields struct {
		Project *config.Project
		App     *config.Alias
	}
	type args struct {
		ui *echelon.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: &config.Project{},
			App: &config.Alias{
				Name:      "test",
				Icon:      "!",
				DependsOn: nil,
			},
		}, args: args{ui: ui}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := a.Push(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Redeploy(t *testing.T) {
	ui, c := logs.GetLogger(false, false, os.Stdout)
	defer c()
	type fields struct {
		Project *config.Project
		App     *config.Alias
	}
	type args struct {
		ui *echelon.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Project: &config.Project{},
			App: &config.Alias{
				Name:      "test",
				Icon:      "!",
				DependsOn: nil,
			},
		}, args: args{ui: ui}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Manager{
				Project: tt.fields.Project,
				App:     tt.fields.App,
			}
			if err := a.Redeploy(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("Redeploy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
