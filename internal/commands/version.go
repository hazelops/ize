package commands

import (
	"io"
	"text/template"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show IZE version",
		Run: func(cmd *cobra.Command, args []string) {
			c := cmd.Parent()
			tmpl(c.OutOrStdout(), c.VersionTemplate(), c)
		},
	}

	return cmd
}

func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
