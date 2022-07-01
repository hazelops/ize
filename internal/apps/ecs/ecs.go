package ecs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/docker"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/mitchellh/mapstructure"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ecsDeployImage = "hazelops/ecs-deploy:latest"

type ecs struct {
	Name                   string
	Unsafe                 bool
	Path                   string
	Image                  string
	Cluster                string
	TaskDefinitionRevision string `mapstructure:"task_definition_revision"`
	Timeout                int
	AwsProfile             string
	AwsRegion              string
	Tag                    string
}

func NewECSApp(name string, app interface{}) *ecs {
	ecsConfig := ecs{}

	raw, ok := app.(map[string]interface{})
	if ok {
		ecsConfig.Name = name
		mapstructure.Decode(raw, &ecsConfig)
	}

	ecsConfig.Name = name

	if ecsConfig.Path == "" {
		appsPath := viper.GetString("APPS_PATH")
		if !filepath.IsAbs(appsPath) {
			appsPath = filepath.Join(os.Getenv("PWD"), appsPath)
		}

		ecsConfig.Path = filepath.Join(appsPath, name)
	} else {
		rootDir := viper.GetString("ROOT_DIR")

		if !filepath.IsAbs(ecsConfig.Path) {
			ecsConfig.Path = filepath.Join(rootDir, ecsConfig.Path)
		}
	}

	ecsConfig.AwsProfile = viper.GetString("aws_profile")
	ecsConfig.AwsRegion = viper.GetString("aws_region")
	ecsConfig.Tag = viper.GetString("tag")
	if len(ecsConfig.Cluster) == 0 {
		ecsConfig.Cluster = fmt.Sprintf("%s-%s", viper.GetString("env"), viper.GetString("namespace"))
	}

	if ecsConfig.Timeout == 0 {
		ecsConfig.Timeout = 300
	}

	return &ecsConfig
}

// Deploy deploys app container to ECS via ECS deploy
func (e *ecs) Deploy(ui terminal.UI) error {
	if e.Unsafe && viper.GetString("prefer-runtime") == "native" {
		pterm.Warning.Println(templates.Dedent(`
			deployment will be accelerated (unsafe):
			- Health Check Interval: 5s
			- Health Check Timeout: 2s
			- Healthy Threshold Count: 2
			- Unhealthy Threshold Count: 2`))
	}

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deploying app container...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	if e.Image == "" {
		e.Image = fmt.Sprintf("%s/%s:%s",
			viper.GetString("DOCKER_REGISTRY"),
			fmt.Sprintf("%s-%s", viper.GetString("namespace"), e.Name),
			fmt.Sprintf("%s-%s", viper.GetString("env"), "latest"))
	}

	if viper.GetString("prefer-runtime") == "native" {
		err := e.deployLocal(s.TermOutput())
		pterm.SetDefaultOutput(os.Stdout)
		if err != nil {
			return fmt.Errorf("unable to deploy app: %w", err)
		}
	} else {
		err := e.deployWithDocker(s.TermOutput())
		if err != nil {
			return fmt.Errorf("unable to deploy app: %w", err)
		}
	}

	s.Done()
	s = sg.Add("%s: deployment completed!", e.Name)
	s.Done()

	return nil
}

func (e *ecs) Redeploy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: redeploying app container...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	if viper.GetString("prefer-runtime") == "native" {
		err := e.redeployLocal(s.TermOutput())
		pterm.SetDefaultOutput(os.Stdout)
		if err != nil {
			return fmt.Errorf("unable to redeploy app: %w", err)
		}
	} else {
		err := e.redeployWithDocker(s.TermOutput())
		if err != nil {
			return fmt.Errorf("unable to redeploy app: %w", err)
		}
	}

	s.Done()
	s = sg.Add("%s: redeployment completed!", e.Name)
	s.Done()

	return nil
}

func (e *ecs) Push(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: push app image...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	image := fmt.Sprintf("%s-%s", viper.GetString("namespace"), e.Name)

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  e.AwsRegion,
		Profile: e.AwsProfile,
	})
	if err != nil {
		return fmt.Errorf("unable to get aws session: %w", err)
	}

	svc := ecr.New(sess)

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

	tagLatest := fmt.Sprintf("%s-latest", viper.GetString("env"))

	dockerRegistry := viper.GetString("DOCKER_REGISTRY")
	imageUri := fmt.Sprintf("%s/%s", dockerRegistry, image)

	r := docker.NewRegistry(*repository.RepositoryUri, token)

	err = r.Push(context.Background(), s.TermOutput(), imageUri, []string{viper.GetString("tag"), tagLatest})
	if err != nil {
		return fmt.Errorf("can't push image: %w", err)
	}

	s.Done()

	return nil
}

func (e *ecs) Build(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: building app container...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	registry := viper.GetString("DOCKER_REGISTRY")
	image := fmt.Sprintf("%s-%s", viper.GetString("namespace"), e.Name)
	imageUri := fmt.Sprintf("%s/%s", registry, image)

	relProjectPath, err := filepath.Rel(viper.GetString("ROOT_DIR"), e.Path)
	if err != nil {
		return fmt.Errorf("unable to get relative path: %w", err)
	}

	buildArgs := map[string]*string{
		"PROJECT_PATH": &relProjectPath,
		"APP_NAME":     &e.Name,
	}

	tags := []string{
		image,
		fmt.Sprintf("%s:%s", imageUri, e.Tag),
		fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", viper.GetString("ENV"))),
	}

	dockerfile := path.Join(e.Path, "Dockerfile")

	cache := []string{fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", viper.GetString("ENV")))}

	b := docker.NewBuilder(
		buildArgs,
		tags,
		dockerfile,
		cache,
	)

	err = b.Build(ui, s)
	if err != nil {
		return fmt.Errorf("unable to build image: %w", err)
	}

	s.Done()

	return nil
}

func (e *ecs) Destroy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	ui.Output("Destroying ECS applications requires destroying the infrastructure.", terminal.WithWarningStyle())
	time.Sleep(time.Millisecond * 100)

	s := sg.Add("%s: destroying completed!", e.Name)
	defer func() { s.Abort() }()
	s.Done()

	return nil
}
