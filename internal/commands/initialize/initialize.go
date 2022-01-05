package initialize

import (
	"fmt"
	"os"
	"reflect"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InitOptions struct {
	Path string
}

func NewInitFlags() *InitOptions {
	return &InitOptions{}
}

func NewCmdInit() *cobra.Command {
	o := NewInitFlags()

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates an IZE configuration file.",
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

	cmd.Flags().StringVar(&o.Path, "path", os.Getenv("IZE_FILE"), "config file path")

	return cmd
}

func (o *InitOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.Path == "" {
		o.Path = "./ize.toml"
	}

	return nil
}

func (o *InitOptions) Validate() error {
	if len(o.Path) == 0 {
		return fmt.Errorf("path must be specified")
	}

	return nil
}

func (o *InitOptions) Run() error {
	env, exist := os.LookupEnv("IZE_ENV")
	if !exist {
		env = os.Getenv("ENV")
	}

	region, exist := os.LookupEnv("IZE_AWS_REGION")
	if !exist {
		region = os.Getenv("AWS_REGION")
	}

	profile, exist := os.LookupEnv("IZE_AWS_PROFILE")
	if !exist {
		profile = os.Getenv("AWS_PROFILE")
	}

	namespace, exist := os.LookupEnv("IZE_NAMESPACE")
	if !exist {
		namespace = os.Getenv("NAMESPACE")
	}

	var qs = []*survey.Question{
		{
			Prompt: &survey.Input{
				Message: " env:",
				Default: env,
			},
			Validate: survey.Required,
			Name:     "env",
		},
		{
			Prompt: &survey.Input{
				Message: " aws region:",
				Default: region,
			},
			Validate: survey.Required,
			Name:     "aws_region",
		},
		{
			Prompt: &survey.Input{
				Message: " aws profile:",
				Default: profile,
			},
			Validate: survey.Required,
			Name:     "aws_profile",
		},
		{
			Prompt: &survey.Input{
				Message: " namespace:",
				Default: namespace,
			},
			Validate: survey.Required,
			Name:     "namespace",
		},
		{
			Prompt: &survey.Input{
				Message: " terraform version:",
			},
			Validate: survey.Required,
			Name:     "terraform_version",
		},
	}

	opts := ConfigOpts{}

	err := survey.Ask(qs, &opts, survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = " ??"
		is.Question.Format = "black:green"
		is.Error.Format = "black:red"
	}))
	if err != nil {
		return err
	}

	v := reflect.ValueOf(opts)
	typeOfOpts := v.Type()

	viper.Reset()
	viper.SetConfigType("toml")

	for i := 0; i < v.NumField(); i++ {
		viper.Set(typeOfOpts.Field(i).Name, v.Field(i).Interface())
	}

	viper.WriteConfigAs(o.Path)

	return nil
}

type ConfigOpts struct {
	Env               string
	Aws_profile       string
	Aws_region        string
	Terraform_version string
	Namespace         string
}
