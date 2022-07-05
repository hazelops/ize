package down

import (
	"context"
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DownInfraOptions struct {
	Config  *config.Config
	Infra   Infra
	SkipGen bool
	ui      terminal.UI
}

func NewDownInfraFlags() *DownInfraOptions {
	return &DownInfraOptions{}
}

func NewCmdDownInfra() *cobra.Command {
	o := NewDownInfraFlags()

	cmd := &cobra.Command{
		Use:   "infra",
		Short: "Destroy infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete(cmd)
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

	cmd.Flags().StringVar(&o.Infra.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.Infra.Profile, "infra.terraform.aws-profile", "", "set aws profile")

	return cmd
}

func BindFlags(flags *pflag.FlagSet) {
	replacer := strings.NewReplacer("-", "_")

	flags.VisitAll(func(flag *pflag.Flag) {
		if err := viper.BindPFlag(replacer.Replace(flag.Name), flag); err != nil {
			panic("unable to bind flag " + flag.Name + ": " + err.Error())
		}
	})
}

func (o *DownInfraOptions) Complete(cmd *cobra.Command) error {
	var err error

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't load options for a command: %w", err)
	}

	BindFlags(cmd.Flags())

	if len(o.Infra.Profile) == 0 {
		o.Infra.Profile = viper.GetString("infra.terraform.aws_profile")
	}

	if len(o.Infra.Profile) == 0 {
		o.Infra.Profile = o.Config.AwsProfile
	}

	if len(o.Infra.Version) == 0 {
		o.Infra.Version = viper.GetString("infra.terraform.terraform_version")
	}

	if len(o.Infra.Version) == 0 {
		o.Infra.Version = viper.GetString("terraform_version")
	}

	o.ui = terminal.ConsoleUI(context.Background(), viper.GetBool("plain_text"))

	return nil
}

func (o *DownInfraOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *DownInfraOptions) Run() error {
	ui := o.ui
	return destroyInfra(ui, o.Infra, *o.Config, o.SkipGen)
}
