package gen

import (
	"fmt"
	"github.com/hazelops/ize/internal/generate"
	"github.com/spf13/cobra"
)

type CIOptions struct {
	Template string
	Source   string
}

func NewCIOptions() *CIOptions {
	return &CIOptions{}
}

func NewCmdCI() *cobra.Command {
	o := NewCIOptions()

	cmd := &cobra.Command{
		Use:   "ci",
		Short: "Generate CI workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if o.Template == "" {
				return fmt.Errorf("'--template' must be specified")
			}

			if o.Source == "" {
				o.Source = "."
			}

			file, err := generate.GetDataFromFile(o.Source, o.Template)
			if err != nil {
				return err
			}

			fmt.Print("template:\n", string(file))

			return nil
		},
	}

	cmd.Flags().StringVar(&o.Template, "template", "", "set template path")
	cmd.Flags().StringVar(&o.Source, "source", "", "set git repository")

	return cmd
}
