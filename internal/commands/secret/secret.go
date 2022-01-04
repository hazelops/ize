package secret

import (
	"github.com/spf13/cobra"
)

func NewCmdSecret() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "secret",
		Short:            "manage secret",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSecretRemove(),
		NewCmdSecretSet(),
	)

	return cmd
}
