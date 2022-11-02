package commands

import (
	"bytes"
	"fmt"
	"io/fs"
	"os/exec"
	"syscall"

	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var explainTunnelDownTmpl = `
# Change to the dir and send an exit request
(cd {{.EnvDir}} && $(aws ssm get-parameter --name "/{{.Env}}/terraform-output" --with-decryption | jq -r '.Parameter.Value' | base64 -d | jq -r '.cmd.value.tunnel.down'))
`

type TunnelDownOptions struct {
	Config  *config.Project
	Explain bool
}

func NewTunnelDownOptions(project *config.Project) *TunnelDownOptions {
	return &TunnelDownOptions{
		Config: project,
	}
}

func NewCmdTunnelDown(project *config.Project) *cobra.Command {
	o := NewTunnelDownOptions(project)

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Close tunnel",
		Long:  "Close tunnel",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if o.Explain {
				err := o.Config.Generate(explainTunnelDownTmpl, nil)
				if err != nil {
					return err
				}

				return nil
			}

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

	cmd.Flags().BoolVar(&o.Explain, "explain", false, "bash alternative shown")

	return cmd
}

func (o *TunnelDownOptions) Complete() error {
	return nil
}

func (o *TunnelDownOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelDownOptions) Run() error {
	c := exec.Command(
		"ssh", "-S", "bastion.sock", "-O", "exit", "",
	)

	if o.Config.LogLevel == "debug" {
		c.Args = append(c.Args, "-vvv")
	}

	out := &bytes.Buffer{}
	c.Stdout = out
	c.Stderr = out
	c.Dir = o.Config.EnvDir

	err := c.Run()
	if err != nil {
		patherr, ok := err.(*fs.PathError)
		if ok {
			return fmt.Errorf("unable to access folder '%s': %w", c.Dir, patherr.Err)
		}
		exiterr := err.(*exec.ExitError)
		status := exiterr.Sys().(syscall.WaitStatus)
		if status.ExitStatus() != 255 {
			logrus.Debug(out.String())
			return fmt.Errorf("unable to bring the tunnel down: %w", err)
		}
		return fmt.Errorf("unable to bring the tunnel down: tunnel is not active")
	}

	pterm.Success.Println("Tunnel is down!")

	return nil
}
