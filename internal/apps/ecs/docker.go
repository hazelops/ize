package ecs

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	dockerutils "github.com/hazelops/ize/internal/docker/utils"

	"github.com/spf13/viper"
)

func (e *ecs) deployWithDocker(w io.Writer) error {
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
