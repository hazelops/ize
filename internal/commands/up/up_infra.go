package up

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type UpInfraOptions struct {
	Config     *config.Project
	SkipGen    bool
	AwsProfile string
	AwsRegion  string
	Version    string
	UI         terminal.UI
}

var upInfraLongDesc = templates.LongDesc(`
	Only deploy infrastructure.
`)

var upInfraExample = templates.Examples(`
	# Deploy infra with flags
	ize up infra --infra.terraform.version <version> --infra.terraform.aws-region <region> --infra.terraform.aws-profile <profile>

	# Deploy infra with explicitly specified config file
	ize --config-file /path/to/config up infra

	# Deploy infra with explicitly specified config file passed via environment variable
	export IZE_CONFIG_FILE=/path/to/config
	ize up infra
`)

func NewUpInfraFlags() *UpInfraOptions {
	return &UpInfraOptions{}
}

func NewCmdUpInfra() *cobra.Command {
	o := NewUpInfraFlags()

	cmd := &cobra.Command{
		Use:     "infra",
		Short:   "Manage infra deployments",
		Long:    upInfraLongDesc,
		Example: upInfraExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
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

	cmd.Flags().BoolVar(&o.SkipGen, "skip-gen", false, "skip generating terraform files")
	cmd.Flags().StringVar(&o.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.AwsRegion, "infra.terraform.aws-region", "", "set aws region")
	cmd.Flags().StringVar(&o.AwsProfile, "infra.terraform.aws-profile", "", "set aws profile")

	return cmd
}

func (o *UpInfraOptions) Complete() error {
	var err error

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't deploy your stack: %w", err)
	}

	if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile()); err != nil {
		return err
	}

	if o.Config.Terraform == nil {
		o.Config.Terraform = map[string]*config.Terraform{}
		o.Config.Terraform["infra"] = &config.Terraform{}
	}

	if len(o.AwsProfile) != 0 {
		o.Config.Terraform["infra"].AwsProfile = o.AwsProfile
	}

	if len(o.Config.Terraform["infra"].AwsProfile) == 0 {
		o.Config.Terraform["infra"].AwsProfile = o.Config.AwsProfile
	}

	if len(o.AwsProfile) != 0 {
		o.Config.Terraform["infra"].AwsRegion = o.AwsRegion
	}

	if len(o.Config.Terraform["infra"].AwsRegion) == 0 {
		o.Config.Terraform["infra"].AwsRegion = o.Config.AwsRegion
	}

	if len(o.Version) != 0 {
		o.Config.Terraform["infra"].Version = o.Version
	}

	if len(o.Config.Terraform["infra"].Version) == 0 {
		o.Config.Terraform["infra"].Version = o.Config.TerraformVersion
	}

	o.UI = terminal.ConsoleUI(context.Background(), o.Config.PlainText)

	return nil
}

func (o *UpInfraOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified")
	}

	return nil
}

func (o *UpInfraOptions) Run() error {
	ui := o.UI

	return deployInfra(ui, o.Config, o.SkipGen)
}
