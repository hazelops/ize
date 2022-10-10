package ecs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/templates"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/hazelops/ize/internal/docker"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

const ecsDeployImage = "hazelops/ecs-deploy:latest"

type Manager struct {
	Project *config.Project
	App     *config.Ecs
}

func (e *Manager) prepare() {
	if e.App.Path == "" {
		appsPath := e.Project.AppsPath
		if !filepath.IsAbs(appsPath) {
			appsPath = filepath.Join(os.Getenv("PWD"), appsPath)
		}

		e.App.Path = filepath.Join(appsPath, e.App.Name)
	} else {
		rootDir := e.Project.RootDir

		if !filepath.IsAbs(e.App.Path) {
			e.App.Path = filepath.Join(rootDir, e.App.Path)
		}
	}

	if len(e.App.Cluster) == 0 {
		e.App.Cluster = fmt.Sprintf("%s-%s", e.Project.Env, e.Project.Namespace)
	}

	if len(e.App.DockerRegistry) == 0 {
		e.App.DockerRegistry = e.Project.DockerRegistry
	}

	if e.App.Timeout == 0 {
		e.App.Timeout = 300
	}
}

// Deploy deploys app container to ECS via ECS deploy
func (e *Manager) Deploy(ui *echelon.Logger) error {
	e.prepare()

	if len(e.App.AwsRegion) != 0 && len(e.App.AwsProfile) != 0 {
		sess, err := utils.GetSession(&utils.SessionConfig{
			Region:  e.App.AwsRegion,
			Profile: e.App.AwsProfile,
		})
		if err != nil {
			return fmt.Errorf("can't get session: %w", err)
		}

		e.Project.SettingAWSClient(sess)
	}

	if e.App.SkipDeploy {
		s := ui.Scoped(fmt.Sprintf("%s: deploy will be skipped", e.App.Name))
		s.Finish(true)
		return nil
	}

	if e.App.Unsafe && e.Project.PreferRuntime == "native" {
		pterm.Warning.Println(templates.Dedent(`
			deployment will be accelerated (unsafe):
			- Health Check Interval: 5s
			- Health Check Timeout: 2s
			- Healthy Threshold Count: 2
			- Unhealthy Threshold Count: 2`))
	}

	s := ui.Scoped(fmt.Sprintf("%s: deploying app container...", e.App.Name))
	defer s.Finish(false)

	if e.App.Image == "" {
		e.App.Image = fmt.Sprintf("%s/%s:%s",
			e.App.DockerRegistry,
			fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name),
			fmt.Sprintf("%s-%s", e.Project.Env, "latest"))
	}

	if e.Project.PreferRuntime == "native" {
		err := e.deployLocal(s.AsWriter(echelon.InfoLevel))
		pterm.SetDefaultOutput(os.Stdout)
		if err != nil {
			return fmt.Errorf("unable to deploy app: %w", err)
		}
	} else {
		err := e.deployWithDocker(s.AsWriter(echelon.InfoLevel))
		if err != nil {
			return fmt.Errorf("unable to deploy app: %w", err)
		}
	}

	s.Finish(true)

	return nil
}

func (e *Manager) Redeploy(ui *echelon.Logger) error {
	e.prepare()

	if len(e.App.AwsRegion) != 0 && len(e.App.AwsProfile) != 0 {
		sess, err := utils.GetSession(&utils.SessionConfig{
			Region:  e.App.AwsRegion,
			Profile: e.App.AwsProfile,
		})
		if err != nil {
			return fmt.Errorf("can't get session: %w", err)
		}

		e.Project.SettingAWSClient(sess)
	}

	s := ui.Scoped(fmt.Sprintf("%s: redeploying app container...", e.App.Name))

	if e.Project.PreferRuntime == "native" {
		err := e.redeployLocal(s.AsWriter(echelon.InfoLevel))
		pterm.SetDefaultOutput(os.Stdout)
		if err != nil {
			return fmt.Errorf("unable to redeploy app: %w", err)
		}
	} else {
		err := e.redeployWithDocker(s.AsWriter(echelon.InfoLevel))
		if err != nil {
			return fmt.Errorf("unable to redeploy app: %w", err)
		}
	}

	s.Finish(true)

	return nil
}

func (e *Manager) Push(ui *echelon.Logger) error {
	e.prepare()

	if len(e.App.Image) != 0 {
		s := ui.Scoped(fmt.Sprintf("%s: pushing app image... (skipped, using %s) ", e.App.Name, e.App.Image))
		s.Finish(true)
		return nil
	}

	s := ui.Scoped(fmt.Sprintf("%s: push app image...", e.App.Name))

	image := fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name)

	svc := e.Project.AWSClient.ECRClient

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

	upToken := *gat.AuthorizationData[0].AuthorizationToken
	data, err := base64.StdEncoding.DecodeString(upToken)
	if err != nil {
		return fmt.Errorf("unable to decode authorization token: %w", err)
	}

	auth := types.AuthConfig{
		Username: "AWS",
		Password: string(data[4:]),
	}

	authBytes, _ := json.Marshal(auth)

	token := base64.URLEncoding.EncodeToString(authBytes)

	tagLatest := fmt.Sprintf("%s-latest", e.Project.Env)
	imageUri := fmt.Sprintf("%s/%s", e.App.DockerRegistry, image)
	platform := "linux/amd64"
	if e.Project.PreferRuntime == "docker-arm64" {
		platform = "linux/arm64"
	}

	r := docker.NewRegistry(*repository.RepositoryUri, token, platform)

	err = r.Push(context.Background(), s.AsWriter(echelon.InfoLevel), imageUri, []string{e.Project.Tag, tagLatest})
	if err != nil {
		return fmt.Errorf("can't push image: %w", err)
	}

	s.Finish(true)

	return nil
}

func (e *Manager) Build(ui *echelon.Logger) error {
	e.prepare()

	if len(e.App.Image) != 0 {
		s := ui.Scoped(fmt.Sprintf("%s: building app container... (skipped, using %s)", e.App.Name, e.App.Image))
		s.Finish(true)
		return nil
	}

	s := ui.Scoped(fmt.Sprintf("%s: building app container...", e.App.Name))

	image := fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name)
	imageUri := fmt.Sprintf("%s/%s", e.App.DockerRegistry, image)

	relProjectPath, err := filepath.Rel(e.Project.RootDir, e.App.Path)
	if err != nil {
		return fmt.Errorf("unable to get relative path: %w", err)
	}

	buildArgs := map[string]*string{
		"PROJECT_PATH": &relProjectPath,
		"APP_PATH":     &relProjectPath,
		"APP_NAME":     &e.App.Name,
	}

	tags := []string{
		image,
		fmt.Sprintf("%s:%s", imageUri, e.Project.Tag),
		fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", e.Project.Env)),
	}

	dockerfile := path.Join(e.App.Path, "Dockerfile")

	cache := []string{fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", e.Project.Env))}

	platform := "linux/amd64"
	if e.Project.PreferRuntime == "docker-arm64" {
		platform = "linux/arm64"
	}

	b := docker.NewBuilder(
		buildArgs,
		tags,
		dockerfile,
		cache,
		platform,
	)

	err = b.Build(s.AsWriter(echelon.InfoLevel), e.Project.RootDir)
	if err != nil {
		return fmt.Errorf("unable to build image: %w", err)
	}

	s.Finish(true)

	return nil
}

func (e *Manager) Destroy(ui *echelon.Logger) error {
	ui.Infof("Destroying ECS applications requires destroying the infrastructure.")
	time.Sleep(time.Millisecond * 50)

	return nil
}
