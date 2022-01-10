package tunnel

import (
	"github.com/pterm/pterm"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

func NewCmdTunnelStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "status tunnel",
		Long:  "Status tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return printDaemonStatus(daemonContext(cmd.Context()))
		},
	}

	return cmd
}

func printDaemonStatus(dCtx *daemon.Context) error {
	process, running, err := daemonRunning(dCtx)
	if err != nil {
		return err
	}
	if !running {
		pterm.Info.Println("no background tunnel found")
		return nil
	}
	pterm.Info.Printf("tunnel is up(pid: %d)\n", process.Pid)
	return err
}
