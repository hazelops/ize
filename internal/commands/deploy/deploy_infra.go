package deploy

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
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
}

func NewDeployInfraFlags() *DeployInfraOptions {
	return &DeployInfraOptions{}
}

func NewCmdDeployInfra(ui terminal.UI) *cobra.Command {
	o := NewDeployInfraFlags()

	cmd := &cobra.Command{
		Use:   "infra",
		Short: "manage infra deployments",
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

	return nil
}

func (o *DeployInfraOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified\n")
	}

	return nil
}

func (o *DeployInfraOptions) Run(ui terminal.UI) error {
	ui.Output("Running deploy infra...", terminal.WithHeaderStyle())

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

	//terraform init run options
	opts := terraform.Options{
		ContainerName:    "terraform",
		Cmd:              []string{"init", "-input=true"},
		Env:              env,
		TerraformVersion: o.Terraform.Version,
	}

	ui.Output("Execution terraform init...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return fmt.Errorf("can't deploy all: %w", err)
	}

	ui.Output("Execution terraform plan...", terminal.WithHeaderStyle())

	//terraform plan run options
	opts.Cmd = []string{"plan"}

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return err
	}

	//terraform apply run options
	opts.Cmd = []string{"apply", "-auto-approve"}

	ui.Output("Execution terraform apply...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return err
	}

	//terraform output run options
	outputPath := fmt.Sprintf("%s/.terraform/output.json", viper.Get("ENV_DIR"))

	opts.Cmd = []string{"output", "-json"}
	opts.OutputPath = outputPath

	ui.Output("Execution terraform output...", terminal.WithHeaderStyle())

	err = terraform.RunUI(ui, opts)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("/%s/terraform-output", o.Config.Env)

	outputFile, err := os.Open(outputPath)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	byteValue, _ := ioutil.ReadAll(outputFile)
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

	ui.Output("deploy infra completed!\n", terminal.WithSuccessStyle())

	return nil
}

type terraformInfraConfig struct {
	Version string `mapstructure:"terraform_version,optional"`
	Region  string `mapstructure:"aws_region,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
