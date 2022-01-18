package ecsdeploy

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Option struct {
	Tags       []string
	BuildArgs  map[string]*string
	Dockerfile string
	CacheFrom  []string
	ContextDir string
}

type Service struct {
	Type              string
	Path              string
	Image             string
	EcsCluster        string
	TaskDefinitionArn string
}

const ansi = `\x1B(?:[@-Z\\-_]|\[[0-?]*[-\]*[@-~])`

func DeployService(s Service, sname string, tag string, cfg *config.Config, sess *session.Session) error {
	var err error

	skipBuildAndPush := true

	if len(s.Image) == 0 {
		skipBuildAndPush = false
	}

	if !skipBuildAndPush {
		dockerImageName := fmt.Sprintf("%s-%s", cfg.Namespace, sname)
		dockerRegistry := viper.GetString("DOCKER_REGISTRY")
		tag := tag
		tagLatest := fmt.Sprintf("%s-latest", cfg.Env)
		contextDir := s.Path

		if !filepath.IsAbs(contextDir) {
			if contextDir, err = filepath.Abs(contextDir); err != nil {
				return fmt.Errorf("cat't deploy service %s: %w", sname, err)
			}
		}

		projectPath, err := filepath.Rel(viper.GetString("ROOT_DIR"), contextDir)
		if err != nil {
			return fmt.Errorf("cat't deploy service %s: %w", sname, err)
		}

		dockerfile := contextDir + "/Dockerfile"

		if _, err := os.Stat(dockerfile); err != nil {
			return fmt.Errorf("cat't deploy service %s: %w", sname, err)
		}

		s.Image = fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, strings.Trim(tag, "\n"))

		err = buildImage(Option{
			Tags: []string{
				dockerImageName,
				s.Image,
				fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest),
			},
			Dockerfile: dockerfile,
			BuildArgs: map[string]*string{
				"DOCKER_REGISTRY":   &dockerRegistry,
				"DOCKER_IMAGE_NAME": &dockerImageName,
				"ENV":               &cfg.Env,
				"PROJECT_PATH":      &projectPath,
			},
			CacheFrom: []string{
				fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest),
			},
			ContextDir: viper.GetString("ROOT_DIR"),
		})
		if err != nil {
			return fmt.Errorf("cat't deploy service %s: %w", sname, err)
		}

		svc := ecr.New(sess)

		repOut, err := svc.DescribeRepositories(&ecr.DescribeRepositoriesInput{
			RepositoryNames: []*string{aws.String(dockerImageName)},
		})
		if err != nil {
			_, ok := err.(*ecr.RepositoryNotFoundException)
			if !ok {
				return fmt.Errorf("cat't deploy service %s: %w", sname, err)
			}
		}

		if repOut == nil || len(repOut.Repositories) == 0 {
			logrus.Info("no ECR repository detected, creating", "name", dockerImageName)

			_, err := svc.CreateRepository(&ecr.CreateRepositoryInput{
				RepositoryName: aws.String(dockerImageName),
			})
			if err != nil {
				return fmt.Errorf("cat't deploy service %s: %w", sname, err)
			}
		}

		gat, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
		if err != nil {
			return fmt.Errorf("cat't deploy service %s: %w", sname, err)
		}

		if len(gat.AuthorizationData) == 0 {
			return fmt.Errorf("cat't deploy service %s: not found authorization data", sname)
		}

		token := *gat.AuthorizationData[0].AuthorizationToken

		err = push(
			[]string{
				fmt.Sprintf("%s/%s:%s", dockerRegistry, dockerImageName, tagLatest),
			},
			token,
			dockerRegistry,
		)
		if err != nil {
			return fmt.Errorf("cat't deploy service %s: %w", sname, err)
		}
	} else {
		tag = strings.Split(s.Image, ":")[1]
	}

	if s.TaskDefinitionArn == "" {
		stdo, err := ecs.New(sess).DescribeTaskDefinition(&ecs.DescribeTaskDefinitionInput{
			TaskDefinition: aws.String(fmt.Sprintf("%s-%s", cfg.Env, sname)),
		})
		if err != nil {
			return fmt.Errorf("cat't deploy service %s: %w", sname, err)
		}

		s.TaskDefinitionArn = *stdo.TaskDefinition.TaskDefinitionArn
	}

	err = deploy(DeployOpts{
		Service:           sname,
		Cluster:           s.EcsCluster,
		TaskDefinitionArn: s.TaskDefinitionArn,
		AwsProfile:        cfg.AwsProfile,
		Tag:               tag,
		Timeout:           "600",
		Image:             s.Image,
		EcsService:        fmt.Sprintf("%s-%s", cfg.Env, sname),
	})
	if err != nil {
		return fmt.Errorf("cat't deploy service %s: %w", sname, err)
	}

	return nil
}

