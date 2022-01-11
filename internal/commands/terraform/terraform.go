package terraform

import (
	"fmt"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TerraformOptions struct {
	Env     string
	Profile string
	Version string
	Command []string
}

func NewTerraformFlags() *TerraformOptions {
	return &TerraformOptions{}
}

func NewCmdTerraform() *cobra.Command {
	o := NewTerraformFlags()

	cmd := &cobra.Command{
		Use:   "terraform",
		Short: "generate terraform files",
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
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Command = args[argsLenAtDash:]
	o.Env = viper.GetString("env")
	o.Profile = viper.GetString("aws_profile")

	if o.Profile == "" {
		o.Profile = viper.GetString("aws-profile")
	}

	if o.Version == "" {
		o.Version = viper.GetString("terraform_version")
	}

	return nil
}

func (o *TerraformOptions) Validate() error {
	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Version) == 0 {
		return fmt.Errorf("terraform version must be specified")
	}

	if len(o.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}
	return nil
}

func (o *TerraformOptions) Run(args []string) error {
	opts := terraform.Options{
		ContainerName: "terraform",
		Cmd:           o.Command,
		Env: []string{
			fmt.Sprintf("ENV=%v", o.Env),
			fmt.Sprintf("AWS_PROFILE=%v", o.Profile),
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

	pterm.DefaultSection.Printfln("Terraform %s completed", args[0])

	return nil
}
