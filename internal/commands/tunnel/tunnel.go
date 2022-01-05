package tunnel

import (
	"github.com/spf13/cobra"
)

func NewCmdTunnel() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "tunnel",
		Short:            "tunnel management",
		Long:             "Tunnel management.",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSSHKey(),
		NewCmdTunnelUp(),
	)

	return cmd
}
