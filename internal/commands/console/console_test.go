//go:build !e2e
// +build !e2e

package console

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestRunConsole(t *testing.T) {
	tests := []struct {
		name, expectedErr string
		args              []string
		flags             map[string]string
	}{
		{
			name: "successful run",
			args: []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "app not found",
			expectedErr: "ServiceNotFoundException: Service not found.",
			args:        []string{"unknow_app"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "env not set",
			expectedErr: "can't validate: env must be specified\n",
			args:        []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"namespace":   "nutcorp",
			},
		},
		{
			name:        "namespace not set",
			expectedErr: "can't validate: namespace must be specified\n",
			args:        []string{"squibby"},
			flags: map[string]string{
				"aws_profile": "default",
				"aws_region":  "us-east-1",
				"env":         "dev",
			},
		},
		{
			name:        "app name not set",
			expectedErr: "can't validate: app name must be specified\n",
			args:        []string{""},
			flags: map[string]string{
				"aws_region": "us-east-1",
				"env":        "dev",
				"namespace":  "nutcorp",
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
			err := cmd.Execute()

			if test.expectedErr != "" {
				require.EqualError(t, err, test.expectedErr)
			}
		})
	}
}

func setGlobalFlags(flags map[string]string) {
	for k, v := range flags {
		viper.Set(k, v)
	}
}
