package commands

import (
	"github.com/hazelops/ize/internal/config"
	"github.com/spf13/cobra"
)

func NewCmdGen(project *config.Project) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate something",
	}

	cmd.AddCommand(
		NewCmdCI(project),
		NewCmdDoc(),
		NewCmdTfenv(project),
		NewCmdCompletion(),
		NewCmdMfa(project),
		NewCmdAWSProfile(),
	)

	return cmd
}
