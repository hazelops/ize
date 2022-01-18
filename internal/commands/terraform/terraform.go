package terraform

import (
	"fmt"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TerraformOptions struct {
	Config  *config.Config
	Version string
	Command []string
}

var terrafromLongDesc = templates.LongDesc(`
	Run terraform command via terraform docker container.
`)

var terraformExample = templates.Examples(`
	# Run terraform init
	ize -e dev -p default -r us-east-1 -n hazelops terraform --version 1.0.10 -- init -input=true

	# Run terraform plan via config file
	ize --config-file (or -c) /path/to/config terraform -- plan -out=$(ENV_DIR)/.terraform/tfplan -input=false

	# Run terraform init via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize terraform -- init -input=true
`)

func NewTerraformFlags() *TerraformOptions {
	return &TerraformOptions{}
}

func NewCmdTerraform() *cobra.Command {
	o := NewTerraformFlags()

	cmd := &cobra.Command{
		Use:                   "terraform [flags] -- <terraform command> [terraform flags]",
		Short:                 "run terraform",
		Long:                  terrafromLongDesc,
		Example:               terraformExample,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			argsLenAtDash := cmd.ArgsLenAtDash()

			err := o.Complete(cmd, args, argsLenAtDash)
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

	cmd.Flags().StringVar(&o.Version, "version", "", "set terraform version")

	return cmd
}

func (o *TerraformOptions) Complete(cmd *cobra.Command, args []string, argsLenAtDash int) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	o.Command = args[argsLenAtDash:]

	if o.Version == "" {
		o.Version = viper.GetString("terraform_version")
	}

	return nil
}

func (o *TerraformOptions) Validate() error {
	if len(o.Version) == 0 {
		return fmt.Errorf("terraform version must be specified")
	}

	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}
	return nil
}

func (o *TerraformOptions) Run(args []string) error {
	opts := terraform.Options{
		ContainerName: "terraform",
		Cmd:           o.Command,
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Config.AwsProfile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
		TerraformVersion: o.Version,
	}

	logrus.Debug("starting terraform")

	err := terraform.Run(opts)
	if err != nil {
		logrus.Errorf("terraform %s not completed", args[0])
		return err
	}

	logrus.Infof("terraform %s completed", args[0])

	return nil
}
