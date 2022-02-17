package tunnel

import (
	"bytes"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelStatusOptions struct {
	Config *config.Config
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
	cfg, err := config.InitializeConfig()
	if err != nil {
		return fmt.Errorf("can't complete options: %w", err)
	}

	o.Config = cfg

	return nil
}

func (o *TunnelStatusOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelStatusOptions) Run(cmd *cobra.Command) error {
	c := exec.Command(
		"ssh", "-S", "bastion.sock", "-O", "check", "",
	)
	out := &bytes.Buffer{}
	c.Stdout = out
	c.Stderr = out
	c.Dir = viper.GetString("ENV_DIR")

	err := c.Run()
	if err != nil {
		exiterr := err.(*exec.ExitError)
		status := exiterr.Sys().(syscall.WaitStatus)
		if status.ExitStatus() != 255 {
			logrus.Debug(out.String())
			return fmt.Errorf("can't get tunnel status: %w", err)
		}
		logrus.Debug(out.String())
		pterm.Info.Printfln("tunnel is down")
		return nil
	}

	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))
	sshConfig, err := getSSHConfig(sshConfigPath)
	if err != nil {
		return fmt.Errorf("can't get tunnel status: %w", err)
	}
	hosts := getHosts(sshConfig)

	pterm.Info.Printfln("tunnel is up. Forwarded ports:")
	for _, h := range hosts {
		pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
	}

	return nil
}
