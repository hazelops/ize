//go:build !e2e
// +build !e2e

package requirements

import (
	"bytes"
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckCommand(t *testing.T) {
	type args struct {
		command    string
		subcommand []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantOut string
	}{
		{
			name: "failed",
			args: args{
				command:    "error_case",
				subcommand: []string{},
			},
			want: false,
		},
		{
			name: "success echo",
			args: args{
				command:    "echo",
				subcommand: []string{"test"},
			},
			want:    true,
			wantOut: "test",
		},
		{name: "failed ssm plugin", args: args{command: "session-manager-plugin"}, want: false, wantOut: ""},
		{name: "success ssm plugin", args: args{command: "session-manager-plugin"}, want: true, wantOut: "session-manager-plugin"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success ssm plugin" {
				temp, err := os.MkdirTemp("", "test")
				if err != nil {
					t.Fail()
				}
				err = os.WriteFile(filepath.Join(temp, "session-manager-plugin"), []byte("#!/bin/bash\necho \"session-manager-plugin\""), 0777)
				if err != nil {
					t.Error(err)
				}
				err = os.Setenv("PATH", fmt.Sprintf("%s:$PATH", temp))
				if err != nil {
					t.Fail()
				}
			}
			fmt.Println(tt.wantOut)
			exist, out := CheckCommand(tt.args.command, tt.args.subcommand)
			if exist != tt.want {
				t.Errorf("CheckCommand() got = %v, want %v", exist, tt.want)
				return
			}
			if out != tt.wantOut {
				t.Errorf("CheckCommand() got = %v, want %v", out, tt.wantOut)
				return
			}
		})
	}
}

func Test_isStructured(t *testing.T) {
	tests := []struct {
		name        string
		dirName     string
		wantWarning bool
	}{
		{name: "success .infra", dirName: ".infra", wantWarning: false},
		{name: "success .ize", dirName: ".ize", wantWarning: false},
		{name: "with warning", dirName: ".test", wantWarning: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "example")
			if err != nil {
				log.Fatal(err)
			}
			defer os.RemoveAll(dir) // clean up
			projectDir := filepath.Join(dir, tt.dirName)
			err = os.Mkdir(projectDir, 0665)
			if err != nil {
				return
			}

			err = os.Chdir(dir)
			if err != nil {
				return
			}

			buffer := &bytes.Buffer{}

			pterm.SetDefaultOutput(buffer)
			isStructured()

			if buffer.String() == pterm.Warning.Sprint("is not an ize-structured directory. Please run ize init or cd into an ize-structured directory.\n") {
				if !tt.wantWarning {
					t.Fail()
				}
			} else {
				return
			}
		})
	}
}
