package apps

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type serverless struct {
	Name                    string
	File                    string
	NodeVersion             string `mapstructure:"node_version"`
	Env                     []string
	Path                    string
	SLSNodeModuleCacheMount string
	CreateDomain            bool `mapstructure:"creare_domain"`
}

func NewServerlessDeployment(name string, app interface{}) *serverless {
	slsConfig := serverless{}

	raw, ok := app.(map[string]interface{})
	if ok {
		mapstructure.Decode(raw, &slsConfig)
	}

	slsConfig.Name = name

	if slsConfig.Path == "" {
		slsConfig.Path = fmt.Sprintf("./projects/%s", name)
	}

	if len(slsConfig.File) == 0 {
		slsConfig.File = "serverless.yml"
	}
	if len(slsConfig.SLSNodeModuleCacheMount) == 0 {
		slsConfig.SLSNodeModuleCacheMount = fmt.Sprintf("%s-node-modules", slsConfig.Name)
	}

	slsConfig.Env = append(slsConfig.Env, "SLS_DEBUG=*")

	return &slsConfig
}

func (sls *serverless) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s: initializing Docker client...", sls.Name)
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	image := "node:" + sls.NodeVersion

	s.Update("%s: checking for Docker image: %s", sls.Name, image)

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
		s.Update("%s: pulling image: %s", sls.Name, image)

		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		stdout, _, err := ui.OutputWriters()
		if err != nil {
			return err
		}

		var termFd uintptr
		if f, ok := stdout.(*os.File); ok {
			termFd = f.Fd()
		}

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), termFd, true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}

		s.Done()
		s = sg.Add("")
	}

	s.Update("%s: downloading npm modules...", sls.Name)

	err = sls.npm(cli, []string{"npm", "install", "--save-dev"}, s, ui)
	if err != nil {
		return fmt.Errorf("can't deploy %s: %w", sls.Name, err)
	}

	s.Done()

	if sls.CreateDomain {
		s = sg.Add("%s: creating domain...", sls.Name)
		err = sls.serverless(cli, []string{
			"create_domain",
			"--verbose",
			"--region", viper.GetString("aws_region"),
			"--env", viper.GetString("env"),
			"--profile", viper.GetString("aws_profile"),
		}, s)
		if err != nil {
			return err
		}

		s.Done()
	}

	s = sg.Add("%s: deloying app...", sls.Name)

	err = sls.serverless(cli, []string{
		"deploy",
		"--config", sls.File,
		"--service", sls.Name,
		"--verbose",
		"--region", viper.GetString("aws_region"),
		"--env", viper.GetString("env"),
		"--profile", viper.GetString("aws_profile"),
	}, s)
	if err != nil {
		s.Abort()
		return err
	}

	s.Done()
	s = sg.Add("%s: deployment completed!", sls.Name)
	s.Done()

	return nil
}

func (sls *serverless) Destroy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s: initializing Docker client...", sls.Name)
	defer func() { s.Abort() }()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	image := "node:" + sls.NodeVersion

	s.Update("%s: checking for Docker image: %s", sls.Name, image)

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
		s.Update("%s: pulling image: %s", sls.Name, image)

		resp, err := cli.ImagePull(context.Background(), reference.FamiliarString(imageRef), types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer resp.Close()

		stdout, _, err := ui.OutputWriters()
		if err != nil {
			return err
		}

		var termFd uintptr
		if f, ok := stdout.(*os.File); ok {
			termFd = f.Fd()
		}

		err = jsonmessage.DisplayJSONMessagesStream(resp, s.TermOutput(), termFd, true, nil)
		if err != nil {
			return fmt.Errorf("unable to stream pull logs to the terminal: %s", err)
		}

		s.Done()
		s = sg.Add("")
	}

	s.Done()
	s = sg.Add("%s: destroying app...", sls.Name)

	err = sls.serverless(cli, []string{
		"remove",
		"--config", sls.File,
		"--service", sls.Name,
		"--verbose",
		"--region", viper.GetString("aws_region"),
		"--env", viper.GetString("env"),
		"--profile", viper.GetString("aws_profile"),
	}, s)
	if err != nil {
		s.Abort()
		return err
	}

	s.Done()
	s = sg.Add("%s: deployment completed!", sls.Name)
	s.Done()

	return nil
}

func (sls *serverless) serverless(cli *client.Client, cmd []string, step terminal.Step) error {
	command := []string{"serverless"}
	command = append(command, cmd...)

	contConfig := &container.Config{
		Image:        fmt.Sprintf("node:%v", sls.NodeVersion),
		Entrypoint:   strslice.StrSlice{"/usr/local/bin/npx"},
		WorkingDir:   "/app",
		Env:          sls.Env,
		Tty:          true,
		Cmd:          command,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	contHostConfig := sls.getHostConfig()

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

func (sls *serverless) npm(cli *client.Client, cmd []string, s terminal.Step, ui terminal.UI) error {
	contConfig := &container.Config{
		WorkingDir:   "/app",
		Image:        fmt.Sprintf("node:%v", sls.NodeVersion),
		Tty:          true,
		Cmd:          cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
	}

	contHostConfig := sls.getHostConfig()

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

func (sls *serverless) getHostConfig() *container.HostConfig {
	return &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				ReadOnly: true,
				Source:   fmt.Sprintf("%v/.aws", viper.Get("HOME")),
				Target:   "/root/.aws",
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%s/%s", viper.Get("ROOT_DIR"), sls.Path),
				Target: "/app",
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%s/%s/.serverless/", viper.Get("ROOT_DIR"), sls.Path),
				Target: "/root/.config",
			},
			{
				Type:   mount.TypeBind,
				Source: fmt.Sprintf("%s/.npm/", viper.Get("ROOT_DIR")),
				Target: "/root/.npm",
			},
			{
				Type:   mount.TypeVolume,
				Source: sls.SLSNodeModuleCacheMount,
				Target: "/app/node_modules",
			},
		},
	}
}
