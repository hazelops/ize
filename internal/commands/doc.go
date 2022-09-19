package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"os"
)

func NewCmdDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "doc",
		Short:                 "Create docs",
		DisableFlagsInUseLine: true,
		Long:                  "Create docs with ize commands description",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := os.MkdirAll("./website/commands", 0777)
			if err != nil {
				return err
			}

			err = doc.GenMarkdownTree(cmd.Root(), "./website/commands")
			if err != nil {
				return err
			}

			pterm.Success.Printfln("Docs generated")

			return nil
		},
	}

	return cmd
}
