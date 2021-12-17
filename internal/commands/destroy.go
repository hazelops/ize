package commands

import (
	"fmt"

	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/mitchellh/mapstructure"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type destroyCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newDestroyCmd() *destroyCmd {
	cc := &destroyCmd{}

	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy anything",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "infra",
		Short: "destroy infrastructures",
		Long:  "Destroy infrastructures.",
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
						Cmd:           []string{"destroy", "-auto-approve"},
						Env: []string{
							fmt.Sprintf("ENV=%v", cc.config.Env),
							fmt.Sprintf("AWS_PROFILE=%v", tic.Profile),
							fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
							fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
						},
						TerraformVersion: tic.Version,
					}

					spinner := &pterm.SpinnerPrinter{}

					if cc.log.Level < 4 {
						spinner, _ = pterm.DefaultSpinner.Start("execution terraform destroy")
					}

					err = terraform.Run(&cc.log, opts)
					if err != nil {
						cc.log.Errorf("terraform %s not completed", "destroy")
						return err
					}

					if cc.log.Level < 4 {
						spinner.Success("terrafrom destroy completed")
					} else {
						pterm.Success.Println("terrafrom destroy completed")
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
