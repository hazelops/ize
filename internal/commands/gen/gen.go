package gen

import "github.com/spf13/cobra"

type GenOptions struct {
}

// func NewGenFlags() *GenOptions {
// 	return &GenOptions{}
// }

func NewCmdGen() *cobra.Command {
	// o := NewGenFlags()

	cmd := &cobra.Command{
		Use:   "gen",
		Short: "generates something",
	}

	cmd.AddCommand(
		NewCmdDoc(),
		NewCmdEnv(),
		NewCmdCompletion(),
		NewCmdMfa(),
	)

	return cmd
}
