package terraform

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const ansi = `\x1B(?:[@-Z\\-_]|\[[0-?]*[-\]*[@-~])`

type Options struct {
	Env              []string
	TerraformVersion string
	ContainerName    string
	Cmd              []string
	OutputPath       string
}

func cleanupOldContainers(cli *client.Client, opts Options) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], opts.ContainerName) {
			err = cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Run(opts Options) error {
	logrus.Debugf("terraform run options: %s", opts)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.Error("docker client initialization")
		return err
	}

	logrus.Debug("docker client initialization")

	err = cleanupOldContainers(cli, opts)
	if err != nil {
		return err
	}

	logrus.Debug("cleanup old containers successfully")

	imageName := "hashicorp/terraform"
	imageTag := opts.TerraformVersion

	logrus.Infof("image name: %s, image tag: %s", imageName, imageTag)

	out, err := cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		logrus.Errorf("pulling terraform image %v:%v", imageName, imageTag)
		return err
	}

	wr := ioutil.Discard
	if logrus.GetLevel() >= 4 {
		wr = os.Stderr
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
		logrus.Errorf("pulling terraform image %v:%v", imageName, imageTag)
		return err
	}

	logrus.Debugf("pulling terraform image %v:%v", imageName, imageTag)

	contConfig := &container.Config{
		User:         fmt.Sprintf("%v:%v", os.Getuid(), os.Getgid()),
		Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
		Tty:          true,
		Cmd:          opts.Cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		WorkingDir:   fmt.Sprintf("%v", viper.Get("ENV_DIR")),
		Env:          opts.Env,
	}

	contHostConfig := &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v", viper.Get("ENV_DIR")),
				Target: fmt.Sprintf("%v", viper.Get("ENV_DIR")),
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v", viper.Get("INFRA_DIR")),
				Target: fmt.Sprintf("%v", viper.Get("INFRA_DIR")),
			},
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
		opts.ContainerName,
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

	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	var f *os.File

	if opts.OutputPath != "" {
		f, err = os.Create(opts.OutputPath)
		if err != nil {
			return err
		}

		defer f.Close()
	}

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Error: ") {
			r := regexp.MustCompile(ansi)
			strErr := r.ReplaceAllString(scanner.Text(), "")
			strErr = strErr[strings.LastIndex(strErr, "Error: "):]
			strErr = strings.TrimPrefix(strErr, "Error: ")
			strErr = strings.ToLower(string(strErr[0])) + strErr[1:]
			err = fmt.Errorf(strErr)
		}
		if logrus.GetLevel() >= 4 || strings.Contains(strings.Join(opts.Cmd, ""), "output") {
			fmt.Println(scanner.Text())
		}
		if f != nil {
			f.WriteString(scanner.Text())
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
