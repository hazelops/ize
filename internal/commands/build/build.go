package build

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type BuildOptions struct {
	Config  *config.Config
	AppName string
	Tag     string
	App     interface{}
}

var buildLongDesc = templates.LongDesc(`
	Build sevice.
    App name must be specified for a app build. 
`)

var buildExample = templates.Examples(`
	# Build app (config file required)
	ize build <app name>

	# Build app via config file
	ize --config-file (or -c) /path/to/config build <app name>

	# Build app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize build <app name>
`)

func NewBuildFlags() *BuildOptions {
	return &BuildOptions{}
}

func NewCmdBuild() *cobra.Command {
	o := NewBuildFlags()

	cmd := &cobra.Command{
		Use:     "build [flags] <app name>",
		Example: buildExample,
		Short:   "manage builds",
		Long:    buildLongDesc,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

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

	return cmd
}

func (o *BuildOptions) Complete(cmd *cobra.Command, args []string) error {
	var err error
	o.Config, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("can`t complete options: %w", err)
	}

	viper.BindPFlags(cmd.Flags())
	o.AppName = cmd.Flags().Args()[0]
	viper.UnmarshalKey(fmt.Sprintf("app.%s", o.AppName), &o.App)

	o.Tag = viper.GetString("tag")

	return nil
}

func (o *BuildOptions) Validate() error {

	return nil
}

func (o *BuildOptions) Run() error {
	ui := terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: building app container...", o.AppName)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	projectPath := o.App.(map[string]interface{})["path"].(string)
	if projectPath == "" {
		return fmt.Errorf("can't build image: path is not specifed")
	}

	registry := viper.GetString("DOCKER_REGISTRY")
	image := fmt.Sprintf("%s-%s", viper.GetString("namespace"), o.AppName)
	imageUri := fmt.Sprintf("%s/%s", registry, image)

	buildArgs := map[string]*string{
		"PROJECT_PATH": &projectPath,
		"APP_NAME":     &o.AppName,
	}

	tags := []string{
		image,
		fmt.Sprintf("%s:%s", imageUri, o.Tag),
		fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", viper.GetString("ENV"))),
	}

	dockerfile := path.Join(projectPath, "Dockerfile")

	cache := []string{fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", viper.GetString("ENV")))}

	b := docker.NewBuilder(
		buildArgs,
		tags,
		dockerfile,
		cache,
	)

	b.Build(ui, s)

	s.Done()

	return nil
}
