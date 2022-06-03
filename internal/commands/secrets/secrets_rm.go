package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type SecretsRemoveOptions struct {
	Config      *config.Config
	AppName     string
	Backend     string
	SecretsPath string
	ui          terminal.UI
}

var secretsRemoveExample = templates.Examples(`
	# Remove secrets:

    # This will remove your secrets for "squibby" service
 	ize secrets rm squibby
`)

func NewSecretsRemoveFlags() *SecretsRemoveOptions {
	return &SecretsRemoveOptions{}
}

func NewCmdSecretsRemove() *cobra.Command {
	o := NewSecretsRemoveFlags()

	cmd := &cobra.Command{
		Use:              "rm <app>",
		Example:          secretsRemoveExample,
		Short:            "Remove secrets from storage",
		Long:             "This command removes secrets from storage",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete(cmd, args)
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

	cmd.Flags().StringVar(&o.Backend, "backend", "ssm", "backend type")
	cmd.Flags().StringVar(&o.SecretsPath, "path", "", "path to secrets")

	return cmd
}

func (o *SecretsRemoveOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg
	o.AppName = cmd.Flags().Args()[0]

	if o.SecretsPath == "" {
		o.SecretsPath = fmt.Sprintf("/%s/%s", o.Config.Env, o.AppName)
	}

	o.ui = terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)

	return nil
}

func (o *SecretsRemoveOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *SecretsRemoveOptions) Run() error {
	s, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Removing secrets for %s...", o.AppName))
	if o.Backend == "ssm" {
		err := o.rm(s)
		if err != nil {
			pterm.DefaultSection.Sprintfln("Secrets have been removed from %s", o.SecretsPath)
			return err
		}
	} else {
		return fmt.Errorf("backend %s is not found or not supported", o.Backend)
	}

	s.Success("Removing secrets complete!")

	return nil
}

func (o *SecretsRemoveOptions) rm(s *pterm.SpinnerPrinter) error {
	if o.SecretsPath == "" {
		s.UpdateText("Path was not set...")
		time.Sleep(2 * time.Second)
		return nil
	}

	s.UpdateText(fmt.Sprintf("Removing secrets from %s://%s...", o.Backend, o.SecretsPath))

	ssmSvc := ssm.New(o.Config.Session)

	out, err := ssmSvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: &o.SecretsPath,
	})
	if err != nil {
		return err
	}

	s.UpdateText("Getting secrets...")

	if len(out.Parameters) == 0 {
		s.UpdateText("No values found...")
		time.Sleep(2 * time.Second)
		s.UpdateText("Removing secrets...")
		time.Sleep(1 * time.Second)
		return nil
	}

	var names []*string

	for _, p := range out.Parameters {
		names = append(names, p.Name)
	}

	_, err = ssmSvc.DeleteParameters(&ssm.DeleteParametersInput{
		Names: names,
	})
	if err != nil {
		return err
	}

	s.UpdateText("Removing secrets...")

	return nil
}
