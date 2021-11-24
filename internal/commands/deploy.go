package commands

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	*baseBuilderCmd

	filePath string
}

func (b *commandsBuilder) newDeployCmd() *deployCmd {
	cc := &deployCmd{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Manage deployments.",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "infra",
		Short: "Deploy infrastructures.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			for pName, pBlock := range cc.config.Infra {

				switch pName {
				case "terraform":
					for _, block := range pBlock {

						var terrafromConfig hclTerraform
						mapstructure.Decode(block, &terrafromConfig)

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

					}
				default:
					return fmt.Errorf("provider %s is not supported", pName)
				}
			}
			return nil
		},
	})

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

type hclTerraform struct {
	RootDir string `mapstructure:"root_dir,optional"`
	Version string `mapstructure:"terraform_version,optional"`
	Region  string `mapstructure:"aws_region,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
