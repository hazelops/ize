package ecs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/hazelops/ize/internal/docker"
	"github.com/hazelops/ize/pkg/templates"
	"github.com/hazelops/ize/pkg/terminal"
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

	if e.App.Timeout == 0 {
		e.App.Timeout = 300
	}
}

// Deploy deploys app container to ECS via ECS deploy
func (e *Manager) Deploy(ui terminal.UI) error {
	e.prepare()

	if e.App.Unsafe && e.Project.PreferRuntime == "native" {
		pterm.Warning.Println(templates.Dedent(`
			deployment will be accelerated (unsafe):
			- Health Check Interval: 5s
			- Health Check Timeout: 2s
			- Healthy Threshold Count: 2
			- Unhealthy Threshold Count: 2`))
	}

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deploying app container...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	if e.App.Image == "" {
		e.App.Image = fmt.Sprintf("%s/%s:%s",
			e.Project.DockerRegistry,
			fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name),
			fmt.Sprintf("%s-%s", e.Project.Env, "latest"))
	}

	if e.Project.PreferRuntime == "native" {
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
	s = sg.Add("%s: deployment completed!", e.App.Name)
	s.Done()

	return nil
}

func (e *Manager) Redeploy(ui terminal.UI) error {
	e.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: redeploying app container...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	if e.Project.PreferRuntime == "native" {
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
	s = sg.Add("%s: redeployment completed!", e.App.Name)
	s.Done()

	return nil
}

func (e *Manager) Push(ui terminal.UI) error {
	e.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: push app image...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	image := fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name)

	svc := ecr.New(e.Project.Session)

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

	dockerRegistry := e.Project.DockerRegistry
	imageUri := fmt.Sprintf("%s/%s", dockerRegistry, image)

	r := docker.NewRegistry(*repository.RepositoryUri, token)

	err = r.Push(context.Background(), s.TermOutput(), imageUri, []string{e.Project.Tag, tagLatest})
	if err != nil {
		return fmt.Errorf("can't push image: %w", err)
	}

	s.Done()

	return nil
}

func (e *Manager) Build(ui terminal.UI) error {
	e.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: building app container...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	registry := e.Project.DockerRegistry
	image := fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name)
	imageUri := fmt.Sprintf("%s/%s", registry, image)

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

	b := docker.NewBuilder(
		buildArgs,
		tags,
		dockerfile,
		cache,
	)

	err = b.Build(ui, s, e.Project.RootDir)
	if err != nil {
		return fmt.Errorf("unable to build image: %w", err)
	}

	s.Done()

	return nil
}

func (e *Manager) Destroy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	ui.Output("Destroying ECS applications requires destroying the infrastructure.", terminal.WithWarningStyle())
	time.Sleep(time.Millisecond * 100)

	s := sg.Add("%s: destroying completed!", e.App.Name)
	defer func() { s.Abort() }()
	s.Done()

	return nil
}
