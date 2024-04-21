package commands

import (
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdBoostrap(project *config.Project) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "boostrap",
		Short:            "Boostrap resources",
		TraverseChildren: true,
	}

	cmd.AddCommand(
		NewBoostrapTerraformState(project),
	)

	return cmd
}
