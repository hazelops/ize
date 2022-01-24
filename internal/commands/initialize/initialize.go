package initialize

import (
	"fmt"

	_ "embed"

	"github.com/hazelops/ize/internal/generate"
	"github.com/spf13/cobra"
)

type InitOptions struct {
	Path   string
	Output string
}

func NewInitFlags() *InitOptions {
	return &InitOptions{}
}

func NewCmdInit() *cobra.Command {
	o := NewInitFlags()

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates an IZE configuration file.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
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

	cmd.Flags().StringVar(&o.Path, "path", "", "set path to template (url)")
	cmd.Flags().StringVar(&o.Output, "output", "", "set output dir")
	cmd.MarkFlagRequired("path")

	return cmd
}

func (o *InitOptions) Complete(cmd *cobra.Command, args []string) error {
	if o.Path == "" {
		o.Path = "./ize.toml"
	}

	return nil
}

func (o *InitOptions) Validate() error {
	if len(o.Path) == 0 {
		return fmt.Errorf("path must be specified")
	}

	return nil
}

func (o *InitOptions) Run() error {
	_, err := generate.GenerateFiles(o.Path, o.Output)
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
