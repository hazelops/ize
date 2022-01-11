package tunnel

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pterm/pterm"
	"github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCmdTunnel() *cobra.Command {

	cmd := &cobra.Command{
		Use:              "tunnel",
		Short:            "tunnel management",
		Long:             "Tunnel management.",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSSHKey(),
		NewCmdTunnelUp(),
		NewCmdTunnelDown(),
		NewCmdTunnelStatus(),
	)

	return cmd
}

func daemonContext(c context.Context) *daemon.Context {
	return &daemon.Context{
		PidFileName: "tunnel.pid",
		PidFilePerm: 0644,
		LogFileName: "tunnel.log",
		LogFilePerm: 0640,
	}
}

func tunnelDown(dCtx *daemon.Context) error {
	p, err := dCtx.Search()
	if err != nil {
		return fmt.Errorf("search for daemon process: %w", err)
	}
	pterm.Info.Printf("killing daemon process(pid: %d)\n", p.Pid)

	if err := p.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("kill daemon process(pid: %d): %w", p.Pid, err)
	}
	return os.Remove(dCtx.PidFileName)
}

func daemonRunning(dCtx *daemon.Context) (process *os.Process, running bool, err error) {
	p, err := dCtx.Search()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("search daemon process: %w", err)
	}
	err = p.Signal(syscall.Signal(0))
	if err != nil {
		return p, false, nil
	}
	return p, true, nil
}

func getForwardConfig(sess *session.Session, env string) (terraformOutput, error) {
	resp, err := ssm.New(sess).GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		logrus.Error("getting SSH forward config")
		return terraformOutput{}, err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return terraformOutput{}, err
	}

	var config terraformOutput

	err = json.Unmarshal(value, &config)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return terraformOutput{}, err
	}

	logrus.Debugf("output: %s", config)

	return config, nil
}

type terraformOutput struct {
	BastionInstanceID struct {
		Value string `json:"value,omitempty"`
	} `json:"bastion_instance_id,omitempty"`
	SSHForwardConfig struct {
		Value []string `json:"value,omitempty"`
	} `json:"ssh_forward_config,omitempty"`
}
