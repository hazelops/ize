package services

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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
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
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/mitchellh/mapstructure"
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
}

func NewECSDeployment(service App) *ecs {
	var ecsConfig ecs

	mapstructure.Decode(service, &ecsConfig)
	mapstructure.Decode(service.Body, &ecsConfig)
	ecsConfig.AwsProfile = viper.GetString("aws_profile")
	ecsConfig.AwsRegion = viper.GetString("aws_region")
	if len(ecsConfig.Cluster) == 0 {
		ecsConfig.Cluster = fmt.Sprintf("%s-%s", viper.GetString("env"), viper.GetString("namespace"))
	}

	return &ecsConfig
}

// Deploy deploys app container to ECS via ECS deploy
func (e *ecs) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s: initializing Docker client...", e.Name)
	defer func() { s.Abort() }()

	skipBuildAndPush := true
	tag := viper.GetString("tag")
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
		dockerImageName := fmt.Sprintf("%s-%s", env, e.Name)
		dockerRegistry := viper.GetString("DOCKER_REGISTRY")

		e.Image = fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, strings.Trim(tag, "\n"))
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
			return fmt.Errorf("can't deploy service %s: %w", e.Name, err)
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
			return fmt.Errorf("can't deploy service %s: %w", e.Name, err)
		}
	} else {
		tag = strings.Split(e.Image, ":")[1]
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
	s = sg.Add("%s: deploying app container...", e.Name)

	cmd := []string{"ecs", "deploy",
		"--profile", e.AwsProfile,
		e.Cluster,
		fmt.Sprintf("%s-%s", env, e.Name),
		"--task", e.TaskDefinitionArn,
		"--image", e.Name,
		e.Image,
		"--diff",
		"--timeout", strconv.Itoa(e.Timeout),
		"--rollback",
		"-e", e.Name,
		"DD_VERSION", tag,
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

	docker.SetupSignalHandlers(cli, cr.ID)

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
			s = sg.Add("%s deployment completed!", e.Name)
			s.Done()
			return nil
		}
		s.Status(terminal.ErrorStyle)
		return fmt.Errorf("container exit status code %d", status.StatusCode)
	case err := <-errC:
		return err
	}
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
