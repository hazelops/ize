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
	"github.com/spf13/viper"
)

type TfenvOptions struct {
	Config                   *config.Config
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

func (o *TfenvOptions) Run() error {
	return GenerateTerraformFiles(
		o.Config.AwsRegion,
		o.Config.AwsProfile,
		o.Config.Env,
		o.Config.Namespace,
		o.TerraformStateBucketName,
	)

}

func GenerateTerraformFiles(region, profile, env, namespace, stateBucketName string) error {
	pterm.DefaultSection.Printfln("Starting generate terraform files")

	if len(stateBucketName) == 0 {
		stateBucketName = viper.GetString("infra.terraform.state_bucket_name")
		if len(stateBucketName) == 0 {
			stateBucketName = fmt.Sprintf("%s-tf-state", namespace)
		}
	}

	awsStateRegion := region
	if len(viper.GetString("infra.terraform.state_bucket_region")) > 0 {
		awsStateRegion = viper.GetString("infra.terraform.state_bucket_region")
	}

	backendOpts := template.BackendOpts{
		ENV:                            env,
		LOCALSTACK_ENDPOINT:            "",
		TERRAFORM_STATE_BUCKET_NAME:    stateBucketName,
		TERRAFORM_STATE_KEY:            fmt.Sprintf("%v/terraform.tfstate", env),
		TERRAFORM_STATE_REGION:         awsStateRegion,
		TERRAFORM_STATE_PROFILE:        profile,
		TERRAFORM_STATE_DYNAMODB_TABLE: "tf-state-lock",
		TERRAFORM_AWS_PROVIDER_VERSION: "",
	}
	envDir := viper.GetString("ENV_DIR")

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
		ENV:               env,
		AWS_PROFILE:       profile,
		AWS_REGION:        region,
		EC2_KEY_PAIR_NAME: fmt.Sprintf("%v-%v", env, namespace),
		ROOT_DOMAIN_NAME:  viper.GetString("infra.terraform.root_domain_name"),
		TAG:               fmt.Sprintf("%s-latest", env),
		SSH_PUBLIC_KEY:    string(key)[:len(string(key))-1],
		DOCKER_REGISTRY:   viper.GetString("DOCKER_REGISTRY"),
		NAMESPACE:         namespace,
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
