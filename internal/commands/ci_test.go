package commands

import (
	_ "embed"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
	"os"
	"path/filepath"
	"testing"
	"text/template"
)

func TestTemplate_Execute(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "test")
	_, err := generate.GenerateFiles("boilerplate-template", filepath.Join(tempDir, "template"))
	if err != nil {
		t.Error(err)
		return
	}

	type fields struct {
		Path    string
		FuncMap template.FuncMap
		Data    interface{}
	}
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "success", fields: fields{
			Path: filepath.Join(tempDir, "template"),
			FuncMap: map[string]any{
				"env":       func() string { return "dev" },
				"namespace": func() string { return "testnut" },
			},
			Data: config.Project{
				Env:        "dev",
				Namespace:  "testnut",
				AwsProfile: "default",
				AwsRegion:  "us-east-2",
			},
		}, args: args{dir: filepath.Join(tempDir, "target")}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			tmpl := &Template{
				Path:    tt.fields.Path,
				FuncMap: tt.fields.FuncMap,
				Data:    tt.fields.Data,
			}
			if err := tmpl.Execute(tt.args.dir); (err != nil) != tt.wantErr {
				t1.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
