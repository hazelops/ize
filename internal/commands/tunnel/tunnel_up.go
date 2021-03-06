package tunnel

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelUpOptions struct {
	Config                *config.Project
	PrivateKeyFile        string
	BastionHostID         string
	ForwardHost           []string
	StrictHostKeyChecking bool
}

func NewTunnelUpFlags() *TunnelUpOptions {
	return &TunnelUpOptions{}
}

func NewCmdTunnelUp() *cobra.Command {
	o := NewTunnelUpFlags()

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Open tunnel",
		Long:  "Open tunnel",
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

	cmd.Flags().StringVar(&o.BastionHostID, "bastion-instance-id", "", "set bastion host instance id (i-xxxxxxxxxxxxxxxxx)")
	cmd.Flags().StringSliceVar(&o.ForwardHost, "forward-host", nil, "set forward host for redirect with next format: <remote-host>:<remote-port>. In this case a free local port will be selected automatically.  It's possible to set local manually using <remote-host>:<remote-port>:<local-port>")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")
	cmd.PersistentFlags().BoolVar(&o.StrictHostKeyChecking, "strict-host-key-checking", true, "set strict host key checking")

	return cmd
}

func (o *TunnelUpOptions) Complete() error {
	if err := config.CheckRequirements(config.WithSSMPlugin()); err != nil {
		return err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't configure tunnel: %w", err)
	}

	o.Config = cfg

	isUp, err := checkTunnel(o.Config.EnvDir)
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
		return fmt.Errorf("can't load options for a command: --forward-host parameter requires --bastion-instance-id")
	}

	if len(o.ForwardHost) == 0 && len(o.BastionHostID) != 0 {
		return fmt.Errorf("can't load options for a command: --bastion-instance-id requires --forward-host parameter")
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		viper.UnmarshalKey("infra.tunnel.bastion_instance_id", &o.BastionHostID)
		viper.UnmarshalKey("infra.tunnel.forward_host", &o.ForwardHost)
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		bastionHostID, forwardHost, err := writeSSHConfigFromSSM(o.Config.Session, o.Config.Env, o.Config.EnvDir)
		if err != nil {
			return err
		}

		o.BastionHostID = bastionHostID
		o.ForwardHost = forwardHost
		pterm.Success.Println("Tunnel forwarding configuration obtained from SSM")
	} else {
		err := writeSSHConfigFromConfig(o.ForwardHost, o.Config.EnvDir)
		if err != nil {
			return err
		}
		pterm.Success.Println("Tunnel forwarding configuration obtained from the config file")
	}

	return nil
}

func (o *TunnelUpOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	for _, h := range o.ForwardHost {
		p, _ := strconv.Atoi(strings.Split(h, ":")[2])
		if err := checkPort(p); err != nil {
			return fmt.Errorf("tunnel forwarding config validation failed: %w", err)
		}
	}

	return nil
}

func (o *TunnelUpOptions) Run() error {
	sshConfigPath := fmt.Sprintf("%s/ssh.config", o.Config.EnvDir)

	if err := setAWSCredentials(o.Config.Session); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	args := []string{"-M", "-t", "-S", "bastion.sock", "-fN",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("ubuntu@%s", o.BastionHostID),
		"-F", sshConfigPath}

	if _, err := os.Stat(o.PrivateKeyFile); !os.IsNotExist(err) {
		args = append(args, "-i", o.PrivateKeyFile)
	}

	c := exec.Command("ssh", args...)

	c.Dir = o.Config.EnvDir

	_, _, code, err := runCommand(c)
	if err != nil {
		return err
	}

	if code != 0 {
		return fmt.Errorf("exit status: %d", code)
	}

	pterm.Success.Println("Tunnel is up! Forwarded ports:")

	var fconfig string
	for _, h := range o.ForwardHost {
		ss := strings.Split(h, ":")
		fconfig += fmt.Sprintf("%s:%s ??? localhost:%s\n", ss[0], ss[1], ss[2])
	}
	pterm.Println(fconfig)

	return nil
}
