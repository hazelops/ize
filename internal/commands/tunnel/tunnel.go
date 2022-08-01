package tunnel

import (
	"github.com/spf13/cobra"
)

func NewCmdTunnel() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "tunnel",
		Short:            "Tunnel management",
		Long:             "Tunnel management",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdTunnelUp(),
		NewCmdTunnelDown(),
		NewCmdTunnelStatus(),
	)

	return cmd
}
