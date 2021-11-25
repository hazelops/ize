package commands

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hazelops/ize/internal/template"
	"github.com/spf13/cobra"
)

type initCmd struct {
	*baseBuilderCmd

	filePath string
}

func (b *commandsBuilder) newInitCmd() *initCmd {
	cc := &initCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Сreates an IZE configuration file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := InitConfigFile(cc.filePath)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&cc.filePath, "path", os.Getenv("IZE_FILE"), "config file path")

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func InitConfigFile(path string) error {
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

	opts := template.ConfigOpts{}

	err := survey.Ask(qs, &opts, survey.WithIcons(func(is *survey.IconSet) {
		is.Question.Text = " ??"
		is.Question.Format = "black:green"
		is.Error.Format = "black:red"
	}))
	if err != nil {
		return err
	}

	err = template.GenerateConfigFile(opts, path)
	if err != nil {
		return err
	}

	return nil
}
