package serverless

import (
	"context"
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/pkg/terminal"
	"os"
	"time"
)

func (sls *Manager) deployWithDocker(s terminal.Step) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	image := "node:" + sls.App.NodeVersion

	s.Update("%s: checking for Docker image: %s", sls.App.Name, image)

	imageRef, err := reference.ParseNormalizedNamed(image)
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
		s.Update("%s: pulling image: %s", sls.App.Name, image)

		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), os.Stdout.Fd(), true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}
	}

	s.Update("%s: downloading npm modules...", sls.App.Name)

	err = sls.npm(cli, []string{"npm", "install", "--save-dev"}, s)
	if err != nil {
		return fmt.Errorf("can't deploy %s: %w", sls.App.Name, err)
	}

	s.Done()

	if sls.App.CreateDomain {
		s.Update("%s: creating domain...", sls.App.Name)
		err = sls.serverless(cli, []string{
			"create_domain",
			"--verbose",
			"--region", sls.Project.AwsRegion,
			"--profile", sls.Project.AwsProfile,
			"--stage", sls.Project.Env,
		}, s)
		if err != nil {
			return err
		}

		s.Done()
	}

	s.Update("%s: deploying app...", sls.App.Name)

	err = sls.serverless(cli, []string{
		"deploy",
		"--config", sls.App.File,
		"--service", sls.App.Name,
		"--verbose",
		"--region", sls.Project.AwsRegion,
		"--profile", sls.Project.AwsProfile,
		"--env", sls.Project.Env,
	}, s)
	if err != nil {
		s.Abort()
		time.Sleep(time.Second)
		return err
	}

	return nil
}

func (sls *Manager) removeWithDocker(s terminal.Step) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	image := "node:" + sls.App.NodeVersion

	s.Update("%s: checking for Docker image: %s", sls.App.Name, image)

	imageRef, err := reference.ParseNormalizedNamed(image)
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
		s.Update("%s: pulling image: %s", sls.App.Name, image)

		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), os.Stdout.Fd(), true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}
	}

	s.Done()
	s.Update("%s: destroying app...", sls.App.Name)

	err = sls.serverless(cli, []string{
		"remove",
		"--config", sls.App.File,
		"--service", sls.App.Name,
		"--verbose",
		"--region", sls.Project.AwsRegion,
		"--env", sls.Project.Env,
		"--profile", sls.Project.AwsProfile,
	}, s)
	if err != nil {
		s.Abort()
		return err
	}

	return nil
}

func (sls *Manager) serverless(cli *client.Client, cmd []string, step terminal.Step) error {
	command := []string{"serverless"}
	command = append(command, cmd...)

	contConfig := &container.Config{
		Image:        fmt.Sprintf("node:%v", sls.App.NodeVersion),
		Entrypoint:   strslice.StrSlice{"/usr/local/bin/npx"},
		WorkingDir:   "/app",
		Env:          sls.App.Env,
		Tty:          true,
		Cmd:          command,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	contHostConfig := sls.getHostConfig(sls.Project.Home, sls.Project.RootDir)

	cont, err := cli.ContainerCreate(
		context.Background(),
		contConfig,
		contHostConfig,
		nil,
		nil,
		"ize_sls",
	)
	if err != nil {
		return fmt.Errorf("can't deploy app: %w", err)
	}

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("can't deploy app: %w", err)
	}

	body, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream: true, Stdout: true, Stderr: true, Stdin: true})
	if err != nil {
		return err
	}

	msgs := make(chan []byte)
	msgsErr := make(chan error)

	go func() {
		for {
			msg, er := body.Reader.ReadBytes('\n')
			if er != nil {
				msgsErr <- er
				return
			}
			msgs <- msg
		}
	}()

msgLoop:
	for {
		select {
		case msg := <-msgs:
			fmt.Fprintf(step.TermOutput(), "%s", msg)
		case <-msgsErr:
			break msgLoop
		}
	}

	wait, errC := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionRemoved)

	select {
	case status := <-wait:
		if status.StatusCode == 0 {
			return nil
		}
		return fmt.Errorf("container exit status code %d", status.StatusCode)
	case err := <-errC:
		return fmt.Errorf("can't deploy app: %w", err)
	}
}

func (sls *Manager) npm(cli *client.Client, cmd []string, s terminal.Step) error {
	contConfig := &container.Config{
		WorkingDir:   "/app",
		Image:        fmt.Sprintf("node:%v", sls.App.NodeVersion),
		Tty:          true,
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
	}

	contHostConfig := sls.getHostConfig(sls.Project.Home, sls.Project.RootDir)

	cont, err := cli.ContainerCreate(
		context.Background(),
		contConfig,
		contHostConfig,
		nil,
		nil,
		"sls",
	)
	if err != nil {
		return fmt.Errorf("can't deploy app: %w", err)
	}

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("can't deploy app: %w", err)
	}

	body, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream: true, Stdout: true})
	if err != nil {
		return err
	}

	msgs := make(chan []byte)
	msgsErr := make(chan error)

	go func() {
		for {
			msg, er := body.Reader.ReadBytes('\n')
			if er != nil {
				msgsErr <- er
				return
			}
			msgs <- msg
		}
	}()

msgLoop:
	for {
		select {
		case msg := <-msgs:
			fmt.Fprintf(s.TermOutput(), "%s", msg)
		case <-msgsErr:
			break msgLoop
		}
	}

	defer close(msgs)
	defer close(msgsErr)
	defer body.Close()

	wait, errC := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionNotRunning)

	select {
	case status := <-wait:
		if status.StatusCode == 0 {
			return nil
		}
		return fmt.Errorf("container exit status code %d", status.StatusCode)
	case err := <-errC:
		return fmt.Errorf("can't deploy app: %w", err)
	}
}

func (sls *Manager) getHostConfig(homeDir, rootDir string) *container.HostConfig {
	return &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				ReadOnly: true,
				Source:   fmt.Sprintf("%v/.aws", homeDir),
				Target:   "/root/.aws",
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%s/%s", rootDir, sls.App.Path),
				Target: "/app",
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%s/%s/.serverless/", rootDir, sls.App.Path),
				Target: "/root/.config",
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%s/.npm/", rootDir),
				Target: "/root/.npm",
			},
			{
				Type:   mount.TypeVolume,
				Source: sls.App.SLSNodeModuleCacheMount,
				Target: "/app/node_modules",
			},
		},
	}
}
