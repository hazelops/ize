package apps

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	ecssvc "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/docker"
	dockerutils "github.com/hazelops/ize/internal/docker/utils"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/mitchellh/mapstructure"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ecsDeployImage = "hazelops/ecs-deploy:latest"

type ecs struct {
	Name              string
	Path              string
	Image             string
	Cluster           string
	TaskDefinitionArn string
	Timeout           int
	AwsProfile        string
	AwsRegion         string
	Tag               string
}

func NewECSDeployment(name string, app interface{}) *ecs {
	ecsConfig := ecs{}

	raw, ok := app.(map[string]interface{})
	if ok {
		ecsConfig.Name = name
		mapstructure.Decode(raw, &ecsConfig)
	}

	ecsConfig.Name = name

	if ecsConfig.Path == "" {
		ecsConfig.Path = fmt.Sprintf("./projects/%s", name)
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
func (e *ecs) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s: initializing Docker client...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	skipBuildAndPush := true
	env := viper.GetString("env")
	tagLatest := fmt.Sprintf("%s-latest", env)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	if len(e.Image) == 0 {
		skipBuildAndPush = false
	}

	if !skipBuildAndPush {
		s.Update("%s: building app container...", e.Name)
		dockerImageName := fmt.Sprintf("%s-%s", viper.GetString("namespace"), e.Name)
		dockerRegistry := viper.GetString("DOCKER_REGISTRY")

		e.Image = fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, strings.Trim(e.Tag, "\n"))
		b := docker.NewBuilder(
			map[string]*string{
				"DOCKER_REGISTRY":   &dockerRegistry,
				"DOCKER_IMAGE_NAME": &dockerImageName,
				"ENV":               &env,
				"PROJECT_PATH":      &e.Path,
			},
			[]string{
				dockerImageName,
				e.Image,
				fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest),
			},
			path.Join(e.Path, "Dockerfile"),
			[]string{
				fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest),
			},
		)

		err = b.Build(ui, s)
		if err != nil {
			return fmt.Errorf("can't deploy app %s: %w", e.Name, err)
		}

		s.Done()

		s = sg.Add("%s: push app container...", e.Name)

		repo, token, err := setupRepo(dockerImageName, e.AwsRegion, e.AwsProfile)
		if err != nil {
			return err
		}

		err = cli.ImageTag(context.Background(), e.Image, fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest))
		if err != nil {
			return err
		}

		err = push(
			cli,
			ui,
			s,
			fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest),
			token,
			repo,
		)
		if err != nil {
			return fmt.Errorf("can't deploy app %s: %w", e.Name, err)
		}
	} else {
		e.Tag = strings.Split(e.Image, ":")[1]
	}

	imageRef, err := reference.ParseNormalizedNamed(ecsDeployImage)
	if err != nil {
		return fmt.Errorf("error parsing Docker image: %s", err)
	}

	imageList, err := cli.ImageList(context.Background(), types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "reference",
			Value: reference.FamiliarString(imageRef),
		}),
	})
	if err != nil {
		return err
	}

	if len(imageList) == 0 {
		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		if err != nil {
			return err
		}

		err = jsonmessage.DisplayJSONMessagesStream(resp, os.Stderr, os.Stderr.Fd(), true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}
	}

	s.Done()

	if viper.GetString("prefer-runtime") == "native" {
		err := e.deployLocal(sg)
		if err != nil {
			return err
		}
	} else {
		err := e.deployWithDocker(cli, sg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *ecs) Destroy(sg terminal.StepGroup, ui terminal.UI) error {
	ui.Output("Destroying ECS applications requires destroying the infrastructure.", terminal.WithWarningStyle())
	time.Sleep(time.Millisecond * 100)

	s := sg.Add("%s: destroying completed!", e.Name)
	defer func() { s.Abort() }()
	s.Done()

	return nil
}

// TODO: refactor
func push(cli *client.Client, ui terminal.UI, s terminal.Step, image string, ecrToken string, registry string) error {
	authBase64, err := getAuthToken(ecrToken)
	if err != nil {
		return err
	}

	resp, err := cli.ImagePush(context.Background(), image, types.ImagePushOptions{
		All:          true,
		RegistryAuth: authBase64,
	})
	if err != nil {
		return fmt.Errorf("can't push image %s: %w", image, err)
	}

	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return err
	}

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	err = jsonmessage.DisplayJSONMessagesStream(
		resp,
		s.TermOutput(),
		termFd,
		true,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func getAuthToken(ecrToken string) (string, error) {
	auth := types.AuthConfig{
		Username: "AWS",
		Password: ecrToken,
	}

	authBytes, _ := json.Marshal(auth)

	return base64.URLEncoding.EncodeToString(authBytes), nil
}

func setupRepo(repoName string, region string, profile string) (string, string, error) {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  region,
		Profile: profile,
	})
	if err != nil {
		return "", "", err
	}
	svc := ecr.New(sess)

	gat, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", "", err
	}
	if len(gat.AuthorizationData) == 0 {
		return "", "", fmt.Errorf("no authorization tokens provided")
	}

	repOut, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(repoName)},
	})
	if err != nil {
		_, ok := err.(*ecr.RepositoryNotFoundException)
		if !ok {
			return "", "", err
		}
	}

	var repo *ecr.Repository
	if repOut == nil || len(repOut.Repositories) == 0 {
		logrus.Info("no ECR repository detected, creating", "name", repoName)

		out, err := svc.CreateRepository(&ecr.CreateRepositoryInput{
			RepositoryName: aws.String(repoName),
		})
		if err != nil {
			return "", "", fmt.Errorf("unable to create repository: %w", err)
		}

		repo = out.Repository
	} else {
		repo = repOut.Repositories[0]
	}

	uptoken := *gat.AuthorizationData[0].AuthorizationToken
	data, err := base64.StdEncoding.DecodeString(uptoken)
	if err != nil {
		return "", "", err
	}

	return *repo.RepositoryUri, string(data[4:]), nil
}

