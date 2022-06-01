package push

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PushOptions struct {
	Config  *config.Config
	AppName string
	Tag     string
	App     interface{}
}

var pushLongDesc = templates.LongDesc(`
	Push app image (so far only ECR).
    App name must be specified for a app image push. 
`)

var pushExample = templates.Examples(`
	# Push image app (config file required)
	ize push <app name>

	# Push image app via config file
	ize --config-file (or -c) /path/to/config push <app name>

	# Push image app via config file installed from env
	export IZE_CONFIG_FILE=/path/to/config
	ize push <app name>
`)

func NewPushFlags() *PushOptions {
	return &PushOptions{}
}

func NewCmdPush() *cobra.Command {
	o := NewPushFlags()

	cmd := &cobra.Command{
		Use:     "push [flags] <app name>",
		Example: pushExample,
		Short:   "push app image",
		Long:    pushLongDesc,
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

func (o *PushOptions) Complete(cmd *cobra.Command, args []string) error {
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

func (o *PushOptions) Validate() error {

	return nil
}

func (o *PushOptions) Run() error {
	ui := terminal.ConsoleUI(context.Background(), o.Config.IsPlainText)
	sg := ui.StepGroup()
	defer sg.Wait()

	image := fmt.Sprintf("%s-%s", viper.GetString("namespace"), o.AppName)

	svc := ecr.New(o.Config.Session)

	var repository *ecr.Repository

	dro, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(image)},
	})
	if err != nil {
		return fmt.Errorf("can't describe repositories: %w", err)
	}

	if dro == nil || len(dro.Repositories) == 0 {
		logrus.Info("no ECR repository detected, creating", "name", image)

		out, err := svc.CreateRepository(&ecr.CreateRepositoryInput{
			RepositoryName: aws.String(image),
		})
		if err != nil {
			return fmt.Errorf("unable to create repository: %w", err)
		}

		repository = out.Repository
	} else {
		repository = dro.Repositories[0]
	}

	gat, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return fmt.Errorf("unable to get authorization token: %w", err)
	}

	if len(gat.AuthorizationData) == 0 {
		return fmt.Errorf("no authorization tokens provided")
	}

	uptoken := *gat.AuthorizationData[0].AuthorizationToken
	data, err := base64.StdEncoding.DecodeString(uptoken)
	if err != nil {
		return fmt.Errorf("unable to decode authorization token: %w", err)
	}

	auth := types.AuthConfig{
		Username: "AWS",
		Password: string(data[4:]),
	}

	authBytes, _ := json.Marshal(auth)

	token := base64.URLEncoding.EncodeToString(authBytes)

	s := sg.Add("%s: building app container...", o.AppName)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	tagLatest := fmt.Sprintf("%s-latest", o.Config.Env)

	dockerRegistry := viper.GetString("DOCKER_REGISTRY")
	imageUri := fmt.Sprintf("%s/%s", dockerRegistry, image)

	r := docker.NewRegistry(*repository.RepositoryUri, token)
	err = r.Push(context.Background(), ui, imageUri, []string{o.Tag, tagLatest})
	if err != nil {
		return fmt.Errorf("can't push image: %w", err)
	}

	s.Done()

	return nil
}
