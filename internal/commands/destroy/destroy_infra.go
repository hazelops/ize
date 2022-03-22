package destroy

import (
	"fmt"
	"strings"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DestroyInfraOptions struct {
	Config    *config.Config
	Terraform terraformInfraConfig
	Local     bool
}

func NewDestroyInfraFlags() *DestroyInfraOptions {
	return &DestroyInfraOptions{}
}

func NewCmdDestroyInfra(ui terminal.UI) *cobra.Command {
	o := NewDestroyInfraFlags()

	cmd := &cobra.Command{
		Use:   "infra",
		Short: "destroy infra",
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

			err = o.Run(ui)
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
	var err error
	o.Local = viper.GetBool("local-terraform")

	if o.Local {
		o.Config, err = config.InitializeConfig()
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
		}
	} else {
		o.Config, err = config.InitializeConfig(config.WithDocker())
		if err != nil {
			return fmt.Errorf("can`t complete options: %w", err)
		}
	}

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

	return nil
}

func (o *DestroyInfraOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	return nil
}

func (o *DestroyInfraOptions) Run(ui terminal.UI) error {
	var tf terraform.Terraform

	logrus.Infof("infra: %s", o.Terraform)

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", o.Terraform.Profile),
		fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
		fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	if o.Local {
		tf = terraform.NewLocalTerraform(o.Terraform.Version, []string{"destroy", "-auto-approve"}, env, "")
		tf.Prepare()
	} else {
		tf = terraform.NewDockerTerraform(o.Terraform.Version, []string{"destroy", "-auto-approve"}, env, "")
	}

	ui.Output("Running terraform destoy...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	ui.Output("terraform destoy completed!\n", terminal.WithSuccessStyle())

	return nil
}

type terraformInfraConfig struct {
	Version string `mapstructure:"terraform_version,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
