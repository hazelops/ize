package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var Version = "development"

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show IZE version",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Info.Printfln("Version: %s", GetVersionNumber())
		},
	}

	return cmd
}

func GetVersionNumber() string {
	return Version
}
