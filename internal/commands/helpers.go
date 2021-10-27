package commands

import "github.com/spf13/cobra"

type cmder interface {
	getCommand() *cobra.Command
}

func CheckCommand(command string, subcommand []string) error {
	err := exec.Command(command, subcommand...).Run()
	if err != nil {
		return err
	}

	return nil
}
