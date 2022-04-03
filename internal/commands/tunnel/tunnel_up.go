package tunnel

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelUpOptions struct {
	Config         *config.Config
	PrivateKeyFile string
	BastionHostID  string
	ForwardHost    []string
}

func NewTunnelUpFlags() *TunnelUpOptions {
	return &TunnelUpOptions{}
}

func NewCmdTunnelUp(ui terminal.UI) *cobra.Command {
	o := NewTunnelUpFlags()

	cmd := &cobra.Command{
		Use:   "up",
		Short: "open tunnel",
		Long:  "Open tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete(ui, cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run(ui, cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.BastionHostID, "bastion-instance-id", "", "set bastion host instance id (i-xxxxxxxxxxxxxxxxx)")
	cmd.Flags().StringSliceVar(&o.ForwardHost, "forward-host", nil, "set forward host for redirect with next format: <remote-host>:<remote-port>. In this case a free local port will be selected automatically.  It's possible to set local manually using <remote-host>:<remote-port>:<local-port>")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")

	return cmd
}

func (o *TunnelUpOptions) Complete(ui terminal.UI, cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig(config.WithSSMPlugin())
	if err != nil {
		return fmt.Errorf("can't configure tunnel: %w", err)
	}

	o.Config = cfg

	isUp, err := checkTunnel(ui)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}
	if isUp {
		os.Exit(0)
	}

	if o.PrivateKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PrivateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", home)
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) != 0 {
		return fmt.Errorf("cat't complete options: --forward-host parameter requires --bastion-instance-id\n")
	}

	if len(o.ForwardHost) == 0 && len(o.BastionHostID) != 0 {
		return fmt.Errorf("cat't complete options: --bastion-instance-id requires --forward-host parameter\n")
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		viper.UnmarshalKey("infra.tunnel.bastion_instance_id", &o.BastionHostID)
		viper.UnmarshalKey("infra.tunnel.forward_host", &o.ForwardHost)
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		bastionHostID, forwardHost, err := writeSSHConfigFromSSM(o.Config.Session, o.Config.Env)
		if err != nil {
			return err
		}

		o.BastionHostID = bastionHostID
		o.ForwardHost = forwardHost
		ui.Output("tunnel forwarding configuration obtained from SSM", terminal.WithSuccessStyle())
	} else {
		err := writeSSHConfigFromConfig(o.ForwardHost)
		if err != nil {
			return err
		}
		ui.Output("tunnel forwarding configuration obtained from the config file", terminal.WithSuccessStyle())
	}

	return nil
}

func (o *TunnelUpOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	for _, h := range o.ForwardHost {
		p, _ := strconv.Atoi(strings.Split(h, ":")[2])
		if err := checkPort(p); err != nil {
			return fmt.Errorf("tunnel forwarding config validation failed: %w\n", err)
		}
	}

	return nil
}

func (o *TunnelUpOptions) Run(ui terminal.UI, cmd *cobra.Command) error {
	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))

	if err := setAWSCredentials(o.Config.Session); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	c := exec.Command(
		"ssh", "-M", "-S", "bastion.sock", "-fNT",
		fmt.Sprintf("ubuntu@%s", o.BastionHostID),
		"-F", sshConfigPath,
		"-i", getPrivateKey(o.PrivateKeyFile),
	)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = viper.GetString("ENV_DIR")
	if err := c.Run(); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	ui.Output("tunnel is up! Forwarded ports:", terminal.WithSuccessStyle())

	var fconfig string
	for _, h := range o.ForwardHost {
		ss := strings.Split(h, ":")
		fconfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", ss[0], ss[1], ss[2])
	}
	ui.Output(fconfig)

	return nil
}
