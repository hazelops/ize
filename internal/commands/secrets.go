package commands

import (
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdSecrets(project *config.Project) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "secrets",
		Short:            "Manage secrets",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewCmdSecretsRemove(project),
		NewCmdSecretsPush(project),
		NewCmdSecretsEdit(project),
		NewCmdSecretsPull(project),
	)

	return cmd
}
