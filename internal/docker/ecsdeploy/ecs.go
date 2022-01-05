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

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
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

const ansi = `\x1B(?:[@-Z\\-_]|\[[0-?]*[-\]*[@-~])`

func Build(opts Option) error {
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

func Push(images []string, ecrToken string, registry string) error {
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

func Deploy(opts DeployOpts) error {
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
		if strings.Contains(scanner.Text(), "Error: ") {
			r := regexp.MustCompile(ansi)
			strErr := r.ReplaceAllString(scanner.Text(), "")
			strErr = strErr[strings.LastIndex(strErr, "Error: "):]
			strErr = strings.TrimPrefix(strErr, "Error: ")
			strErr = strings.ToLower(string(strErr[0])) + strErr[1:]
			err = fmt.Errorf(strErr)
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
