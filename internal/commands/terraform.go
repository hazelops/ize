package commands

import (
	"fmt"

	"github.com/hazelops/ize/internal/docker/terraform"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type terraformCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newTerraformCmd() *terraformCmd {
	cc := &terraformCmd{}

	cmd := &cobra.Command{
		Use:                "terraform <terraform command> [terraform flags]",
		Short:              "terraform management",
		Long:               "This command contains subcommands for work with terraform.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			if len(args) != 0 {
				if args[0] == "-h" || args[0] == "--help" {
					return cmd.Help()
				}
			}

			err := cc.Init()
			if err != nil {
				return err
			}

			opts := terraform.Options{
				ContainerName: "terraform",
				Cmd:           args,
				Env: []string{
					fmt.Sprintf("ENV=%v", cc.config.Env),
					fmt.Sprintf("AWS_PROFILE=%v", cc.config.AwsProfile),
					fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
					fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
				},
				TerraformVersion: cc.config.TerraformVersion,
			}

			cc.log.Debug("starting terraform")

			err = terraform.Run(&cc.log, opts)
			if err != nil {
				cc.log.Errorf("terraform %s not completed", args[0])
				return err
			}

			pterm.DefaultSection.Printfln("Terraform %s completed", args[0])

			return nil
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}
