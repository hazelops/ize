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

			cc.log.Infof("infra: %s", cc.config.Infra)

			for pname, provider := range cc.config.Infra {
				switch pname {
				case "terraform":
					var terraform map[string]hclTerraform = make(map[string]hclTerraform)
					var t hclTerraform
					mapstructure.Decode(provider, &t)
					terraform["main"] = t

					for sname, subitem := range provider {
						i, ok := subitem.(map[string]interface{})
						if ok {
							mapstructure.Decode(i, &t)
							terraform[sname] = t
						}
					}

					cc.log.Debugf("%s block: %s", pname, terraform)
				default:
					return fmt.Errorf("provider %s is not supported", pname)
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
