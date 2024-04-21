package commands

import (
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CIOptions struct {
	Template      string
	Source        string
	Config        *config.Project
	PublicKeyFile string
}

var ciLongDesc = templates.LongDesc(`
	Generate CI workflow.
    Template file and source url must be specified for a CI workflow generate. 
`)

var ciExample = templates.Examples(`
	# Generate CI workflow
	ize gen ci --template github.tmpl --source https://github.com/hazelops/ize-ci-templates
`)

func NewCIOptions(project *config.Project) *CIOptions {
	return &CIOptions{
		Config: project,
	}
}

func NewCmdCI(project *config.Project) *cobra.Command {
	o := NewCIOptions(project)

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
	cmd.Flags().StringVar(&o.PublicKeyFile, "ssh-public-key", "", "set ssh key public path")

	return cmd
}

func (o *CIOptions) Complete() error {
	if len(o.PublicKeyFile) == 0 {
		o.PublicKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", o.Config.Home)
	}

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

	key, err := getPublicKey(o.PublicKeyFile)
	if err != nil {
		return err
	}

	err = t.Execute(os.Stdout, struct {
		config.Project
		PublicKey string
		Apps      map[string]*interface{}
	}{
		Project:   *o.Config,
		Apps:      o.Config.GetApps(),
		PublicKey: key,
	})
	if err != nil {
		return err
	}

	return nil
}

type Template struct {
	Path    string
	FuncMap template.FuncMap
	Data    interface{}
}

func (t *Template) Execute(dir string) error {
	isOnlyWhitespace := func(buf []byte) bool {
		wsre := regexp.MustCompile(`\S`)

		return !wsre.Match(buf)
	}

	err := filepath.Walk(t.Path, func(filename string, info fs.FileInfo, err error) error {
		oldName, err := filepath.Rel(t.Path, filename)
		if err != nil {
			return err
		}

		buf := bytes.NewBufferString("")
		fnameTmpl := template.Must(template.
			New("file name template").
			Funcs(t.FuncMap).
			Parse(oldName))

		if err := fnameTmpl.Execute(buf, t.Data); err != nil {
			return err
		}

		newName := buf.String()

		target := filepath.Join(dir, newName)

		if info.IsDir() {
			if err := os.Mkdir(target, 0755); err != nil {
				if !os.IsExist(err) {
					return err
				}
			}
		} else {
			fi, err := os.Lstat(filename)
			if err != nil {
				return err
			}

			// Delete target file if it exists
			if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
				return err
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, fi.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			defer func(fname string) {
				contents, err := ioutil.ReadFile(fname)
				if err != nil {
					logrus.Debug(fmt.Sprintf("couldn't read the contents of file %q, got error %q", fname, err))
					return
				}

				if isOnlyWhitespace(contents) {
					os.Remove(fname)
					return
				}
			}(f.Name())

			contentsTmpl := template.Must(template.
				New("file contents template").
				Funcs(t.FuncMap).
				ParseFiles(filename))

			fileTemplateName := filepath.Base(filename)

			if err := contentsTmpl.ExecuteTemplate(f, fileTemplateName, t.Data); err != nil {
				return err
			}
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
