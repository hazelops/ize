package gomplate

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/hazelops/ize/pkg/logger"
)

const (
	imageName = "hairyhenderson/gomplate"
	imageTag  = "latest"
)

// Options for run gomplate
type GomplateOptions struct {
	OutputFileDir  string
	InputFileDir   string
	InputFileName  string
	OutputFileName string
	Env            []string
}

// Run gomplate container and generate file
func RunGomplate(opts GomplateOptions, log logger.StandartLogger) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	_, err = cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		return err
	}

	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		Cmd: strslice.StrSlice{"-f", fmt.Sprintf("%v/%v", opts.InputFileDir, opts.InputFileName),
			"-o", fmt.Sprintf("%v/%v", opts.OutputFileDir, opts.OutputFileName)},
		Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		WorkingDir:   opts.InputFileDir,
		Env:          opts.Env,
	},
		&container.HostConfig{

			AutoRemove: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: fmt.Sprintf("%v", opts.InputFileDir),
					Target: fmt.Sprintf("%v", opts.InputFileDir),
				},
				{
					Type:   mount.TypeBind,
					Source: fmt.Sprintf("%v", opts.OutputFileDir),
					Target: fmt.Sprintf("%v", opts.OutputFileDir),
				},
			},
		},
		nil, nil, "",
	)

	if err != nil {
		return err
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
