package initialize

import (
	"fmt"

	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	Output   string
	Template string
}

var initLongDesc = templates.LongDesc(`
	Initialize project from git or internal examples.
`)

var initExample = templates.Examples(`
	# Init project from url
	ize init --template https://github.com/<org>/<repo>

	# Init project from internal examples (https://github.com/hazelops/ize/tree/main/examples)
	ize init --template simple-monorepo


`)

func NewInitFlags() *InitOptions {
	return &InitOptions{}
}

func NewCmdInit() *cobra.Command {
	o := NewInitFlags()

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "initialize project",
		Long:    initLongDesc,
		Example: initExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			err := o.Validate()
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
	cmd.MarkFlagRequired("template")

	return cmd
}

func (o *InitOptions) Validate() error {
	if len(o.Template) == 0 {
		return fmt.Errorf("template must be specified")
	}

	return nil
}

func (o *InitOptions) Run() error {
	_, err := generate.GenerateFiles(o.Template, o.Output)
	if err != nil {
		return err
	}

	return nil
}

type ConfigOpts struct {
	Env               string
	Aws_profile       string
	Aws_region        string
	Terraform_version string
	Namespace         string
}
