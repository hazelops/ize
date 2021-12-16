package commands

import (
	"fmt"

	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
					var tic terraformInfraConfig
					mapstructure.Decode(provider, &tic)

					opts := terraform.Options{
						ContainerName: "terraform",
						Cmd:           []string{"init"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
					}

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "init")
						return err
					}

					opts = terraform.Options{
						ContainerName: "terraform",
						Cmd:           []string{"plan"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
					}

					cc.log.Debugf("terraform run opts: %s", opts)

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "plan")
						return err
					}

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

type terraformInfraConfig struct {
	RootDir string `mapstructure:"root_dir,optional"`
	Version string `mapstructure:"terraform_version,optional"`
	Region  string `mapstructure:"aws_region,optional"`
	Profile string `mapstructure:"aws_profile,optional"`
}
