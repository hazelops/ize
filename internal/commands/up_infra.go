package commands

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/manager"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type UpInfraOptions struct {
	Config     *config.Project
	SkipGen    bool
	AwsProfile string
	AwsRegion  string
	Version    string
	UI         terminal.UI
	Explain    bool
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

func NewUpInfraFlags(project *config.Project) *UpInfraOptions {
	return &UpInfraOptions{
		Config: project,
	}
}

func NewCmdUpInfra(project *config.Project) *cobra.Command {
	o := NewUpInfraFlags(project)

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
	cmd.Flags().BoolVar(&o.Explain, "explain", false, "bash alternative shown")
	cmd.Flags().StringVar(&o.Version, "infra.terraform.version", "", "set terraform version")
	cmd.Flags().StringVar(&o.AwsRegion, "infra.terraform.aws-region", "", "set aws region")
	cmd.Flags().StringVar(&o.AwsProfile, "infra.terraform.aws-profile", "", "set aws profile")

	return cmd
}

func (o *UpInfraOptions) Complete() error {
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

		if len(o.Config.Terraform["infra"].StateBucketRegion) == 0 {
			o.Config.Terraform["infra"].StateBucketRegion = o.Config.Terraform["infra"].AwsRegion
		}

		if len(o.Version) != 0 {
			o.Config.Terraform["infra"].Version = o.Version
		}

		if len(o.Config.Terraform["infra"].Version) == 0 {
			o.Config.Terraform["infra"].Version = o.Config.TerraformVersion
		}
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
	if o.Explain {
		tmpl := `# Change to the dir
cd {{.EnvDir}}

# Generate Backend.tf
cat << EOF > backend.tf
provider "aws" {
  profile = var.aws_profile
  region  = var.aws_region
  default_tags {
    tags = {
      env       = "{{.Env}}"
      namespace = "{{.Namespace}}"
      terraform = "true"
    }
  }
}

terraform {
  backend "s3" {
    bucket         = "{{.Namespace}}-tf-state"
    key            = "{{.Env}}/terraform.tfstate"
    region         = "{{.Terraform.infra.StateBucketRegion}}"
    profile        = "{{.AwsProfile}}"
    dynamodb_table = "tf-state-lock"
  }
}

EOF

# Generate variables.tfvars
cat << EOF > variables.tfvars
env               = "{{.Env}}"
aws_profile       = "{{.AwsProfile}}"
aws_region        = "{{.AwsRegion}}"
ec2_key_pair_name = "{{.Env}}-{{.Namespace}}"
docker_image_tag  = "{{.Tag}}"
ssh_public_key    = ""
docker_registry   = "{{.DockerRegistry}}"
namespace         = "{{.Namespace}}"
root_domain_name  = "{{.Terraform.infra.RootDomainName}}"

EOF

# Ensure Terraform is v {{.TerraformVersion}}
terraform --version

# Terraform Plan
terraform plan

# Terraform Apply
terraform apply
`
		err := o.Config.Generate(tmpl, nil)
		if err != nil {
			return err
		}

		return nil
	}

	ui := o.UI

	if _, ok := o.Config.Terraform["infra"]; ok {
		err := deployInfra("infra", ui, o.Config, o.SkipGen)
		if err != nil {
			return err
		}
	}

	err := manager.InDependencyOrder(aws.BackgroundContext(), o.Config.GetStates(), func(c context.Context, name string) error {
		return deployInfra(name, ui, o.Config, o.SkipGen)
	})
	if err != nil {
		return err
	}

	return nil
}

func deployInfra(name string, ui terminal.UI, config *config.Project, skipGen bool) error {
	if !skipGen {
		err := GenerateTerraformFiles(name, "", config)
		if err != nil {
			return err
		}
	}

	var tf terraform.Terraform

	logrus.Infof("infra: %s", config.Terraform[name])

	v, err := config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't get AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", config.Env),
		fmt.Sprintf("AWS_PROFILE=%v", config.Terraform[name].AwsProfile),
		fmt.Sprintf("TF_LOG=%v", config.TFLog),
		fmt.Sprintf("TF_LOG_PATH=%v", config.TFLogPath),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	switch config.PreferRuntime {
	case "docker":
		tf = terraform.NewDockerTerraform(name, []string{"init", "-input=true"}, env, nil, config)
	case "native":
		tf = terraform.NewLocalTerraform(name, []string{"init", "-input=true"}, env, nil, config)
		err = tf.Prepare()
		if err != nil {
			return fmt.Errorf("can't deploy infra: %w", err)
		}
	default:
		return fmt.Errorf("can't supported %s runtime", config.PreferRuntime)
	}

	ui.Output(fmt.Sprintf("[%s][%s] Running deploy infra...", config.Env, name), terminal.WithHeaderStyle())
	ui.Output("Execution terraform init...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy infra: %w", err)
	}

	ui.Output("Execution terraform plan...", terminal.WithHeaderStyle())

	outPath := filepath.Join(config.EnvDir, name, ".terraform", "tfplan")
	if name == "infra" {
		outPath = filepath.Join(config.EnvDir, ".terraform", "tfplan")
	}

	//terraform plan run options
	tf.NewCmd([]string{"plan", fmt.Sprintf("-out=%s", outPath)})

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy infra: %w", err)
	}

	//terraform apply run options
	tf.NewCmd([]string{"apply", "-auto-approve", outPath})

	ui.Output("Execution terraform apply...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy infra: %w", err)
	}

	//terraform output run options

	tf.NewCmd([]string{"output", "-json"})

	var output bytes.Buffer

	tf.SetOut(&output)

	ui.Output("Execution terraform output...", terminal.WithHeaderStyle())

	err = tf.RunUI(ui)
	if err != nil {
		return fmt.Errorf("can't deploy infra: %w", err)
	}

	parameterName := fmt.Sprintf("/%s/terraform-output", config.Env)

	byteValue, _ := ioutil.ReadAll(&output)
	sDec := base64.StdEncoding.EncodeToString(byteValue)
	if err != nil {
		return err
	}

	_, err = ssm.New(config.Session).PutParameter(&ssm.PutParameterInput{
		Name:      &parameterName,
		Value:     aws.String(sDec),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Overwrite: aws.Bool(true),
		Tier:      aws.String(ssm.ParameterTierIntelligentTiering),
		DataType:  aws.String("text"),
	})
	if err != nil {
		return err
	}

	ui.Output("Deploy infra completed!\n", terminal.WithSuccessStyle())

	return nil
}
