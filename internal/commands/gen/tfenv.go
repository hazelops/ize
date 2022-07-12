package gen

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/template"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TfenvOptions struct {
	Config                   *config.Project
	TerraformStateBucketName string
}

var tfenvLongDesc = templates.LongDesc(`
	tfenv generates backend.tf and variable.tfvars files.
`)

var tfenvExample = templates.Examples(`
	# Generate files
	ize tfenv

	# Generate files via config file
	ize --config-file /path/to/config tfenv

	# Generate files via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize tfenv
`)

func NewTfenvFlags() *TfenvOptions {
	return &TfenvOptions{}
}

func NewCmdTfenv() *cobra.Command {
	o := NewTfenvFlags()

	cmd := &cobra.Command{
		Use:     "tfenv",
		Short:   "Generate terraform files",
		Long:    tfenvLongDesc,
		Example: tfenvExample,
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

	cmd.Flags().StringVar(&o.TerraformStateBucketName, "terraform-state-bucket-name", "", "set terraform state bucket name (default <NAMESPACE>-tf-state)")

	return cmd
}

func (o *TfenvOptions) Complete() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *TfenvOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	if len(o.Config.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified\n")
	}

	if len(o.TerraformStateBucketName) == 0 {
		o.TerraformStateBucketName = fmt.Sprintf("%s-tf-state", o.Config.Namespace)
	}

	return nil
}

func (o *TfenvOptions) Run() error {
	pterm.DefaultSection.Printfln("Starting generate terraform files")

	var tf config.Terraform
	if o.Config.Terraform != nil {
		tf = *o.Config.Terraform["infra"]
	}

	if len(o.TerraformStateBucketName) != 0 {
		tf.StateBucketName = o.TerraformStateBucketName
	}

	if len(tf.StateBucketRegion) == 0 {
		tf.StateBucketRegion = o.Config.AwsRegion

	}

	backendOpts := template.BackendOpts{
		ENV:                            o.Config.Env,
		LOCALSTACK_ENDPOINT:            "",
		TERRAFORM_STATE_BUCKET_NAME:    tf.StateBucketName,
		TERRAFORM_STATE_KEY:            fmt.Sprintf("%v/terraform.tfstate", o.Config.Env),
		TERRAFORM_STATE_REGION:         tf.StateBucketRegion,
		TERRAFORM_STATE_PROFILE:        o.Config.AwsProfile,
		TERRAFORM_STATE_DYNAMODB_TABLE: "tf-state-lock",
		TERRAFORM_AWS_PROVIDER_VERSION: "",
	}
	envDir := o.Config.EnvDir

	logrus.Debugf("backend opts: %s", backendOpts)
	logrus.Debugf("ENV dir path: %s", envDir)

	err := template.GenerateBackendTf(
		backendOpts,
		envDir,
	)

	if err != nil {
		pterm.DefaultSection.Println("Generate terraform file not completed")
		return err
	}

	pterm.Success.Println("backend.tf generated")

	pterm.Success.Printfln("Read SSH public key")

	home, _ := os.UserHomeDir()
	key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", home))
	if err != nil {
		pterm.DefaultSection.Println("Generate terraform file not completed")
		return err
	}

	if err != nil {
		pterm.DefaultSection.Println("Generate terraform file not completed")
		return err
	}

	varsOpts := template.VarsOpts{
		ENV:               o.Config.Env,
		AWS_PROFILE:       o.Config.AwsProfile,
		AWS_REGION:        o.Config.AwsRegion,
		EC2_KEY_PAIR_NAME: fmt.Sprintf("%v-%v", o.Config.Env, o.Config.Namespace),
		ROOT_DOMAIN_NAME:  tf.RootDomainName,
		TAG:               o.Config.Env,
		SSH_PUBLIC_KEY:    string(key)[:len(string(key))-1],
		DOCKER_REGISTRY:   o.Config.DockerRegistry,
		NAMESPACE:         o.Config.Namespace,
	}

	logrus.Debugf("backend opts: %s", varsOpts)
	logrus.Debugf("ENV dir path: %s", envDir)

	err = template.GenerateVarsTf(
		varsOpts,
		envDir,
	)

	if err != nil {
		pterm.DefaultSection.Println("Generate terraform file not completed")
		return err
	}

	pterm.Success.Println("terraform.tfvars generated")

	if err != nil {
		pterm.DefaultSection.Println("Generate terraform file not completed")
		return err
	}

	pterm.DefaultSection.Printfln("Generate terraform files completed")

	return nil
}