func (e *ecs) deployWithDocker(cli *client.Client, sg terminal.StepGroup) error {
	s := sg.Add("%s: deploying app container...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	cmd := []string{"ecs", "deploy",
		"--profile", e.AwsProfile,
		e.Cluster,
		fmt.Sprintf("%s-%s", viper.GetString("env"), e.Name),
		"--task", e.TaskDefinitionArn,
		"--image", e.Name,
		e.Image,
		"--diff",
		"--timeout", strconv.Itoa(e.Timeout),
		"--rollback",
		"-e", e.Name,
		"DD_VERSION", e.Tag,
	}

	cfg := container.Config{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          true,
		Image:        ecsDeployImage,
		User:         fmt.Sprintf("%v:%v", os.Getuid(), os.Getgid()),
		WorkingDir:   fmt.Sprintf("%v", viper.Get("ENV_DIR")),
		Cmd:          cmd,
	}

	hostconfig := container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v/.aws", viper.Get("HOME")),
				Target: "/.aws",
			},
		},
	}

	cr, err := cli.ContainerCreate(context.Background(), &cfg, &hostconfig, &network.NetworkingConfig{}, nil, e.Name)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), cr.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	dockerutils.SetupSignalHandlers(cli, cr.ID)

	out, err := cli.ContainerLogs(context.Background(), cr.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}

	defer out.Close()

	io.Copy(s.TermOutput(), out)

	wait, errC := cli.ContainerWait(context.Background(), cr.ID, container.WaitConditionRemoved)

	select {
	case status := <-wait:
		if status.StatusCode == 0 {
			s.Done()
			s = sg.Add("%s: deployment completed!", e.Name)
			s.Done()
			return nil
		}
		s.Status(terminal.ErrorStyle)
		return fmt.Errorf("container exit status code %d", status.StatusCode)
	case err := <-errC:
		return err
	}

	return nil
}

func (e *ecs) deployLocal(sg terminal.StepGroup) error {
	s := sg.Add("%s: deploying app container...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()
	pterm.SetDefaultOutput(s.TermOutput())

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  e.AwsRegion,
		Profile: e.AwsProfile,
	})
	if err != nil {
		return err
	}

	svc := ecssvc.New(sess)

	name := fmt.Sprintf("%s-%s", viper.GetString("env"), e.Name)

	dso, err := svc.DescribeServices(&ecssvc.DescribeServicesInput{
		Cluster:  &e.Cluster,
		Services: []*string{&name},
	})
	if err != nil {
		return err
	}

	if len(dso.Services) == 0 {
		return fmt.Errorf("app %s not found", name)
	}

	dtdo, err := svc.DescribeTaskDefinition(&ecssvc.DescribeTaskDefinitionInput{
		TaskDefinition: dso.Services[0].TaskDefinition,
	})
	if err != nil {
		return err
	}

	pterm.Printfln("Deploying based on task definition: %s:%d", *dtdo.TaskDefinition.Family, *dtdo.TaskDefinition.Revision)

	oldTaskDef := *dtdo.TaskDefinition

	var image string
	for i := 0; i < len(dtdo.TaskDefinition.ContainerDefinitions); i++ {
		container := dtdo.TaskDefinition.ContainerDefinitions[i]
		if container.Name == &e.Name {
			image = e.Image
			container.Image = &image
			pterm.Printfln(`Changed image of container "%s" to : "%s" (was: "%s")`, *container.Name, image, *container.Image)
		} else if len(e.Tag) != 0 {
			name := strings.Split(*container.Image, ":")[0]
			image = fmt.Sprintf("%s:%s", name, e.Tag)
			container.Image = &image
			pterm.Printfln(`Changed image of container "%s" to : "%s" (was: "%s")`, *container.Name, image, *container.Image)

		}
	}

	pterm.Println("Creating new task definition revision")

	rtdo, err := svc.RegisterTaskDefinition(&ecssvc.RegisterTaskDefinitionInput{
		ContainerDefinitions:    dtdo.TaskDefinition.ContainerDefinitions,
		Family:                  dtdo.TaskDefinition.Family,
		Volumes:                 dtdo.TaskDefinition.Volumes,
		TaskRoleArn:             dtdo.TaskDefinition.TaskRoleArn,
		ExecutionRoleArn:        dtdo.TaskDefinition.ExecutionRoleArn,
		RuntimePlatform:         dtdo.TaskDefinition.RuntimePlatform,
		RequiresCompatibilities: dtdo.TaskDefinition.RequiresCompatibilities,
		NetworkMode:             dtdo.TaskDefinition.NetworkMode,
		Cpu:                     dtdo.TaskDefinition.Cpu,
		Memory:                  dtdo.TaskDefinition.Memory,
	})
	if err != nil {
		return err
	}

	pterm.Printfln("Successfully created revision: %s:%d", *rtdo.TaskDefinition.Family, *rtdo.TaskDefinition.Revision)

	if err = e.updateTaskDefinition(svc, rtdo.TaskDefinition, name, "Deploying new task definition"); err != nil {
		pterm.Printfln("Rolling back to old task definition: %s:%d", *oldTaskDef.Family, *oldTaskDef.Revision)
		e.Timeout = 600
		if err = e.updateTaskDefinition(svc, &oldTaskDef, name, "Deploying previous task definition"); err != nil {
			s.TermOutput().Write([]byte(err.Error() + " - try rollback\n"))
			return err
		}

		pterm.Println("Rollback successful")

		if err = deregisterTaskDefinition(svc, &oldTaskDef); err != nil {
			return err
		}

		pterm.Println("Deployment failed, but service has been rolled back to previous task definition:", *oldTaskDef.Family)
	}

	if err = deregisterTaskDefinition(svc, &oldTaskDef); err != nil {
		return err
	}

	s.Done()

	return nil
}

