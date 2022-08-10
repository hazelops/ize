package gen

import (
	"bufio"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type CIOptions struct {
	Template string
	Source   string
	Config   *config.Project
}

var ciLongDesc = templates.LongDesc(`
	Generate CI workflow.
    Template file and source url must be specified for a CI workflow generate. 
`)

var ciExample = templates.Examples(`
	# Generate CI workflow
	ize gen ci --template github.tmpl --source https://github.com/hazelops/ize-ci-templates
`)

func NewCIOptions() *CIOptions {
	return &CIOptions{}
}

func NewCmdCI() *cobra.Command {
	o := NewCIOptions()

	cmd := &cobra.Command{
		Use:     "ci",
		Short:   "Generate CI workflow",
		Long:    ciLongDesc,
		Example: ciExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete()
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run(cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&o.Template, "template", "", "set template path")
	cmd.Flags().StringVar(&o.Source, "source", "", "set git repository")

	return cmd
}

func (o *CIOptions) Complete() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *CIOptions) Validate() error {
	if o.Template == "" {
		return fmt.Errorf("'--template' must be specified")
	}

	if o.Source == "" {
		return fmt.Errorf("'--source' must be specified")
	}

	return nil
}

func (o *CIOptions) Run(cmd *cobra.Command) error {
	cmd.SilenceUsage = true

	file, err := generate.GetDataFromFile(o.Source, o.Template)
	if err != nil {
		return err
	}

	t := template.New("template")
	t, err = t.Parse(string(file))
	if err != nil {
		return err
	}

	err = t.Execute(os.Stdout, struct {
		Env       string
		AwsRegion string
		PublicKey string
		Namespace string
		Apps      map[string]*interface{}
	}{
		Env:       o.Config.Env,
		AwsRegion: o.Config.AwsRegion,
		Apps:      o.Config.GetApps(),
		Namespace: o.Config.Namespace,
		PublicKey: getPublicKey(fmt.Sprintf("%s/.ssh/id_rsa.pub", o.Config.Home)),
	})
	if err != nil {
		return err
	}

	return nil
}

func getPublicKey(path string) string {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	if _, err := os.Stat(path); err != nil {
		logrus.Fatalf("%s does not exist", path)
	}

	var key string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return key
}
