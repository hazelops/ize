package commands

import (
	"github.com/spf13/cobra"
)

type Response struct {
	Err error

	Cmd *cobra.Command
}

func Execute(args []string) Response {
	izeCmd := newCommandBuilder().addAll().build()
	cmd := izeCmd.getCommand()
	cmd.SetArgs(args)

	c, err := cmd.ExecuteC()

	var resp Response

	resp.Err = err
	resp.Cmd = c

	return resp
}
