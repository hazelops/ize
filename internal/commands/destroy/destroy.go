package destroy

import (
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
)

func NewCmdDestroy(ui terminal.UI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy anything",
	}

	cmd.AddCommand(NewCmdDestroyInfra(ui))

	return cmd
}
