package terraform

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"io"
	"os"
	"testing"
)

func TestInstall(t *testing.T) {
	type args struct {
		tfversion string
		mirrorURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "success", args: args{
			tfversion: "1.1.9",
			mirrorURL: defaultMirror,
		}, wantErr: false},
		{name: "success exist", args: args{
			tfversion: "1.1.9",
			mirrorURL: defaultMirror,
		}, wantErr: false},
		{name: "invalid mirror", args: args{
			tfversion: "1.2.9",
			mirrorURL: "invalid.com/mirror",
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Install(tt.args.tfversion, tt.args.mirrorURL); (err != nil) != tt.wantErr {
				t.Errorf("Install() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_installVersion(t *testing.T) {
	mirror := defaultMirror
	hd := home()

	type args struct {
		version   string
		mirrorURL *string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "success", args: args{
			version:   "1.1.3",
			mirrorURL: &mirror,
		}, want: fmt.Sprintf("%s/%s", hd, ".ize/versions/terraform/terraform_1.1.3"), wantErr: false},
		{name: "success exist", args: args{
			version:   "1.1.3",
			mirrorURL: &mirror,
		}, want: fmt.Sprintf("%s/%s", hd, ".ize/versions/terraform/terraform_1.1.3"), wantErr: false},
		{name: "invalid tf version", args: args{
			version:   "1.1.X",
			mirrorURL: &mirror,
		}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := installVersion(tt.args.version, tt.args.mirrorURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("installVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("installVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_local_Prepare(t *testing.T) {
	type fields struct {
		version string
		command []string
		env     []string
		output  io.Writer
		tfpath  string
		project *config.Project
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "success", fields: fields{
			version: "1.1.3",
			command: []string{"version"},
			env:     nil,
			output:  os.Stdout,
			tfpath:  "",
			project: &config.Project{},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &local{
				version: tt.fields.version,
				command: tt.fields.command,
				env:     tt.fields.env,
				output:  tt.fields.output,
				tfpath:  tt.fields.tfpath,
				project: tt.fields.project,
			}
			if err := l.Prepare(); (err != nil) != tt.wantErr {
				t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_local_Run(t *testing.T) {
	mirror := defaultMirror
	version, err := installVersion("1.1.3", &mirror)
	if err != nil {
		t.Error(err)
		return
	}

	type fields struct {
		version string
		command []string
		env     []string
		output  io.Writer
		tfpath  string
		project *config.Project
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "success", fields: fields{
			version: "1.1.3",
			command: []string{"version"},
			env:     nil,
			output:  os.Stdout,
			tfpath:  version,
			project: &config.Project{},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &local{
				version: tt.fields.version,
				command: tt.fields.command,
				env:     tt.fields.env,
				output:  tt.fields.output,
				tfpath:  tt.fields.tfpath,
				project: tt.fields.project,
			}
			if err := l.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_local_RunUI(t *testing.T) {
	mirror := defaultMirror
	version, err := installVersion("1.1.3", &mirror)
	if err != nil {
		t.Error(err)
		return
	}

	type fields struct {
		version string
		command []string
		env     []string
		output  io.Writer
		tfpath  string
		project *config.Project
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
		{
			name: "success",
			fields: fields{
				version: "1.1.3",
				command: []string{"version"},
				env:     nil,
				output:  os.Stdout,
				tfpath:  version,
				project: &config.Project{},
			},
			args:    args{ui: terminal.ConsoleUI(context.TODO(), true)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &local{
				version: tt.fields.version,
				command: tt.fields.command,
				env:     tt.fields.env,
				output:  tt.fields.output,
				tfpath:  tt.fields.tfpath,
				project: tt.fields.project,
			}
			if err := l.RunUI(tt.args.ui); (err != nil) != tt.wantErr {
				t.Errorf("RunUI() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func home() string {
	h, _ := os.UserHomeDir()
	return h
}
