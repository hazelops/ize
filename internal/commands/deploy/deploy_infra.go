package deploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DeployInfraOptions struct {
	Config *config.Config
	Type   string
	Infra  Infra
	UI     terminal.UI
}

var deployInfraLongDesc = templates.LongDesc(`
	Only deploy infrastructure.
`)

var deployInfraExample = templates.Examples(`
	# Deploy infra via flags
	ize deploy infra --infra.terraform.version <version> --infra.terraform.aws-region <region> --infra.terraform.aws-profile <profile>

	# Deploy infra via config file
	ize --config-file /path/to/config deploy infra

	# Deploy infra via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize deploy infra
`)

func NewDeployInfraFlags() *DeployInfraOptions {
	return &DeployInfraOptions{}
}

func NewCmdDeployInfra() *cobra.Command {
	o := NewDeployInfraFlags()

	cmd := &cobra.Command{
		Use:     "infra",
		Short:   "Manage infra deployments",
		Long:    deployInfraLongDesc,
		Example: deployInfraExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
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

	cmd.Flags().StringVar(&o.Infra.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.Infra.Region, "infra.terraform.aws-region", "", "set aws region")
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

func (o *DeployInfraOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't deploy your stack: %w", err)
	}

	BindFlags(cmd.Flags())

	if len(o.Infra.Profile) == 0 {
		o.Infra.Profile = viper.GetString("infra.terraform.aws_profile")
	}

	if len(o.Infra.Profile) == 0 {
		o.Infra.Profile = o.Config.AwsProfile
	}

	if len(o.Infra.Region) == 0 {
		o.Infra.Region = viper.GetString("infra.terraform.aws_region")
	}

	if len(o.Infra.Region) == 0 {
		o.Infra.Region = o.Config.AwsRegion
	}

	if len(o.Infra.Version) == 0 {
		o.Infra.Version = viper.GetString("infra.terraform.terraform_version")
	}

	if len(o.Infra.Version) == 0 {
		o.Infra.Version = viper.GetString("terraform_version")
	}

	o.UI = terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)

	return nil
}

func (o *DeployInfraOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified")
	}

	return nil
}

func (o *DeployInfraOptions) Run(cmd *cobra.Command) error {
	ui := o.UI

	return deployInfra(ui, o.Infra, *o.Config)
}
