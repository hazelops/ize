package down

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

type DownInfraOptions struct {
	Config     *config.Project
	ui         terminal.UI
	Version    string
	AwsProfile string
	AwsRegion  string
	SkipGen    bool
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

	cmd.Flags().StringVar(&o.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.AwsProfile, "infra.terraform.aws-profile", "", "set aws profile")
	cmd.Flags().StringVar(&o.AwsRegion, "infra.terraform.aws-region", "", "set aws region")
	cmd.Flags().BoolVar(&o.SkipGen, "skip-gen", false, "skip generating terraform files")

	return cmd
}

func (o *DownInfraOptions) Complete() error {
	var err error

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't load options for a command: %w", err)
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

	o.ui = terminal.ConsoleUI(context.Background(), o.Config.PlainText)

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
	return destroyInfra(ui, o.Config, o.SkipGen)
}
