package destroy

import (
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DestroyInfraOptions struct {
	Config    *config.Config
	Terraform terraformInfraConfig
}

func NewDestroyInfraFlags() *DestroyInfraOptions {
	return &DestroyInfraOptions{}
}

func NewCmdDestroyInfra() *cobra.Command {
	o := NewDestroyInfraFlags()

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy anything",
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

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.Terraform.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.Terraform.Profile, "infra.terraform.aws-profile", "", "set aws profile")

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

func (o *DestroyInfraOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	BindFlags(cmd.Flags())

	if len(o.Terraform.Profile) == 0 {
		o.Terraform.Profile = viper.GetString("infra.terraform.aws_profile")
	}

	if len(o.Terraform.Profile) == 0 {
		o.Terraform.Profile = o.Config.AwsProfile
	}

	if len(o.Terraform.Version) == 0 {
		o.Terraform.Version = viper.GetString("infra.terraform.terraform_version")
	}

	if len(o.Terraform.Version) == 0 {
		o.Terraform.Version = viper.GetString("terraform_version")
	}

	fmt.Println(o.Terraform)

	return nil
}

func (o *DestroyInfraOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *DestroyInfraOptions) Run() error {
	logrus.Infof("infra: %s", o.Terraform)

	opts := terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"destroy", "-auto-approve"},
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Terraform.Profile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Terraform.Version,
	}

	spinner := &pterm.SpinnerPrinter{}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform destroy")
	}

	err := terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "destroy")
		return err
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform destroy completed")
	} else {
		pterm.Success.Println("terraform destroy completed")
	}

	return nil
}

type terraformInfraConfig struct {
	Version string `mapstructure:"terraform_version,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
