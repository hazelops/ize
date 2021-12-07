package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var Version = "development"

type versionCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newVersionCmd() *versionCmd {
	cc := &versionCmd{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "show IZE version",
		Run: func(cmd *cobra.Command, args []string) {
			pterm.Info.Printfln("Version: %s", GetVersionNumber())
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func GetVersionNumber() string {
	return Version
}
