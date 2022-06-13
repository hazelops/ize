package initialize

import (
	"fmt"

	"github.com/hazelops/ize/examples"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	Output   string
	Template string
	ShowList bool
}

var initLongDesc = templates.LongDesc(`
	Initialize a new ize project from a template
`)

var initExample = templates.Examples(`
	# Init project from url
	ize init --template https://github.com/<org>/<repo>

	# Init project from internal examples (https://github.com/hazelops/ize/tree/main/examples)
	ize init --template simple-monorepo

	# Display all internal templates
	ize init --list
`)

func NewInitFlags() *InitOptions {
	return &InitOptions{}
}

func NewCmdInit() *cobra.Command {
	o := NewInitFlags()

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize project",
		Long:    initLongDesc,
		Example: initExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if o.ShowList {
				internalTemplates()
				return nil
			}

			err := o.Validate(cmd)
			if err != nil {
				return err
			}

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.Template, "template", "", "set template (url or internal)")
	cmd.Flags().StringVar(&o.Output, "output", "", "set output dir")
	cmd.Flags().BoolVar(&o.ShowList, "list", false, "show list of examples")

	return cmd
}

func (o *InitOptions) Validate(cmd *cobra.Command) error {
	if len(o.Template) == 0 {
		cmd.Help()
		return fmt.Errorf("template must be specified\n")
	}

	return nil
}

func (o *InitOptions) Run() error {
	dest, err := generate.GenerateFiles(o.Template, o.Output)
	if err != nil {
		return err
	}

	pterm.Success.Printfln(`Initialized project from template "%s" to %s`, o.Template, dest)

	return nil
}

type ConfigOpts struct {
	Env               string
	Aws_profile       string
	Aws_region        string
	Terraform_version string
	Namespace         string
}

func internalTemplates() {
	dirs, err := examples.Examples.ReadDir(".")
	if err != nil {
		logrus.Fatal(err)
	}

	for _, d := range dirs {
		if d.IsDir() {
			fmt.Println(d.Name())
		}
	}
}
