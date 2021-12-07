package commands

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type gendocCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newGendocCmd() *gendocCmd {
	cc := &gendocCmd{}

	cmd := &cobra.Command{
		Use:   "gendoc",
		Short: "Create Docs",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := doc.GenMarkdownTree(rootCmd, "./commands")
			if err != nil {
				return err
			}

			pterm.Success.Printfln("Docs generated")

			return nil
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}
