package tunnel

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewCmdTunnelDown() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "close tunnel",
		Long:  "Close tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := exec.Command(
				"ssh", "-S", "bastion.sock", "-O", "exit", "",
			)
			out := &bytes.Buffer{}
			c.Stdout = out
			c.Stderr = out

			err := c.Run()
			if err != nil {
				exiterr := err.(*exec.ExitError)
				status := exiterr.Sys().(syscall.WaitStatus)
				if status.ExitStatus() != 255 {
					return err
				}
				return fmt.Errorf("can't run tunnel down: %s", out.String())
			}

			pterm.Success.Println("tunnel is down")

			return nil
		},
	}

	return cmd
}
