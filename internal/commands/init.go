package commands

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hazelops/ize/internal/template"
	"github.com/spf13/cobra"
)

type initCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newInitCmd() *mfaCmd {
	cc := &mfaCmd{}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := InitConfigFile()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func InitConfigFile() error {
	var qs = []*survey.Question{
		{
			Prompt: &survey.Input{
				Message: " Env:",
				Default: os.Getenv("ENV"),
			},
			Validate: survey.Required,
			Name:     "env",
		},
		{
			Prompt: &survey.Input{
				Message: " aws region:",
				Default: os.Getenv("AWS_REGION"),
			},
			Validate: survey.Required,
			Name:     "aws_region",
		},
		{
			Prompt: &survey.Input{
				Message: " aws profile:",
				Default: os.Getenv("AWS_PROFILE"),
			},
			Validate: survey.Required,
			Name:     "aws_profile",
		},
		{
			Prompt: &survey.Input{
				Message: " namespace:",
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

	err = template.GenerateConfigFile(opts)
	if err != nil {
		return err
	}

	return nil
}
