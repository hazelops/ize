package tunnel

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TunnelSSHKeyOptions struct {
	Config        *config.Project
	PublicKeyFile string
}

func NewSSHKeyFlags() *TunnelSSHKeyOptions {
	return &TunnelSSHKeyOptions{}
}

func NewCmdSSHKey() *cobra.Command {
	o := NewSSHKeyFlags()

	cmd := &cobra.Command{
		Use:   "ssh-key",
		Short: "Send ssh key to remote server",
		Long:  "Send ssh key to remote server",
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

	cmd.Flags().StringVar(&o.PublicKeyFile, "ssh-public-key", "", "set ssh key public path")

	return cmd
}

func (o *TunnelSSHKeyOptions) Complete() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't load options for a command: %w", err)
	}

	o.Config = cfg

	if o.PublicKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PublicKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	return nil
}

func (o *TunnelSSHKeyOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelSSHKeyOptions) Run() error {
	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	to, err := getTerraformOutput(o.Config.Session, o.Config.Env)
	if err != nil {
		return fmt.Errorf("can't send ssh key: %w", err)
	}

	err = sendSSHPublicKey(to.BastionInstanceID.Value, getPublicKey(o.PublicKeyFile), o.Config.Session)
	if err != nil {
		return fmt.Errorf("can't send ssh key: %w", err)
	}

	pterm.Success.Printfln("SSH user's public key has been sent!\n", terminal.WithSuccessStyle())

	return nil
}
