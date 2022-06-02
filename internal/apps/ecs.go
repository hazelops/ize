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

func NewECSApp(name string, app interface{}) *ecs {
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
func (e *ecs) Deploy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deploying app container...", e.Name)
	defer func() { s.Abort(); time.Sleep(50 * time.Millisecond) }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
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

	if e.Image == "" {
		e.Image = fmt.Sprintf("%s/%s:%s",
			viper.GetString("DOCKER_REGISTRY"),
			fmt.Sprintf("%s-%s", viper.GetString("namespace"), e.Name),
			fmt.Sprintf("%s-%s", viper.GetString("env"), "latest"))
	}

	if viper.GetString("prefer-runtime") == "native" {
		err := e.deployLocal(s.TermOutput())
		if err != nil {
			return err
		}
	} else {
		err := e.deployWithDocker(cli, s.TermOutput())
		if err != nil {
			return err
		}
	}

	s.Done()
	s = sg.Add("%s: deployment completed!", e.Name)
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

	buildArgs := map[string]*string{
		"PROJECT_PATH": &e.Path,
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

	err := b.Build(ui, s)
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

func (e *ecs) deployWithDocker(cli *client.Client, w io.Writer) error {
	cmd := []string{"ecs", "deploy",
		"--profile", e.AwsProfile,
		"--region", e.AwsRegion,
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

	io.Copy(w, out)

	wait, errC := cli.ContainerWait(context.Background(), cr.ID, container.WaitConditionRemoved)

	select {
	case status := <-wait:
		if status.StatusCode == 0 {
			return nil
		}
		return fmt.Errorf("container exit status code %d", status.StatusCode)
	case err := <-errC:
		return err
	}
}

func (e *ecs) deployLocal(w io.Writer) error {
	pterm.SetDefaultOutput(w)

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
			w.Write([]byte(err.Error() + " - try rollback\n"))
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
