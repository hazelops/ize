package commands

import "testing"

func TestCheckCommand(t *testing.T) {
	type args struct {
		command    string
		subcommand []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error case",
			args: args{
				command: "error_case",
			},
			wantErr: true,
		},
		{
			name: "Success case",
			args: args{
				command: "ls",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckCommand(tt.args.command, tt.args.subcommand); (err != nil) != tt.wantErr {
				t.Errorf("CheckCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
