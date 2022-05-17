//go:build !e2e
// +build !e2e

package config

import "testing"

func TestCheckCommand(t *testing.T) {
	type args struct {
		command    string
		subcommand []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Error case",
			args: args{
				command:    "error_case",
				subcommand: []string{},
			},
			wantErr: true,
		},
		{
			name: "Success case",
			args: args{
				command:    "ls",
				subcommand: []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckCommand(tt.args.command, tt.args.subcommand)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
