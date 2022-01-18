package secret

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type SecretRemoveOptions struct {
	Config *config.Config
	Type   string
	Path   string
}

func NewSecretRemoveFlags() *SecretRemoveOptions {
	return &SecretRemoveOptions{}
}

func NewCmdSecretRemove() *cobra.Command {
	o := NewSecretRemoveFlags()

	cmd := &cobra.Command{
		Use:              "remove",
		Short:            "remove secrets from storage",
		Long:             "This command remove sercrets from storage",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
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

	cmd.Flags().StringVar(&o.Type, "type", "", "vault type")
	cmd.Flags().StringVar(&o.Path, "path", "", "path to secrets")
	cmd.MarkFlagRequired("type")
	cmd.MarkFlagRequired("path")

	return cmd
}

func (o *SecretRemoveOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *SecretRemoveOptions) Run() error {
	pterm.DefaultSection.Printfln("Starting remove secrets")

	if o.Type == "ssm" {
		err := remove(
			utils.SessionConfig{
				Region:  o.Config.AwsRegion,
				Profile: o.Config.AwsRegion,
			},
			o.Path,
		)
		if err != nil {
			pterm.DefaultSection.Println("Remove secrets not completed")
			return err
		}
	} else {
		pterm.DefaultSection.Println("Remove secrets not completed")
		return fmt.Errorf("vault with type %s not found or not supported", o.Type)
	}

	pterm.DefaultSection.Printfln("Remove secrets completed")

	return nil
}

func remove(sessCfg utils.SessionConfig, path string) error {
	if path == "" {
		pterm.Info.Printfln("Path were not set")
		return nil
	}

	sess, err := utils.GetSession(&sessCfg)
	if err != nil {
		return err
	}
	pterm.Success.Printfln("Geting AWS session")

	ssmSvc := ssm.New(sess)

	out, err := ssmSvc.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path: &path,
	})
	if err != nil {
		return err
	}

	pterm.Success.Printfln("Getting secrets from path")

	if len(out.Parameters) == 0 {
		pterm.Info.Printfln("No values found along the path")
		pterm.Success.Printfln("Deleting secrets from path")
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

	pterm.Success.Printfln("Deleting secrets from path")

	return nil
}
