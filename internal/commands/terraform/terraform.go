package terraform

import (
	"fmt"
	"os"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TerraformOptions struct {
	Config  *config.Config
	Version string
	Command []string
	Local   bool
}

var terraformLongDesc = templates.LongDesc(`
	Run terraform command via terraform docker container.
	By default, terraform runs locally.
	At the same time, terraform will be downloaded and launched from ~/.ize/versions/terraform/

	To use a docker terraform, set value of "docker" to the --prefer-runtime global flag.
`)

var terraformExample = templates.Examples(`
	# Run terraform init
	ize -e dev -p default -r us-east-1 -n hazelops terraform --version 1.0.10 init -input=true

	# Run terraform plan via config file
	ize --config-file (or -c) /path/to/config terraform plan -out=$(ENV_DIR)/.terraform/tfplan -input=false

	# Run terraform init via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize terraform init -input=true

	# Run terraform in docker
	ize -e dev -p default -r us-east-1 -n hazelops --prefer-runtime=docker terraform --version 1.0.10 init -input=true
`)

func NewTerraformFlags() *TerraformOptions {
	return &TerraformOptions{}
}

func NewCmdTerraform() *cobra.Command {
	o := NewTerraformFlags()

	cmd := &cobra.Command{
		Use:                   "terraform [flags] <terraform command> [terraform flags]",
		Short:                 "run terraform",
		Long:                  terraformLongDesc,
		Example:               terraformExample,
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
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

			err = o.Run(args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *TerraformOptions) Complete(cmd *cobra.Command, args []string) error {
	var (
		cfg *config.Config
		err error
	)

	cfg, err = config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *TerraformOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}
	return nil
}

func (o *TerraformOptions) Run(args []string) error {
	var tf terraform.Terraform

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("USER=%v", os.Getenv("USER")),
		fmt.Sprintf("AWS_PROFILE=%v", o.Config.AwsProfile),
		fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
		fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	if o.Config.IsDockerRuntime {
		tf = terraform.NewDockerTerraform(viper.GetString("terraform_version"), args, env, nil)
	} else {
		tf = terraform.NewLocalTerraform(viper.GetString("terraform_version"), args, env, nil)
		err = tf.Prepare()
		if err != nil {
			return err
		}
	}

	logrus.Debug("starting terraform")

	err = tf.Run()
	if err != nil {
		return err
	}

	return nil
}
