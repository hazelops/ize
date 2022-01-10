package env

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/template"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type EnvOptions struct {
	Env       string
	Namespace string
	Profile   string
	Region    string
}

func NewEnvFlags() *EnvOptions {
	return &EnvOptions{}
}

func NewCmdEnv() *cobra.Command {
	o := NewEnvFlags()

	cmd := &cobra.Command{
		Use:   "env",
		Short: "generate terraform files",
		RunE: func(cmd *cobra.Command, args []string) error {
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

	return cmd
}

func (o *EnvOptions) Complete(cmd *cobra.Command, args []string) error {
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Env = viper.GetString("env")
	o.Namespace = viper.GetString("namespace")

	o.Profile = viper.GetString("aws_profile")
	o.Region = viper.GetString("aws_region")

	if o.Region == "" {
		o.Region = viper.GetString("aws-region")
	}

	if o.Profile == "" {
		o.Profile = viper.GetString("aws-profile")
	}

	return nil
}

func (o *EnvOptions) Validate() error {
	if len(o.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Namespace) == 0 {
		return fmt.Errorf("namespace must be specified")
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("AWS region must be specified")
	}
	return nil
}

func (o *EnvOptions) Run() error {
	pterm.DefaultSection.Printfln("Starting generate terraform files")

	backendOpts := template.BackendOpts{
		ENV:                            o.Env,
		LOCALSTACK_ENDPOINT:            "",
		TERRAFORM_STATE_BUCKET_NAME:    fmt.Sprintf("%s-tf-state", o.Namespace),
		TERRAFORM_STATE_KEY:            fmt.Sprintf("%v/terraform.tfstate", o.Env),
		TERRAFORM_STATE_REGION:         o.Region,
		TERRAFORM_STATE_PROFILE:        o.Profile,
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
		ENV:               o.Env,
		AWS_PROFILE:       o.Profile,
		AWS_REGION:        o.Region,
		EC2_KEY_PAIR_NAME: fmt.Sprintf("%v-%v", o.Env, o.Namespace),
		TAG:               o.Env,
		SSH_PUBLIC_KEY:    string(key)[:len(string(key))-1],
		DOCKER_REGISTRY:   viper.GetString("DOCKER_REGISTRY"),
		NAMESPACE:         o.Namespace,
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
