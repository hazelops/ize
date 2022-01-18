package tunnel

import (
	"fmt"
	"os"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TunnelSSHKeyOptions struct {
	Config        *config.Config
	PublicKeyFile string
}

func NewSSHKeyFlags() *TunnelSSHKeyOptions {
	return &TunnelSSHKeyOptions{}
}

func NewCmdSSHKey() *cobra.Command {
	o := NewSSHKeyFlags()

	cmd := &cobra.Command{
		Use:   "ssh-key",
		Short: "send ssh key to remote server",
		Long:  "Send ssh key to remote server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
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

func (o *TunnelSSHKeyOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
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
	pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		pterm.Error.Printfln("Getting AWS session")
		logrus.Error("getting AWS session")
		return err
	}

	logrus.Debug("getting AWS session")

	to, err := getTerraformOutput(sess, o.Config.Env)
	if err != nil {
		return fmt.Errorf("can't get forward config: %w", err)
	}

	logrus.Debug("getting bastion instance ID")

	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	err = sendSSHPublicKey(to.BastionInstanceID.Value, getPublicKey(o.PublicKeyFile), sess)

	logrus.Debug("reading user SSH public key")

	if err != nil {
		logrus.Error("sending user SSH public key")
		return err
	}

	logrus.Debug("sending user SSH public key")

	return nil
}
