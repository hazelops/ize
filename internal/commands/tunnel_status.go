package commands

import (
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/cobra"
)

type TunnelStatusOptions struct {
	Config *config.Project
}

func NewTunnelStatusOptions(project *config.Project) *TunnelStatusOptions {
	return &TunnelStatusOptions{
		Config: project,
	}
}

func NewCmdTunnelStatus(project *config.Project) *cobra.Command {
	o := NewTunnelStatusOptions(project)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Tunnel status",
		Long:  "Tunnel running status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete()
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *TunnelStatusOptions) Complete() error {
	return nil
}

func (o *TunnelStatusOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelStatusOptions) Run() error {
	isUp, err := checkTunnel(o.Config.EnvDir)
	if err != nil {
		return fmt.Errorf("can't get tunnel status: %w", err)
	}

	if !isUp {
		return fmt.Errorf("can't get tunnel status: tunnel is down")
	}

	return nil
}
