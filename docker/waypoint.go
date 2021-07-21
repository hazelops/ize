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

// TODO: get current directory
// TODO: Be able to get current directory globally

func WaypointInit()  {
	command := "init"
	RunWaypoint(command)

}


func RunWaypoint(command string)  {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	//TODO: Check if such container exists to use fixed name
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "hazelops/waypoint",
			Tty: true,
			Cmd: strings.Split(command, " "),
			AttachStdin: true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin: true,

		},

		&container.HostConfig{
			AutoRemove: true,
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