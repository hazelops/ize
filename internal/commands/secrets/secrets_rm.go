package secrets

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SecretsRemoveOptions struct {
	Env string
	Region  string
	Profile string
	AppName string
	Backend        string
	SecretsPath string
}

func NewSecretsRemoveFlags() *SecretsRemoveOptions {
	return &SecretsRemoveOptions{}
}

func NewCmdSecretsRemove() *cobra.Command {
	o := NewSecretsRemoveFlags()

	cmd := &cobra.Command{
		Use:              "rm",
		Short:            "remove secrets from storage",
		Long:             "This command removes secrets from storage",
		TraverseChildren: true,
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

	cmd.Flags().StringVar(&o.Backend, "backend", "ssm", "backend type")
	cmd.Flags().StringVar(&o.SecretsPath, "path", "", "path to secrets")

	return cmd
}

func (o *SecretsRemoveOptions) Complete(cmd *cobra.Command, args []string) error {
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Env = viper.GetString("env")
	o.Profile = viper.GetString("aws-profile")
	o.Region = viper.GetString("aws-region")
	o.AppName = cmd.Flags().Args()[0]

	if o.Profile == "" {
		o.Profile = viper.GetString("aws_profile")
	}

	if o.Region == "" {
		o.Region = viper.GetString("aws_region")
	}

	if o.SecretsPath == "" {
		o.SecretsPath = fmt.Sprintf("/%s/%s", o.Env, o.AppName)
	}

	return nil
}

func (o *SecretsRemoveOptions) Validate() error {
	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("AWS region must be specified")
	}
	return nil
}

func (o *SecretsRemoveOptions) Run() error {
	pterm.DefaultSection.Printfln("Removing Secrets for %s", o.AppName)

	if o.Backend == "ssm" {
		err := rm(
			utils.SessionConfig{
				Region:  o.Region,
				Profile: o.Profile,
			},
			o,
		)
		if err != nil {
			pterm.DefaultSection.Sprintfln("Secrets have been removed from %s", o.SecretsPath)
			return err
		}
	} else {
		pterm.DefaultSection.Println("Secrets removal unsuccessful")
		return fmt.Errorf("backend %s is not found or not supported", o.Backend)
	}

	//pterm.DefaultSection.Printfln("Done: removing secrets")


	return nil
}

func rm(sessCfg utils.SessionConfig, o *SecretsRemoveOptions) error {
	if o.SecretsPath == "" {
		pterm.Info.Printfln("Path was not set")
		return nil
	}

	pterm.Info.Printfln("Removing secrets from %s://%s", o.Backend, o.SecretsPath)

	sess, err := utils.GetSession(&sessCfg)
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Establish AWS session")

	ssmSvc := ssm.New(sess)

	out, err := ssmSvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: &o.SecretsPath,
	})
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Getting secrets")

	if len(out.Parameters) == 0 {
		pterm.Info.Printfln("No values found")
		pterm.Success.Printfln("Removing secrets")
		return nil
	}

	var names []*string

	for _, p := range out.Parameters {
		names = append(names, p.Name)
	}

	_, err = ssmSvc.DeleteParameters(&ssm.DeleteParametersInput{
		Names: names,
	})

	ssmSvc.RemoveTagsFromResource(&ssm.RemoveTagsFromResourceInput{})

	if err != nil {
		return err
	}

	pterm.Success.Printfln("Removed secrets")

	return nil
}
