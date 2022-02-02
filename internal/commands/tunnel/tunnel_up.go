package tunnel

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelUpOptions struct {
	Config         *config.Config
	PrivateKeyFile string
	BastionHostID  string
	ForwardHost    []string
	sess           *session.Session
}

func NewTunnelUpFlags() *TunnelUpOptions {
	return &TunnelUpOptions{}
}

func NewCmdTunnelUp() *cobra.Command {
	o := NewTunnelUpFlags()

	cmd := &cobra.Command{
		Use:   "up",
		Short: "open tunnel",
		Long:  "Open tunnel.",
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.Flags().StringVar(&o.BastionHostID, "bastion-host-id", "", "set bastion host id")
	cmd.Flags().StringSliceVar(&o.ForwardHost, "forward-host", nil, "set forward host for redirect with next format: host:port:localport")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")

	return cmd
}

func (o *TunnelUpOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig(config.WithSSMPlugin())
	if err != nil {
		return fmt.Errorf("can't complete options: %w", err)
	}

	o.Config = cfg

	isUp, err := checkTunnel()
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}
	if isUp {
		os.Exit(0)
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		return fmt.Errorf("can't complete options: %w", err)
	}

	o.sess = sess

	if o.PrivateKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PrivateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", home)
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) != 0 {
		return fmt.Errorf("cat't complete options: for forward-host should be specified bastion host id")
	}

	if len(o.ForwardHost) == 0 && len(o.BastionHostID) != 0 {
		return fmt.Errorf("cat't complete options: for bastion host id should be specified forward-host")
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		bastionHostID, forwardHost, err := writeSSHConfigFromSSM(o.sess, o.Config.Env)
		if err != nil {
			return err
		}

		o.BastionHostID = bastionHostID
		o.ForwardHost = forwardHost
	} else {
		err := writeSSHConfigFromFlags(o.ForwardHost)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *TunnelUpOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelUpOptions) Run(cmd *cobra.Command) error {
	cmd.SilenceUsage = true

	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))

	c := exec.Command(
		"ssh", "-M", "-S", "bastion.sock", "-fNT",
		fmt.Sprintf("ubuntu@%s", o.BastionHostID),
		"-F", sshConfigPath,
		"-i", getPrivateKey(o.PrivateKeyFile),
	)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	pterm.Success.Printfln("tunnel is up. Forwarded ports:")
	for _, h := range o.ForwardHost {
		ss := strings.Split(h, ":")
		pterm.Info.Printfln("%s:%s ➡ localhost:%s", ss[0], ss[1], ss[2])
	}

	return nil
}
