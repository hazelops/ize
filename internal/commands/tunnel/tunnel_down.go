package tunnel

import (
	"github.com/spf13/cobra"
)

func NewCmdTunnelDown() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "close tunnel",
		Long:  "Close tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return killDaemon(daemonContext(cmd.Context()))
		},
	}

	return cmd
}
