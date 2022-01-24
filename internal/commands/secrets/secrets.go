package secrets

import (
	"github.com/spf13/cobra"
)

func NewCmdSecrets() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "secrets",
		Short:            "manage secrets",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSecretsRemove(),
		NewCmdSecretsPush(),
		NewCmdSecretsEdit(),
	)

	return cmd
}
