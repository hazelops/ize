package serverless

import (
	"context"
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"os"
	"path/filepath"
	"time"
)

func (sls *Manager) deployWithDocker(ui *echelon.Logger) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	image := "node:" + sls.App.NodeVersion

	s := ui.Scoped(fmt.Sprintf("%s: checking for Docker image: %s", sls.App.Name, image))

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
	s.Finish(true)

	if len(imageList) == 0 {
		s = ui.Scoped(fmt.Sprintf("%s: pulling image: %s", sls.App.Name, image))

		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		err = jsonmessage.DisplayJSONMessagesStream(resp, &Writer{logger: s}, os.Stdout.Fd(), true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}
		s.Finish(true)
		time.Sleep(time.Millisecond * 200)
	}

	s = ui.Scoped(fmt.Sprintf("%s: downloading npm modules...", sls.App.Name))
	err = sls.npm(cli, []string{"npm", "install", "--save-dev"}, s)
	time.Sleep(time.Millisecond * 50)
	if err != nil {
		return fmt.Errorf("can't deploy %s: %w", sls.App.Name, err)
	}
	s.Finish(true)

	if sls.App.CreateDomain {
		s = ui.Scoped(fmt.Sprintf("%s: creating domain...", sls.App.Name))
		err = sls.serverless(cli, []string{
			"create_domain",
			"--verbose",
			"--region", sls.App.AwsRegion,
			"--profile", sls.App.AwsProfile,
			"--stage", sls.Project.Env,
		}, s)
		if err != nil {
			return err
		}

		s.Finish(true)
	}

	s = ui.Scoped(fmt.Sprintf("%s: deploying app...", sls.App.Name))

	err = sls.serverless(cli, []string{
		"deploy",
		"--config", sls.App.File,
		"--service", sls.App.Name,
		"--verbose",
		"--region", sls.App.AwsRegion,
		"--profile", sls.App.AwsProfile,
		"--stage", sls.Project.Env,
	}, s)
	if err != nil {
		time.Sleep(time.Second)
		return err
	}

	s.Finish(true)

	return nil
}

func (sls *Manager) removeWithDocker(ui *echelon.Logger) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	image := "node:" + sls.App.NodeVersion

	s := ui.Scoped(fmt.Sprintf("%s: checking for Docker image: %s", sls.App.Name, image))

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
	s.Finish(true)

	if len(imageList) == 0 {
		s := ui.Scoped(fmt.Sprintf("%s: pulling image: %s", sls.App.Name, image))

		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		err = jsonmessage.DisplayJSONMessagesStream(resp, &Writer{logger: s}, os.Stdout.Fd(), true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}
		s.Finish(true)
	}

	s = ui.Scoped(fmt.Sprintf("%s: destroying app...", sls.App.Name))

	err = sls.serverless(cli, []string{
		"remove",
		"--config", sls.App.File,
		"--service", sls.App.Name,
		"--verbose",
		"--region", sls.App.AwsRegion,
		"--stage", sls.Project.Env,
		"--profile", sls.App.AwsProfile,
	}, s)
	if err != nil {
		return err
	}

	s.Finish(true)

	return nil
}

func (sls *Manager) serverless(cli *client.Client, cmd []string, ui *echelon.Logger) error {
	command := []string{"serverless"}
	command = append(command, cmd...)

	contConfig := &container.Config{
		Image:        fmt.Sprintf("node:%v", sls.App.NodeVersion),
		Entrypoint:   strslice.StrSlice{"/usr/local/bin/npx"},
		WorkingDir:   "/app",
		Env:          sls.App.Env,
		Tty:          false,
		Cmd:          command,
		AttachStdout: true,
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

	body, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream: true, Stdout: true, Stderr: true, Stdin: false})
	if err != nil {
		return err
	}
	defer body.Close()

	_, err = stdcopy.StdCopy(stdcopy.NewStdWriter(&Writer{logger: ui}, stdcopy.Stdout), stdcopy.NewStdWriter(&Writer{logger: ui}, stdcopy.Stderr), body.Reader)
	if err != nil {
		return err
	}
	if err != nil {
		return err
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

	return nil
}

func (sls *Manager) npm(cli *client.Client, cmd []string, ui *echelon.Logger) error {
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

	body, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{Stream: true, Stdout: true, Stderr: true})
	if err != nil {
		return err
	}
	defer body.Close()
	//_, err = stdcopy.StdCopy(&Writer{logger: ui}, &Writer{logger: ui}, body.Reader)
	go func() {
		_, err = io.Copy(&Writer{logger: ui}, body.Reader)
		if err != nil {
			ui.Errorf("error copying: %v", err)
		}
	}()

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

	return nil
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
				Source: sls.App.Path,
				Target: "/app",
			},
			{
				Type:   mount.TypeBind,
				Source: filepath.Join(sls.App.Path, ".serverless"),
				Target: "/root/.config",
			},
			{
				Type:   mount.TypeTmpfs,
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
