//go:build !e2e
// +build !e2e

package requirements

import (
	"fmt"
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
			wantOut: "test\n",
		},
		{name: "failed ssm plugin", args: args{command: "session-manager-plugin"}, want: false, wantOut: ""},
		{name: "success ssm plugin", args: args{command: "session-manager-plugin"}, want: true, wantOut: "session-manager-plugin\n"},
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