func (e *ecs) updateTaskDefinition(svc *ecssvc.ECS, td *ecssvc.TaskDefinition, serviceName string, title string) error {
	pterm.Println("Updating service")

	_, err := svc.UpdateService(&ecssvc.UpdateServiceInput{
		Service:        aws.String(serviceName),
		Cluster:        aws.String(e.Cluster),
		TaskDefinition: aws.String(*td.TaskDefinitionArn),
	})
	if err != nil {
		return fmt.Errorf("unable to update service: %w", err)
	}

	pterm.Printfln("Successfully changed task definition to: %s:%d", *td.Family, *td.Revision)
	pterm.Println(title)

	waitingTimeout := time.Now().Add(time.Duration(e.Timeout) * time.Second)
	waiting := true

	for waiting && time.Now().Before(waitingTimeout) {
		d, err := e.isDeployed(svc, serviceName)
		if err != nil {
			return err
		}

		waiting = !d

		if waiting {
			time.Sleep(time.Second * 5)
		}
	}

	if waiting && time.Now().After(waitingTimeout) {
		pterm.Println("Deployment failed due to timeout")
		return fmt.Errorf("deployment failed due to timeout")
	}

	return nil
}

func (e *ecs) isDeployed(svc *ecssvc.ECS, name string) (bool, error) {
	dso, err := svc.DescribeServices(&ecssvc.DescribeServicesInput{
		Cluster:  &e.Cluster,
		Services: []*string{&name},
	})
	if err != nil {
		return false, err
	}

	if len(dso.Services) == 0 {
		return false, nil
	}

	if len(dso.Services[0].Deployments) != 1 {
		return false, nil
	}

	runningTasks, err := svc.ListTasks(&ecssvc.ListTasksInput{
		Cluster:     &e.Cluster,
		ServiceName: &name,
	})
	if err != nil {
		return false, err
	}

	if runningTasks.TaskArns == nil {
		return *dso.Services[0].DesiredCount == 0, nil
	}

	runningCount, err := getRunningTaskCount(e.Cluster, runningTasks.TaskArns, *dso.Services[0].TaskDefinition, svc)
	if err != nil {
		return false, err
	}

	return runningCount == *dso.Services[0].DesiredCount, nil
}

func getRunningTaskCount(cluster string, tasks []*string, serviceArn string, svc *ecssvc.ECS) (int64, error) {
	count := 0

	dto, err := svc.DescribeTasks(&ecssvc.DescribeTasksInput{
		Cluster: &cluster,
		Tasks:   tasks,
	})
	if err != nil {
		return 0, err
	}

	for _, t := range dto.Tasks {
		if *t.TaskDefinitionArn == serviceArn && *t.LastStatus == "RUNNING" {
			count++
		}
	}

	return int64(count), nil
}

func deregisterTaskDefinition(svc *ecssvc.ECS, td *ecssvc.TaskDefinition) error {
	pterm.Println("Deregister task definition revision")

	_, err := svc.DeregisterTaskDefinition(&ecssvc.DeregisterTaskDefinitionInput{
		TaskDefinition: td.TaskDefinitionArn,
	})
	if err != nil {
		return err
	}

	pterm.Printfln("Successfully deregistered revision: %s:%d", *td.Family, *td.Revision)

	return nil
}
