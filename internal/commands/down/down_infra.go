package down

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type DownInfraOptions struct {
	Config *config.Project
	ui     terminal.UI

	Version    string
	AwsProfile string
	AwsRegion  string
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
	var tf terraform.Terraform

	logrus.Infof("infra: %s", o.Config.Terraform["infra"])

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", o.Config.Terraform["infra"].AwsProfile),
		fmt.Sprintf("TF_LOG=%v", o.Config.TFLog),
		fmt.Sprintf("TF_LOG_PATH=%v", o.Config.TFLogPath),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	switch o.Config.PreferRuntime {
	case "docker":
		tf = terraform.NewDockerTerraform(o.Config.Terraform["infra"].Version, []string{"destroy", "-auto-approve"}, env, nil, o.Config.Home, o.Config.InfraDir, o.Config.EnvDir)
	case "native":
		tf = terraform.NewLocalTerraform(o.Config.Terraform["infra"].Version, []string{"destroy", "-auto-approve"}, env, nil, o.Config.EnvDir)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't destroy infra: %w", err)
		}
	default:
		return fmt.Errorf("can't supported %s runtime", o.Config.PreferRuntime)
	}

	ui.Output("Running terraform destroy...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	ui.Output("Terraform destroy completed!\n", terminal.WithSuccessStyle())

	return nil
}