func buildImage(opts Option) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Error("docker client initialization")
		return err
	}

	logrus.Debugf("build options: %s", opts)

	contextDir := opts.ContextDir
	dockerfile := opts.Dockerfile

	dockerfile, err = filepath.Rel(contextDir, dockerfile)
	if err != nil {
		return err
	}

	excludes, err := build.ReadDockerignore(contextDir)
	if err != nil {
		return status.Errorf(codes.Internal, "unable to read .dockerignore: %s", err)
	}

	if err := build.ValidateContextDirectory(contextDir, excludes); err != nil {
		return status.Errorf(codes.Internal, "error checking context: %s", err)
	}

	excludes = build.TrimBuildFilesFromExcludes(excludes, opts.Dockerfile, false)

	buildCtx, err := archive.TarWithOptions(contextDir, &archive.TarOptions{
		ExcludePatterns: excludes,
		ChownOpts:       &idtools.Identity{UID: 0, GID: 0},
	})
	if err != nil {
		return status.Errorf(codes.Internal, "unable to compress context: %s", err)
	}

	logrus.Debug("CacheFrom:", opts.CacheFrom)
	logrus.Debug("Tags:", opts.Tags)
	logrus.Debug("BuildArgs:", opts.BuildArgs)
	logrus.Debug("Dockerfile:", dockerfile)
	logrus.Debug("buildCtx:", buildCtx)

	resp, err := cli.ImageBuild(context.Background(), buildCtx, types.ImageBuildOptions{
		CacheFrom:  opts.CacheFrom,
		Tags:       opts.Tags,
		BuildArgs:  opts.BuildArgs,
		Dockerfile: dockerfile,
	})
	if err != nil {
		logrus.Error(err)
		return err
	}

	wr := ioutil.Discard
	if logrus.GetLevel() >= 4 {
		wr = os.Stdout
	}

	var termFd uintptr

	err = jsonmessage.DisplayJSONMessagesStream(
		resp.Body,
		wr,
		termFd,
		true,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func push(images []string, ecrToken string, registry string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Error("docker client initialization")
		return err
	}

	authBase64, err := getAuthToken(ecrToken, registry)
	if err != nil {
		return err
	}

	for _, i := range images {
		resp, err := cli.ImagePush(context.Background(), i, types.ImagePushOptions{
			All:          true,
			RegistryAuth: authBase64,
		})
		if err != nil {
			return err
		}

		wr := ioutil.Discard
		if logrus.GetLevel() >= 4 {
			wr = os.Stdout
		}

		var termFd uintptr

		err = jsonmessage.DisplayJSONMessagesStream(
			resp,
			wr,
			termFd,
			true,
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

type DeployOpts struct {
	AwsProfile        string
	Cluster           string
	Service           string
	TaskDefinitionArn string
	Tag               string
	Timeout           string
	Image             string
	EcsService        string
}

func deploy(opts DeployOpts) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Error("docker client initialization")
		return err
	}

	imageName := "hazelops/ecs-deploy"
	imageTag := "latest"

	logrus.Infof("image name: %s, image tag: %s", imageName, imageTag)

	cmd := []string{
		"ecs",
		"deploy",
		"--profile", opts.AwsProfile,
		opts.Cluster,
		opts.EcsService,
		"--task", opts.TaskDefinitionArn,
		"--image", opts.Service,
		opts.Image,
		"--diff",
		"--timeout", opts.Timeout,
		"--rollback",
		"-e", opts.Service,
		"DD_VERSION", opts.Tag,
	}

	out, err := cli.ImagePull(context.Background(), "hazelops/ecs-deploy", types.ImagePullOptions{})
	if err != nil {
		logrus.Errorf("pulling terraform image %v:%v", imageName, imageTag)
		return err
	}

	wr := ioutil.Discard
	if logrus.GetLevel() >= 4 {
		wr = os.Stdout
	}

	var termFd uintptr

	err = jsonmessage.DisplayJSONMessagesStream(
		out,
		wr,
		termFd,
		true,
		nil,
	)
	if err != nil {
		logrus.Errorf("pulling ecs-deploy image %v:%v", imageName, imageTag)
		return err
	}

	logrus.Debugf("pulling ecs-deploy image %v:%v", imageName, imageTag)

	contConfig := &container.Config{
		User:         fmt.Sprintf("%v:%v", os.Getuid(), os.Getgid()),
		Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
		Tty:          true,
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		WorkingDir:   fmt.Sprintf("%v", viper.Get("ENV_DIR")),
	}

	contHostConfig := &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v/.aws", viper.Get("HOME")),
				Target: "/.aws",
			},
		},
	}

	cont, err := cli.ContainerCreate(
		context.Background(),
		contConfig,
		contHostConfig,
		nil,
		nil,
		"ecs-deploy",
	)
	if err != nil {
		logrus.Errorf("creating terraform container from image %v:%v", imageName, imageTag)
		return err
	}

	logrus.Debugf("creating terraform container from image %v:%v", imageName, imageTag)

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		logrus.Errorf("terraform container started: %s", cont.ID)
		return err
	}

	reader, err := cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}

	logrus.Debugf("terraform container started: %s", cont.ID)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "error") {
			r := regexp.MustCompile(ansi)
			err = fmt.Errorf(r.ReplaceAllString(scanner.Text(), ""))
		}
		if logrus.GetLevel() >= 4 {
			fmt.Println(scanner.Text())
		}
	}

	if err != nil {
		return err
	}

	wait, errC := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionRemoved)

	select {
	case status := <-wait:
		logrus.Debugf("container exit status code %d", status.StatusCode)
		if status.StatusCode == 1 {
			return err
		}
		return nil
	case err := <-errC:
		return err
	}
}

func getAuthToken(ecrToken, registry string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ecrToken)
	if err != nil {
		return "", err
	}

	auth := types.AuthConfig{
		Username:      "AWS",
		Password:      string(data[4:]),
		ServerAddress: registry,
	}

	authBytes, _ := json.Marshal(auth)

	return base64.URLEncoding.EncodeToString(authBytes), nil
}
