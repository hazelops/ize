package gen

import (
	"bufio"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/generate"
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

			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}

			if o.Template == "" {
				return fmt.Errorf("'--template' must be specified")
			}

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
				Env:       cfg.Env,
				AwsRegion: cfg.AwsRegion,
				Apps:      cfg.GetApps(),
				Namespace: cfg.Namespace,
				PublicKey: getPublicKey(fmt.Sprintf("%s/.ssh/id_rsa.pub", cfg.Home)),
			})
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
