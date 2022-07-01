package gen

import "github.com/spf13/cobra"

func NewCmdGen() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate something",
	}

	cmd.AddCommand(
		NewCmdCI(),
		NewCmdDoc(),
		NewCmdTfenv(),
		NewCmdCompletion(),
		NewCmdMfa(),
		NewCmdAWSProfile(),
	)

	return cmd
}
