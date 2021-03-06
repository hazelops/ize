package terraform

import (
	"fmt"
	"os"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Options struct {
	Config  *config.Project
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

func NewTerraformFlags() *Options {
	return &Options{}
}

func NewCmdTerraform() *cobra.Command {
	o := NewTerraformFlags()

	cmd := &cobra.Command{
		Use:                   "terraform [flags] <terraform command> [terraform flags]",
		Short:                 "Run terraform",
		Long:                  terraformLongDesc,
		Example:               terraformExample,
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
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

			err = o.Run(args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *Options) Complete() error {
	var (
		cfg *config.Project
		err error
	)

	if err = config.CheckRequirements(config.WithIzeStructure(), config.WithConfigFile()); err != nil {
		return err
	}

	cfg, err = config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *Options) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}
	return nil
}

func (o *Options) Run(args []string) error {
	var tf terraform.Terraform

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	env := []string{
		fmt.Sprintf("ENV=%v", o.Config.Env),
		fmt.Sprintf("USER=%v", os.Getenv("USER")),
		fmt.Sprintf("AWS_PROFILE=%v", o.Config.AwsProfile),
		fmt.Sprintf("TF_LOG=%v", o.Config.TFLog),
		fmt.Sprintf("TF_LOG_PATH=%v", o.Config.TFLogPath),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%v", v.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%v", v.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%v", v.SessionToken),
	}

	if o.Config.Terraform != nil {
		o.Version = o.Config.Terraform["infra"].Version
	}

	if o.Version == "" {
		o.Version = o.Config.TerraformVersion
	}

	switch o.Config.PreferRuntime {
	case "docker":
		tf = terraform.NewDockerTerraform(o.Version, args, env, nil, o.Config.Home, o.Config.InfraDir, o.Config.EnvDir)
	case "native":
		tf = terraform.NewLocalTerraform(o.Version, args, env, nil, o.Config.EnvDir)
		err = tf.Prepare()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("can't supported %s runtime", o.Config.PreferRuntime)
	}

	logrus.Debug("starting terraform")

	err = tf.Run()
	if err != nil {
		return err
	}

	return nil
}
