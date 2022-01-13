package tunnel

import (
	"fmt"
	"os"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelSSHKeyOptions struct {
	Env           string
	Region        string
	Profile       string
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
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Env = viper.GetString("env")
	o.Profile = viper.GetString("aws-profile")
	o.Region = viper.GetString("aws-region")

	if o.Profile == "" {
		o.Profile = viper.GetString("aws_profile")
	}

	if o.Region == "" {
		o.Region = viper.GetString("aws_region")
	}

	if o.PublicKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PublicKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	return nil
}

func (o *TunnelSSHKeyOptions) Validate() error {
	if len(o.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("AWS region must be specified")
	}
	return nil
}

func (o *TunnelSSHKeyOptions) Run() error {
	pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
	})
	if err != nil {
		pterm.Error.Printfln("Getting AWS session")
		logrus.Error("getting AWS session")
		return err
	}

	logrus.Debug("getting AWS session")

	to, err := getTerraformOutput(sess, o.Env)
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
