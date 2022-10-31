package commands

import (
	"context"
	"fmt"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

var explainTunnelStatusTmpl = `
# Change to the dir and get status
(cd {{.EnvDir}} && $(aws ssm get-parameter --name "/{{.Env}}/terraform-output" --with-decryption | jq -r '.Parameter.Value' | base64 -d | jq -r '.cmd.value.tunnel.status'))
`

type TunnelStatusOptions struct {
	Config  *config.Project
	UI      terminal.UI
	Explain bool
}

func NewTunnelStatusOptions(project *config.Project) *TunnelStatusOptions {
	return &TunnelStatusOptions{
		Config: project,
	}
}

func NewCmdTunnelStatus(project *config.Project) *cobra.Command {
	o := NewTunnelStatusOptions(project)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Tunnel status",
		Long:  "Tunnel running status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if o.Explain {
				err := o.Config.Generate(explainTunnelStatusTmpl, nil)
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

func (o *TunnelStatusOptions) Complete() error {
	o.UI = terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	return nil
}

func (o *TunnelStatusOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelStatusOptions) Run() error {
	ui := o.UI
	sg := ui.StepGroup()
	defer sg.Wait()

	isUp, err := checkTunnel(o.Config.EnvDir)
	if err != nil {
		return fmt.Errorf("can't get tunnel status: %w", err)
	}

	if !isUp {
		return fmt.Errorf("can't get tunnel status: tunnel is down")
	}

	return nil
}
