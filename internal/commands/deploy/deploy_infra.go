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
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
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

var deployInfraLongDesc = templates.LongDesc(`
	Deploy only infraftructure.
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
		return fmt.Errorf("env must be specified")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified")
	}

	return nil
}

func (o *DeployInfraOptions) Run() error {
	logrus.Infof("infra: %s", o.Terraform)

	//terraform init
	opts := terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"init", "-input=true"},
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
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform init")
	}

	err := terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "init")
		return err
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform init completed")
	} else {
		pterm.Success.Println("terraform init completed")
	}

	//terraform plan
	opts = terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"plan"},
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Terraform.Profile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Terraform.Version,
	}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform plan")
	}

	err = terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "plan")
		return err
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform plan completed")
	} else {
		pterm.Success.Println("terraform plan completed")
	}

	//terraform apply
	opts = terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"apply", "-auto-approve"},
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Terraform.Profile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Terraform.Version,
	}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform apply")
	}

	err = terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "apply")
		return err
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform apply completed")
	} else {
		pterm.Success.Println("terraform apply completed")
	}

	// terraform output
	outputPath := fmt.Sprintf("%s/.terraform/output.json", viper.Get("ENV_DIR"))

	opts = terraform.Options{
		ContainerName: "terraform",
		Cmd:           []string{"output", "-json"},
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Terraform.Profile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Terraform.Version,
		OutputPath:       outputPath,
	}

	if logrus.GetLevel() < 4 {
		spinner, _ = pterm.DefaultSpinner.Start("execution terraform output")
	}

	err = terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", "output")
		return err
	}

	if logrus.GetLevel() < 4 {
		spinner.Success("terraform output completed")
	} else {
		pterm.Success.Println("terraform output completed")
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

	return nil
}

type terraformInfraConfig struct {
	Version string `mapstructure:"terraform_version,optional"`
	Region  string `mapstructure:"aws_region,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
