package tunnel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/elliotchance/sshtunnel"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		logrus.Error("getting AWS session")
		return err
	}

	logrus.Debug("getting AWS session")

	to, err := getTerraformOutput(sess, o.Config.Env)
	if err != nil {
		logrus.Error("get forward config")
		return err
	}

	hosts := getForwardConfig(to)

	logrus.Debug("getting SSH forward config")

	logrus.Debugf("hosts: %s", hosts)

	if len(hosts) == 0 {
		return fmt.Errorf("can't tunnel up: forward config is not valid")
	}

	logrus.Debug("forwarding config is valid")

	_ = killDaemon(daemonContext(ctx))

	p, err := daemonContext(ctx).Reborn()
	if err != nil {
		return fmt.Errorf("restarted tunnel process: %w", err)
	}
	if p != nil {
		pterm.Info.Printf("tunnel is up(pid: %d)\n", p.Pid)
		return nil
	}
	defer daemonContext(ctx).Release()

	localport, sessionID, err := startPortForwardSession(to, o.Config.AwsRegion, sess)
	if err != nil {
		return fmt.Errorf("can't tunnel up: %w", err)
	}

	logrus.Debugf("private key path: %s", o.PrivateKeyFile)

	for _, h := range hosts {
		destinationHost := h[2] + ":" + h[3]

		localPort := h[1]

		tunnel := sshtunnel.NewSSHTunnel(
			"ubuntu@localhost",
			sshtunnel.PrivateKeyFile(getPrivateKey(o.PrivateKeyFile)),
			destinationHost,
			localPort,
		)

		tunnel.Server.Port = localport

		go func() {
			if err := tunnel.Start(); err != nil {
				pterm.Error.Printfln("Forward destination host to localhost")
				os.Exit(1)
			}
		}()
		pterm.Info.Printfln("%s => localhost:%s", destinationHost, localPort)
		time.Sleep(100 * time.Millisecond)
	}

	for {
		select {
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGTERM:
				ssm.New(sess).TerminateSession(&ssm.TerminateSessionInput{
					SessionId: &sessionID,
				})
				cancel()
				return nil
			}
		}
	}
}
