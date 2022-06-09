package secrets

import (
	"github.com/spf13/cobra"
)

func NewCmdSecrets() *cobra.Command {

	cmd := &cobra.Command{
		Use:              "secrets",
		Short:            "Manage secrets",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSecretsRemove(),
		NewCmdSecretsPush(),
		NewCmdSecretsEdit(),
		NewCmdSecretsPull(),
	)

	return cmd
}
