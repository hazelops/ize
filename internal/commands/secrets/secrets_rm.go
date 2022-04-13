package secrets

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SecretsRemoveOptions struct {
	Config      *config.Config
	AppName     string
	Backend     string
	SecretsPath string
	ui          terminal.UI
}

func NewSecretsRemoveFlags() *SecretsRemoveOptions {
	return &SecretsRemoveOptions{}
}

func NewCmdSecretsRemove() *cobra.Command {
	o := NewSecretsRemoveFlags()

	cmd := &cobra.Command{
		Use:              "rm",
		Short:            "remove secrets from storage",
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
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg
	o.AppName = cmd.Flags().Args()[0]

	if o.SecretsPath == "" {
		o.SecretsPath = fmt.Sprintf("/%s/%s", o.Config.Env, o.AppName)
	}

	o.ui = terminal.ConsoleUI(context.Background(), viper.GetBool("plain-text"))

	return nil
}

func (o *SecretsRemoveOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	return nil
}

func (o *SecretsRemoveOptions) Run() error {
	ui := o.ui
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("removing secrets for %s...", o.AppName)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 50) }()

	if o.Backend == "ssm" {
		err := o.rm(s)
		if err != nil {
			pterm.DefaultSection.Sprintfln("secrets have been removed from %s", o.SecretsPath)
			return err
		}
	} else {
		return fmt.Errorf("backend %s is not found or not supported", o.Backend)
	}

	s.Done()
	ui.Output("Removing secrets complete!\n", terminal.WithSuccessStyle())

	return nil
}

func (o *SecretsRemoveOptions) rm(s terminal.Step) error {
	if o.SecretsPath == "" {
		fmt.Fprintf(s.TermOutput(), "path was not set...")
		return nil
	}

	fmt.Fprintf(s.TermOutput(), "removing secrets from %s://%s...\n", o.Backend, o.SecretsPath)

	ssmSvc := ssm.New(o.Config.Session)

	out, err := ssmSvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: &o.SecretsPath,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(s.TermOutput(), "getting secrets...\n")

	if len(out.Parameters) == 0 {
		fmt.Fprintf(s.TermOutput(), "no values found...\n")
		fmt.Fprintf(s.TermOutput(), "removing secrets...\n")
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

	fmt.Fprintf(s.TermOutput(), "removing secrets...\n")

	return nil
}
