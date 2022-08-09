package terraform

import (
	"context"
	"fmt"
	"github.com/hazelops/ize/internal/config"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/internal/docker/utils"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/sirupsen/logrus"
	t "golang.org/x/crypto/ssh/terminal"
)

const (
	defaultName = "ize-terraform"
)

func cleanupOldContainers(cli *client.Client) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	for _, c := range containers {
		if strings.Contains(c.Names[0], defaultName) {
			err = cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type docker struct {
	version string
	command []string
	env     []string
	output  io.Writer
	project *config.Project
}

func NewDockerTerraform(version string, command []string, env []string, out io.Writer, project *config.Project) *docker {
	return &docker{
		version: version,
		command: command,
		env:     env,
		output:  out,
		project: project,
	}
}

func (d *docker) Prepare() error {
	return nil
}

func (d *docker) NewCmd(cmd []string) {
	d.command = cmd
}

func (d *docker) SetOut(out io.Writer) {
	d.output = out
}

func (d *docker) RunUI(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Initializing Docker client...")
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 50) }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	s.Done()
	s = sg.Add("Cleanuping old containers...")

	err = cleanupOldContainers(cli)
	if err != nil {
		return err
	}

	imageName := "hashicorp/terraform"
	imageTag := d.version

	imageRef, err := reference.ParseNormalizedNamed(fmt.Sprintf("%s:%s", imageName, imageTag))
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

	stdout, _, err := ui.OutputWriters()
	if err != nil {
		return err
	}

	var termFd uintptr
	if f, ok := stdout.(*os.File); ok {
		termFd = f.Fd()
	}

	if len(imageList) == 0 {
		s.Update("pulling terraform image %v:%v...", imageName, imageTag)
		out, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer out.Close()

		if err != nil {
			return err
		}

		err = jsonmessage.DisplayJSONMessagesStream(out, s.TermOutput(), termFd, true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}

		s.Done()
		s = sg.Add("")
	}

	logrus.Infof("image name: %s, image tag: %s", imageName, imageTag)

	contConfig := &container.Config{
		User:         fmt.Sprintf("%v:%v", os.Getuid(), os.Getgid()),
		Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
		Tty:          true,
		Cmd:          d.command,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		WorkingDir:   fmt.Sprintf("%v", d.project.EnvDir),
		Env:          d.env,
	}

	contHostConfig := &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v", d.project.EnvDir),
				Target: fmt.Sprintf("%v", d.project.EnvDir),
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v", d.project.InfraDir),
				Target: fmt.Sprintf("%v", d.project.InfraDir),
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v/.aws", d.project.Home),
				Target: "/.aws",
			},
		},
	}

	s.Update("[%s] running terraform image %v:%v...", d.project.Env, imageName, imageTag)

	cont, err := cli.ContainerCreate(
		context.Background(),
		contConfig,
		contHostConfig,
		nil,
		nil,
		defaultName,
	)

	if err != nil {
		return err
	}

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	utils.SetupSignalHandlers(cli, cont.ID)

	reader, err := cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}

	defer reader.Close()

	if d.output != nil {
		io.Copy(d.output, reader)

	} else {
		io.Copy(s.TermOutput(), reader)
	}

	wait, errC := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionRemoved)

	select {
	case status := <-wait:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exit status code %d\n", status.StatusCode)
		}
		s.Done()
		return nil
	case err := <-errC:
		return err
	}
}

func (d *docker) Run() error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	err = cleanupOldContainers(cli)
	if err != nil {
		return err
	}

	imageName := "hashicorp/terraform"
	imageTag := d.version

	imageRef, err := reference.ParseNormalizedNamed(fmt.Sprintf("%s:%s", imageName, imageTag))
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
		out, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer out.Close()

		if err != nil {
			return err
		}

		err = jsonmessage.DisplayJSONMessagesStream(out, os.Stderr, os.Stderr.Fd(), true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}
	}

	logrus.Infof("image name: %s, image tag: %s", imageName, imageTag)

	contConfig := &container.Config{
		User:         fmt.Sprintf("%v:%v", os.Getuid(), os.Getgid()),
		Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
		Tty:          true,
		Cmd:          d.command,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		WorkingDir:   fmt.Sprintf("%v", d.project.EnvDir),
		Env:          d.env,
	}

	contHostConfig := &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v", d.project.EnvDir),
				Target: fmt.Sprintf("%v", d.project.EnvDir),
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v", d.project.InfraDir),
				Target: fmt.Sprintf("%v", d.project.InfraDir),
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%v/.aws", d.project.Home),
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
		defaultName,
	)

	if err != nil {
		return err
	}

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	waiter, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return err
	}

	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(os.Stderr, waiter.Reader)
	go io.Copy(waiter.Conn, os.Stdin)

	fd := int(os.Stdin.Fd())
	var oldState *t.State
	if t.IsTerminal(fd) {
		oldState, err = t.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer t.Restore(fd, oldState)
	}

	wait, errC := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionRemoved)

	select {
	case status := <-wait:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exit status code %d\n", status.StatusCode)
		}
		return nil
	case err := <-errC:
		return err
	}
}
