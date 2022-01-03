package destroy

import (
	"github.com/spf13/cobra"
)

func NewCmdDestroy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy anything",
	}

	cmd.AddCommand(NewCmdDestroyInfra())

	return cmd
}
