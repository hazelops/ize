package tunnel

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelUpOptions struct {
	Config         *config.Config
	PrivateKeyFile string
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

	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")

	return cmd
}

func (o *TunnelUpOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	if o.PrivateKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PrivateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", home)
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
	c := exec.Command(
		"ssh", "-S", "bastion.sock", "-O", "check", "",
	)
	out := &bytes.Buffer{}
	c.Stdout = out
	c.Stderr = out

	err := c.Run()
	if err == nil {
		sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))
		sshConfig, err := getSSHConfig(sshConfigPath)
		if err != nil {
			return fmt.Errorf("can't run tunnel up: %w", err)
		}
		hosts := getHosts(sshConfig)
		pterm.Info.Printfln("tunnel is already up. Forwarding config:")
		for _, h := range hosts {
			pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
		}

		return nil
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		logrus.Error("getting AWS session")
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	logrus.Debug("getting AWS session")

	to, err := getTerraformOutput(sess, o.Config.Env)
	if err != nil {
		logrus.Error("get forward config")
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))

	f, err := os.Create(sshConfigPath)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	sshConfig := strings.Join(to.SSHForwardConfig.Value, "\n")
	_, err = io.WriteString(f, sshConfig)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	hosts := getHosts(sshConfig)
	if len(hosts) == 0 {
		return fmt.Errorf("can't tunnel up: forwarding config is not valid")
	}
	logrus.Debugf("hosts: %s", hosts)

	c = exec.Command(
		"ssh", "-M", "-S", "bastion.sock", "-fNT",
		fmt.Sprintf("ubuntu@%s", to.BastionInstanceID.Value),
		"-F", sshConfigPath,
	)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	pterm.Success.Printfln("tunnel is up. Forwarded ports:")
	for _, h := range hosts {
		pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
	}

	return nil
}
