package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"
)

func TerraformInit()  {
	command := "init"
	runTerraform(command)

}

func runTerraform(command string)  {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	//TODO: Add Auto Pull Docker image
	//TODO: Check if such container exists to use fixed name
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "hashicorp/terraform",
			Tty: true,
			Cmd: strings.Split(command, " "),
			AttachStdin: true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin: true,

		},
		&container.HostConfig{
			AutoRemove: false,
		}, nil, nil, "")

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: false})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)



}