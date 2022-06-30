package validate

import (
	"github.com/hazelops/ize/internal/schema"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate configuration (only for test)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := schema.Validate(viper.AllSettings())

			if err != nil {
				return err
			}

			pterm.Success.Println("Config structure, env vars and flags look valid. ")

			return nil
		},
	}

	return cmd
}
