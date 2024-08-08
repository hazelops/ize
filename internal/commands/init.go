package commands

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/hazelops/ize/internal/schema"
	"github.com/hazelops/ize/internal/version"
	"golang.org/x/term"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hazelops/ize/examples"
	"github.com/hazelops/ize/internal/generate"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InitOptions struct {
	Output       string
	Template     string
	ShowList     bool
	SkipExamples bool
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

	cmd.Flags().StringVar(&o.Template, "template", "", "set template (url or internal)")
	cmd.Flags().BoolVar(&o.ShowList, "list", false, "show list of examples")
	cmd.Flags().BoolVar(&o.SkipExamples, "skip-examples", false, "generate ize.toml without commented examples")

	return cmd
}

func (o *InitOptions) Complete(cmd *cobra.Command) error {
	o.Output = "."
	if len(cmd.Flags().Args()) != 0 {
		o.Output = cmd.Flags().Args()[0]
	}

	return nil
}

func (o *InitOptions) Validate() error {
	return nil
}

func (o *InitOptions) Run() error {
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	if len(o.Template) != 0 {
		dest, err := generate.GenerateFiles(o.Template, o.Output)
		if err != nil {
			return err
		}

		pterm.Success.Printfln(`Initialized project from template "%s" to %s`, o.Template, dest)

		return nil
	}

	namespace := ""
	var envList []string

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

	if isTTY {
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
	} else {
		envList = append(envList, env)
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

		if o.SkipExamples {
			err := viper.MergeConfigMap(raw)
			if err != nil {
				return err
			}

			err = viper.WriteConfigAs(filepath.Join(envPath, "ize.toml"))
			if err != nil {
				return fmt.Errorf("can't write config: %w", err)
			}
		} else {
			err = writeConfig(filepath.Join(envPath, "ize.toml"), cfg)
			if err != nil {
				return fmt.Errorf("can't write config: %w", err)
			}
		}

		return nil
	}

	pterm.Success.Printfln(`Created ize skeleton for %s in %s`, strings.Join(envList, ", "), dir)
	return nil
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

	_ = pterm.DefaultTable.WithHasHeader().WithData(internal).Render()
}

func writeConfig(path string, existsValues map[string]string) error {
	allSettings := schema.GetJsonSchema()

	var str string

	str += getProperties(allSettings, existsValues)

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = f.WriteString(str)
	if err != nil {
		return err
	}

	return err
}

func getProperties(settings interface{}, existsValues map[string]string) string {
	var strRoot string
	var strBlocks string

	properties, ok := settings.(map[string]interface{})["properties"]
	if ok {
		propertiesMap, ok := properties.(map[string]interface{})
		if ok {
			for pn, pv := range propertiesMap {
				pm, ok := pv.(map[string]interface{})
				if ok {
					_, ok := pm["deprecationMessage"]
					if ok {
						continue
					}
					_, ok = pm["patternProperties"].(map[string]interface{})
					if ok {
						pd, ok := settings.(map[string]interface{})["definitions"].(map[string]interface{})[pn]
						desc := ""
						if ok {
							if pn == "terraform" {
								strBlocks += fmt.Sprintf("\n# [%s.infra]%s\n", pn, desc)
							} else {
								strBlocks += fmt.Sprintf("\n# [%s.<name>]%s\n", pn, desc)
							}
							strBlocks += getProperties(pd, map[string]string{})
						}
					}
					pt, ok := pm["type"].(string)
					if !ok {
						pt = "boolean"
					}
					pd := pm["description"].(string)
					switch pt {
					case "string":
						v, ok := existsValues[pn]
						if ok && pn != "env" {
							strRoot += fmt.Sprintf("%-36s\t# %s\n", fmt.Sprintf("%s = \"%s\"", pn, v), pd)
						} else {
							strRoot += fmt.Sprintf("# %-36s\t# %s\n", fmt.Sprintf("%s = \"%s\"", pn, v), pd)
						}
					case "boolean":
						strRoot += fmt.Sprintf("# %-36s\t# %s\n", fmt.Sprintf("%s = false", pn), pd)
					}
				}
			}
		}
	}

	lines := strings.Split(strRoot, "\n")
	sort.Sort(parametersSort(lines))
	strRoot = strings.Join(lines, "\n")

	strRoot += strBlocks

	return strRoot
}

type parametersSort []string

func (p parametersSort) Less(i, _ int) bool {
	if len(p[i]) == 0 {
		return false
	}
	return (p[i][0]) != '#' || strings.Contains(p[i], "required")
}
func (p parametersSort) Len() int      { return len(p) }
func (p parametersSort) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
