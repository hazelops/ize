package initialize

import (
	"fmt"
	"github.com/hazelops/ize/internal/version"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hazelops/ize/examples"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Use:     "init [flags] <path>",
		Short:   "Initialize project",
		Args:    cobra.MaximumNArgs(1),
		Long:    initLongDesc,
		Example: initExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if o.ShowList {
				internalTemplates()
				return nil
			}

			err := o.Complete(cmd)
			if err != nil {
				return err
			}

			err = o.Validate(cmd)
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
	cmd.Flags().BoolVar(&o.ShowList, "list", false, "show list of examples")

	return cmd
}

func (o *InitOptions) Complete(cmd *cobra.Command) error {
	o.Output = "."
	if len(cmd.Flags().Args()) != 0 {
		o.Output = cmd.Flags().Args()[0]
	}

	return nil
}

func (o *InitOptions) Validate(cmd *cobra.Command) error {
	return nil
}

func (o *InitOptions) Run() error {
	if len(o.Template) != 0 {
		dest, err := generate.GenerateFiles(o.Template, o.Output)
		if err != nil {
			return err
		}

		pterm.Success.Printfln(`Initialized project from template "%s" to %s`, o.Template, dest)

		return nil
	}

	namespace := ""
	envList := []string{}

	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "dev"
	}

	dir := o.Output
	if dir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("can't get current work directory: %s", cwd)
		}
	}

	dir, err := filepath.Abs(o.Output)
	if err != nil {
		return fmt.Errorf("can't init: %w", err)
	}

	namespace = filepath.Base(dir)
	err = survey.AskOne(
		&survey.Input{
			Message: fmt.Sprintf("Namespace:"),
			Default: namespace,
		},
		&namespace,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return fmt.Errorf("can't init: %w", err)
	}

	err = survey.AskOne(
		&survey.Input{
			Message: fmt.Sprintf("Environment:"),
			Default: env,
		},
		&env,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return fmt.Errorf("can't init: %w", err)
	}

	envList = append(envList, env)
	env = ""

	for {
		err = survey.AskOne(
			&survey.Input{
				Message: fmt.Sprintf("Another environment? [enter - skip]"),
				Default: env,
			},
			&env,
		)
		if err != nil {
			return fmt.Errorf("can't init: %w", err)
		}

		if env == "" {
			break
		}

		envList = append(envList, env)
		env = ""
	}

	for _, v := range envList {
		envPath := filepath.Join(dir, ".ize", "env", v)
		err := os.MkdirAll(envPath, 0755)
		if err != nil {
			return fmt.Errorf("can't create dir by path %s: %w", envPath, err)
		}

		viper.Reset()
		cfg := make(map[string]string)
		cfg["namespace"] = namespace
		cfg["env"] = v

		raw := make(map[string]interface{}, len(cfg))
		for k, v := range cfg {
			raw[k] = v
		}

		viper.MergeConfigMap(raw)
		err = viper.WriteConfigAs(filepath.Join(envPath, "ize.toml"))
		if err != nil {
			return fmt.Errorf("can't write config: %w", err)
		}
	}

	pterm.Success.Printfln(`Created ize skeleton for %s in %s`, strings.Join(envList, ", "), dir)
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

	var internal [][]string
	internal = append(internal, []string{"Name", "Version"})

	for _, d := range dirs {
		if d.IsDir() {
			internal = append(internal, []string{d.Name(), version.GitCommit})
		}
	}

	pterm.DefaultTable.WithHasHeader().WithData(internal).Render()

}
