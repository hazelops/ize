package commands

import (
	"github.com/hazelops/ize/internal/config"
	"testing"
	"text/template"
)

func TestTemplate_Execute(t1 *testing.T) {
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
			Path: "./testdata/testrepo",
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
		}, args: args{dir: "./testdata/target"}, wantErr: false},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Template{
				Path:    tt.fields.Path,
				FuncMap: tt.fields.FuncMap,
				Data:    tt.fields.Data,
			}
			if err := t.Execute(tt.args.dir); (err != nil) != tt.wantErr {
				t1.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
