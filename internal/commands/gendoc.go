package commands

import (
	"context"

	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewGendocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "gendoc",
		Short:                 "create docs",
		DisableFlagsInUseLine: true,
		Long:                  "Create docs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			root, err := newApp(terminal.ConsoleUI(context.Background()))
			if err != nil {
				return err
			}

			err = doc.GenMarkdownTree(root, "./commands")
			if err != nil {
				return err
			}

			pterm.Success.Printfln("Docs generated")

			return nil
		},
	}

	return cmd
}
