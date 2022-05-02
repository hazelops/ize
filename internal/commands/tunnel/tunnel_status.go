package tunnel

import (
	"context"
	"fmt"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type TunnelStatusOptions struct {
	Config *config.Config
	UI     terminal.UI
}

func NewTunnelStatusOptions() *TunnelStatusOptions {
	return &TunnelStatusOptions{}
}

func NewCmdTunnelStatus() *cobra.Command {
	o := NewTunnelStatusOptions()

	cmd := &cobra.Command{
		Use:   "status",
		Short: "status tunnel",
		Long:  "Status tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run(cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *TunnelStatusOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't complete options: %w", err)
	}

	o.Config = cfg
	o.UI = terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)

	return nil
}

func (o *TunnelStatusOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelStatusOptions) Run(cmd *cobra.Command) error {
	ui := o.UI
	sg := ui.StepGroup()
	defer sg.Wait()

	isUp, err := checkTunnel()
	if err != nil {
		return fmt.Errorf("can't get tunnel status: %w", err)
	}

	if !isUp {
		return fmt.Errorf("can't get tunnel status: tunnel is down")
	}

	return nil
}
