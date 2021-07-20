package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunWaypoint()  {


	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	//TODO: Check if such container exists to use fixed name
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "hazelops/waypoint",
		},
		&container.HostConfig{
			AutoRemove: true,
		}, nil, nil, "")

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
	fmt.Println("Started successfully Waypoint")

}
func ListContainers() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}