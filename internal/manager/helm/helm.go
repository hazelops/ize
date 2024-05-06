package helm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/hazelops/ize/pkg/term"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/docker"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
)

const helmDeployImage = "hazelops/helm-deploy:latest"

type Manager struct {
	Project *config.Project
	App     *config.Helm
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

	if len(e.App.Namespace) == 0 {
		e.App.Namespace = fmt.Sprintf("%s-%s", e.Project.Env, e.Project.Namespace)
	}

	if len(e.App.DockerRegistry) == 0 {
		e.App.DockerRegistry = e.Project.DockerRegistry
	}

	if e.App.Timeout == 0 {
		e.App.Timeout = 300
	}
}

// Deploy deploys app container to helm via helm deploy
func (e *Manager) Deploy(ui terminal.UI) error {
	e.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	if len(e.App.AwsRegion) != 0 && len(e.App.AwsProfile) != 0 {
		sess, err := utils.GetSession(&utils.SessionConfig{
			Region:      e.App.AwsRegion,
			Profile:     e.App.AwsProfile,
			EndpointUrl: e.Project.EndpointUrl,
		})
		if err != nil {
			return fmt.Errorf("can't get session: %w", err)
		}

		e.Project.SettingAWSClient(sess)
	}

	if e.App.SkipDeploy {
		s := sg.Add("%s: deploy will be skipped", e.App.Name)
		defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()
		s.Done()
		return nil
	}

	//if e.App.Unsafe && e.Project.PreferRuntime == "native" {
	//	pterm.Warning.Println(templates.Dedent(`
	//		deployment will be accelerated (unsafe):
	//		- Health Check Interval: 5s
	//		- Health Check Timeout: 2s
	//		- Healthy Threshold Count: 2
	//		- Unhealthy Threshold Count: 2`))
	//}

	s := sg.Add("%s: deploying app container...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	if e.App.Image == "" {
		e.App.Image = fmt.Sprintf("%s/%s:%s",
			e.App.DockerRegistry,
			fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name),
			fmt.Sprintf("%s-%s", e.Project.Env, "latest"))
	}

	if e.Project.PreferRuntime == "native" {

		//err := e.deployLocal(s.TermOutput())
		pterm.SetDefaultOutput(os.Stdout)
		s = sg.Add("%s: deploying app [run helm deploy]...", e.App.Name)
		err := e.runDeploy(s.TermOutput())

		if err != nil {
			return fmt.Errorf("unable to deploy app: %w", err)
		}
	} else {
		return fmt.Errorf("runtime not implemented. use native")
	}

	s.Done()
	s = sg.Add("%s: deployment completed!", e.App.Name)
	s.Done()

	return nil
}

func (helm *Manager) runDeploy(w io.Writer) error {

	// TODO build image name as a part of helm manager
	logrus.Debugf(helm.Project.Namespace)
	helmImageName := fmt.Sprintf("%s/%s-%s", helm.App.DockerRegistry, helm.Project.Namespace, helm.App.Name)
	helmChartDir := fmt.Sprintf("%s", path.Join(helm.App.Path, "helm/api"))

	var command string
	namespace := fmt.Sprintf("%s-%s", helm.Project.Env, helm.Project.Namespace)
	command = fmt.Sprintf(
		`AWS_PROFILE="localstack-user" KUBECONFIG=/Users/dmitry/.kube/metameetings-local \
		 helm upgrade --atomic --install --namespace "%s" "%s" %s \
			--set image.repository=%s \
			--set image.tag="%s"
		`, namespace, helm.App.Name, helmChartDir, helmImageName, helm.Project.Tag)

	//if sls.App.Force {
	//	command += ` \
	//	--force`
	//}

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	return term.New(
		term.WithDir(helm.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func (e *Manager) Push(ui terminal.UI) error {
	e.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: push app image...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	if len(e.App.Image) != 0 {
		s.Update("%s: pushing app image... (skipped, using %s) ", e.App.Name, e.App.Image)
		s.Done()

		return nil
	}

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

	if len(e.App.Image) != 0 {
		s.Update("%s: building app container... (skipped, using %s)", e.App.Name, e.App.Image)

		s.Done()
		return nil
	}

	image := fmt.Sprintf("%s-%s", e.Project.Namespace, e.App.Name)
	imageUri := fmt.Sprintf("%s/%s", e.App.DockerRegistry, image)

	relProjectPath, err := filepath.Rel(e.Project.RootDir, e.App.Path)
	if err != nil {
		return fmt.Errorf("unable to get relative path: %w", err)
	}

	cache := []string{fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", e.Project.Env))}
	logrus.Debugf("Using CACHE_IMAGE: %s", cache)

	buildArgs := map[string]*string{
		"PROJECT_PATH": &relProjectPath,
		"APP_PATH":     &relProjectPath,
		"APP_NAME":     &e.App.Name,
		"CACHE_IMAGE":  &cache[0],
	}

	tags := []string{
		image,
		fmt.Sprintf("%s:%s", imageUri, e.Project.Tag),
		fmt.Sprintf("%s:%s", imageUri, fmt.Sprintf("%s-latest", e.Project.Env)),
	}

	dockerfile := path.Join(e.App.Path, "Dockerfile")

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

	err = b.Build(ui, s, e.Project.RootDir)
	if err != nil {
		return fmt.Errorf("unable to build image: %w", err)
	}

	s.Done()

	return nil
}

//func definitionsToBulletItems(definitions *Helm.ListTaskDefinitionsOutput) []pterm.BulletListItem {
//	var items []pterm.BulletListItem
//	//for _, arn := range definitions.TaskDefinitionArns {
//	//	items = append(items, pterm.BulletListItem{Level: 0, Text: *arn})
//	//}
//
//	return items
//}

func (e *Manager) Destroy(ui terminal.UI, autoApprove bool) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: destroying Helm application...", e.App.Name)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 200) }()

	//name := fmt.Sprintf("%s-%s", e.Project.Env, e.App.Name)
	//
	//svc := e.Project.AWSClient.helmClient
	//
	//definitions, err := svc.ListTaskDefinitions(&helm.ListTaskDefinitionsInput{
	//	FamilyPrefix: &name,
	//	Sort:         aws.String(helm.SortOrderDesc),
	//})
	//if err != nil {
	//	return fmt.Errorf("can't get list task definitions of '%s': %v", name, err)
	//}
	//
	//if !autoApprove {
	//	pterm.SetDefaultOutput(s.TermOutput())
	//
	//	pterm.Printfln("this will destroy the following:")
	//	pterm.DefaultBulletList.WithItems(definitionsToBulletItems(definitions)).Render()
	//
	//	isContinue, err := pterm.DefaultInteractiveConfirm.WithDefaultText("Continue?").Show()
	//	if err != nil {
	//		return err
	//	}
	//
	//	if !isContinue {
	//		return fmt.Errorf("destroying was canceled")
	//	}
	//}
	//
	//for _, tda := range definitions.TaskDefinitionArns {
	//	_, err := e.Project.AWSClient.helmClient.DeregisterTaskDefinition(&helm.DeregisterTaskDefinitionInput{
	//		TaskDefinition: tda,
	//	})
	//	if err != nil {
	//		return fmt.Errorf("can't deregister task definition '%s': %v", *tda, err)
	//	}
	//}

	s.Done()
	s = sg.Add("%s: destroying completed!", e.App.Name)
	s.Done()

	return nil
}
func (e *Manager) Redeploy(ui terminal.UI) error {
	return nil
}
func (e *Manager) Explain() error {
	return nil
}
