package deploy

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DeployInfraOptions struct {
	Config    *config.Config
	Type      string
	Terraform terraformInfraConfig
	Local     bool
	UI        terminal.UI
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
		Short:   "manage infra deployments",
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

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.Terraform.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.Terraform.Region, "infra.terraform.aws-region", "", "set aws region")
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

func (o *DeployInfraOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error

	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can`t complete options: %w", err)
	}

	BindFlags(cmd.Flags())

	if len(o.Terraform.Profile) == 0 {
		o.Terraform.Profile = viper.GetString("infra.terraform.aws_profile")
	}

	if len(o.Terraform.Profile) == 0 {
		o.Terraform.Profile = o.Config.AwsProfile
	}

	if len(o.Terraform.Region) == 0 {
		o.Terraform.Region = viper.GetString("infra.terraform.aws_region")
	}

	if len(o.Terraform.Region) == 0 {
		o.Terraform.Region = o.Config.AwsRegion
	}

	if len(o.Terraform.Version) == 0 {
		o.Terraform.Version = viper.GetString("infra.terraform.terraform_version")
	}

	if len(o.Terraform.Version) == 0 {
		o.Terraform.Version = viper.GetString("terraform_version")
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

func (o *DeployInfraOptions) Run() error {
	ui := o.UI
	var tf terraform.Terraform

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	logrus.Infof("infra: %s", o.Terraform)

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", o.Terraform.Profile),
		fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
		fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	if o.Config.IsDockerRuntime {
		tf = terraform.NewDockerTerraform(o.Terraform.Version, []string{"init", "-input=true"}, env, nil)
	} else {
		tf = terraform.NewLocalTerraform(o.Terraform.Version, []string{"init", "-input=true"}, env, nil)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't deploy infra: %w", err)
		}
	}

	ui.Output(fmt.Sprintf("[%s] Running deploy infra...", viper.Get("ENV")), terminal.WithHeaderStyle())
	ui.Output("Execution terraform init...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy infra: %w", err)
	}

	ui.Output("Execution terraform plan...", terminal.WithHeaderStyle())

	//terraform plan run options
	tf.NewCmd([]string{"plan"})

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	//terraform apply run options
	tf.NewCmd([]string{"apply", "-auto-approve"})

	ui.Output("Execution terraform apply...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	tf.NewCmd([]string{"output", "-json"})

	var output bytes.Buffer

	tf.SetOut(&output)

	ui.Output("Execution terraform output...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("/%s/terraform-output", o.Config.Env)

	byteValue, _ := ioutil.ReadAll(&output)
	sDec := base64.StdEncoding.EncodeToString(byteValue)
	if err != nil {
		return err
	}

	ssm.New(o.Config.Session).PutParameter(&ssm.PutParameterInput{
		Name:      &name,
		Value:     aws.String(string(sDec)),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String("Intelligent-Tiering"),
		DataType:  aws.String("text"),
	})

	ui.Output("Deploy infra completed!\n", terminal.WithSuccessStyle())

	return nil
}

type terraformInfraConfig struct {
	Version string `mapstructure:"terraform_version,optional"`
	Region  string `mapstructure:"aws_region,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
