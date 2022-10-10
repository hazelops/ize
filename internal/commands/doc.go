package commands

import (
	"github.com/hazelops/ize/internal/schema"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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

			sections := []string{"main", "terraform", "ecs", "alias", "serverless", "tunnel"}

			err = os.MkdirAll("./website/schema", 0777)
			if err != nil {
				return err
			}

			t := template.New("schema")
			t, err = t.Parse(sectionTmpl)
			if err != nil {
				return err
			}

			for _, s := range sections {
				err := generateSection(s, t)
				if err != nil {
					return err
				}
			}

			pterm.Success.Printfln("Docs generated")

			return nil
		},
	}

	return cmd
}

func generateSection(name string, t *template.Template) error {
	s := schema.GetSchema()
	filename := "README.md"
	if name != "main" {
		filename = strings.ToUpper(name) + ".md"
		s = schema.GetSchema()[name].Items
	}
	f, err := os.Create(filepath.Join(".", "website", "schema", filename))
	if err != nil {
		return err
	}

	err = t.Execute(f, Section{
		Name:  cases.Title(language.Und).String(name),
		Items: s,
	})
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

type Section struct {
	Name string
	schema.Items
}

var sectionTmpl = `# {{.Name}} Section
| Parameter | Required | Description |
| --- | --- | --- |
{{range $k, $v := .Items}}| {{$k}} | {{if $v.Required}}yes{{else}}no {{end}}| {{$v.Description}} |
{{end}}
`
