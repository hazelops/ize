package configure

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigureOptions struct {
}

func NewConfigFlags() *ConfigureOptions {
	return &ConfigureOptions{}
}

func NewCmdConfig() *cobra.Command {
	o := NewConfigFlags()

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Generate global configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *ConfigureOptions) Run() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	gc := map[string]map[string]map[string]string{}

	viper.SetConfigFile(fmt.Sprintf("%s/.ize/config.toml", home))
	viper.ReadInConfig()
	viper.Unmarshal(&gc)

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

	err = survey.AskOne(&survey.Input{
		Message: " namespace:",
		Default: namespace,
	}, &namespace, survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = " ??"
		is.Question.Format = "black:green"
		is.Error.Format = "black:red"
	}))
	if err != nil {
		return err
	}

	err = survey.AskOne(&survey.Input{
		Message: " env:",
		Default: env,
	}, &env, survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = " ??"
		is.Question.Format = "black:green"
		is.Error.Format = "black:red"
	}))
	if err != nil {
		return err
	}

	var qs = []*survey.Question{
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
				Message: " terraform version:",
			},
			Validate: survey.Required,
			Name:     "terraform_version",
		},
	}

	opts := Config{}

	err = survey.Ask(qs, &opts, survey.WithIcons(func(is *survey.IconSet) {
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

	if gc[namespace] == nil {
		gc[namespace] = make(map[string]map[string]string)
	}

	if gc[namespace][env] == nil {
		gc[namespace][env] = make(map[string]string)
	}

	for i := 0; i < v.NumField(); i++ {
		gc[namespace][env][strings.ToLower(typeOfOpts.Field(i).Name)] = v.Field(i).String()
	}

	raw := make(map[string]interface{}, len(gc))
	for k, v := range gc {
		raw[k] = v
	}

	viper.MergeConfigMap(raw)

	os.MkdirAll(fmt.Sprintf("%s/.ize", home), 0755)

	err = viper.WriteConfigAs(fmt.Sprintf("%s/.ize/config.toml", home))
	if err != nil {
		return fmt.Errorf("can't write config: %w", err)
	}

	return nil
}

type Namespace struct {
	Env map[string]Config
}

type Config struct {
	Aws_profile       string
	Aws_region        string
	Terraform_version string
}
