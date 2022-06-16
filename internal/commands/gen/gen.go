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
		Short: "Generate something",
	}

	cmd.AddCommand(
		NewCmdDoc(),
		NewCmdTfenv(),
		NewCmdCompletion(),
		NewCmdMfa(),
		NewCmdAWSProfile(),
	)

	return cmd
}
