package tunnel

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var IsNotActive = "tunnel is not active\n"

func NewCmdTunnelDown(ui terminal.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "close tunnel",
		Long:  "Close tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			ui := terminal.ConsoleUI(context.Background())

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
					logrus.Debug(out.String())
					return fmt.Errorf("unable to bring the tunnel down: %w", err)
				}
				return fmt.Errorf("unable to bring the tunnel down: tunnel is not active\n")
			}

			ui.Output("tunnel is down!\n", terminal.WithSuccessStyle())

			return nil
		},
	}

	return cmd
}
