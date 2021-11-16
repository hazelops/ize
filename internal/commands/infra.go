package commands

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/spf13/cobra"
)

type infraCmd struct {
	*baseBuilderCmd

	filePath string
}

func (b *commandsBuilder) newInfraCmd() *infraCmd {
	cc := &infraCmd{}

	cmd := &cobra.Command{
		Use:   "infra",
		Short: "",
		Long:  "",
	}

	cmd.AddCommand(&cobra.Command{
		Use: "deploy",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			for _, i := range cc.config.Infra {
				switch i.Provider {
				case "terraform":
					var terrafromConfig hclTerraform
					diag := gohcl.DecodeBody(i.Body, &hcl.EvalContext{}, &terrafromConfig)
					if diag.HasErrors() {
						return diag
					}

					if terrafromConfig.Profile == "" {
						terrafromConfig.Profile = cc.config.AwsProfile
					}
					if terrafromConfig.Region == "" {
						terrafromConfig.Region = cc.config.AwsRegion
					}
					if terrafromConfig.Version == "" {
						terrafromConfig.Version = cc.config.TerraformVersion
					}

					fmt.Println(terrafromConfig)
				default:
					return fmt.Errorf("provider %s is not supported", i.Provider)
				}
			}

			return nil
		},
	})

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

type hclTerraform struct {
	RootDir string `hcl:"root_dir,optional"`
	Version string `hcl:"terraform_version,optional"`
	Region  string `hcl:"aws_region,optional"`
	Profile string `hcl:"aws_profile,optional"`
}
