package tunnel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		return err
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
	process, running, err := daemonRunning(daemonContext(cmd.Context()))
	if err != nil {
		return err
	}
	if !running {
		pterm.Info.Println("no background tunnel found")
		return nil
	}

	logrus.Debugf("tunnel is up(pid: %d)")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		return err
	}

	config, err := getTerraformOutput(sess, o.Config.Env)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return err
	}

	hosts := re.FindAllStringSubmatch(
		strings.Join(config.SSHForwardConfig.Value, "\n"),
		-1,
	)

	logrus.Debug("getting SSH forward config")

	logrus.Debugf("hosts: %s", hosts)

	if len(hosts) == 0 {
		logrus.Error("getting SSH forward config")
		return err
	}

	logrus.Debug("port forwarding config is valid")

	logrus.Debugf("tunnel is up(pid: %s)", process.Pid)

	pterm.Info.Printfln("tunnel is up with following config:")
	for _, h := range hosts {
		pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
	}
	return err
}
