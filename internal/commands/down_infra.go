package commands

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/requirements"
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
	OnlyInfra  bool
}

func NewDownInfraFlags(project *config.Project) *DownInfraOptions {
	return &DownInfraOptions{
		Config: project,
	}
}

func NewCmdDownInfra(project *config.Project) *cobra.Command {
	o := NewDownInfraFlags(project)

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
	cmd.Flags().BoolVar(&o.OnlyInfra, "only-infra", false, "down only infra state")

	return cmd
}

func (o *DownInfraOptions) Complete() error {
	if err := requirements.CheckRequirements(requirements.WithIzeStructure(), requirements.WithConfigFile()); err != nil {
		return err
	}

	if o.Config.Terraform == nil {
		return fmt.Errorf("you must specify at least one terraform stack in ize.toml")
	}

	if _, ok := o.Config.Terraform["infra"]; ok {
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

	if _, ok := o.Config.Terraform["infra"]; ok {
		err := destroyInfra("infra", o.Config, o.SkipGen, ui)
		if err != nil {
			return err
		}
	}

	err := manager.InReversDependencyOrder(aws.BackgroundContext(), o.Config.GetApps(), func(c context.Context, name string) error {
		return destroyInfra(name, o.Config, o.SkipGen, ui)
	})
	if err != nil {

	}

	return nil
}
