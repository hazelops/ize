package secrets

import (
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

func NewCmdSecrets(ui terminal.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "secrets",
		Short:            "manage secrets",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSecretsRemove(ui),
		NewCmdSecretsPush(ui),
		NewCmdSecretsEdit(),
	)

	return cmd
}
