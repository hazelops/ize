package commands

import (
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdTunnel(project *config.Project) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "tunnel",
		Aliases:          []string{"atun"},
		Short:            "Tunnel management",
		Long:             "Tunnel management",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdTunnelUp(project),
		NewCmdTunnelDown(project),
		NewCmdTunnelStatus(project),
	)

	return cmd
}
