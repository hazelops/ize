package console

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestRunConsole(t *testing.T) {
	tests := []struct {
		name, expectedErr, expectedOut string
		args                           []string
		flags                          map[string]string
	}{
		{
			name:        "successful run",
			expectedOut: "",
			args:        []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "service not found",
			expectedOut: "Error: ServiceNotFoundException: Service not found.",
			args:        []string{"unknow_service"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "env not set",
			expectedOut: "Error: can't validate: env must be specified",
			args:        []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "namespace not set",
			expectedOut: "Error: can't validate: namespace must be specified",
			args:        []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
			},
		},
		{
			name:        "service not set",
			expectedOut: "Error: can't validate: service name must be specified",
			args:        []string{""},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "profile not set",
			expectedOut: "Error: AWS profile must be specified using flags or config file",
			args:        []string{""},
			flags: map[string]string{
				"aws_region": "us-east-1",
				"env":        "dev",
				"namespace":  "nutcorp",
			},
		},
		{
			name:        "invalid creds",
			expectedOut: `Error: can't run console: failed to create aws session`,
			args:        []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "defaultt",
				"aws_region":  "us-east-1",
				"env":         "dev",
				"namespace":   "nutcorp",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := NewCmdConsole()
			cmd.SilenceUsage = true
			out := &bytes.Buffer{}
			cmd.SetOut(out)
			cmd.SetErr(out)
			viper.Reset()
			setGlobalFlags(test.flags)

			cmd.SetArgs(test.args)
			cmd.Execute()

			if strings.TrimSpace(out.String()) != test.expectedOut {
				t.Fatalf("%s: unexpected output: %s\nexpected: %s", test.name, out.String(), test.expectedOut)
			}
		})
	}
}

func setGlobalFlags(flags map[string]string) {
	for k, v := range flags {
		viper.Set(k, v)
	}
}
