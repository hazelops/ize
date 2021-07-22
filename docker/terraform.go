package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/hazelops/ize/tpl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

func TerraformInit()  {
	//command := "init"
	tpl.GenerateBackendTf()
	//runTerraform(command)
}

func TerraformPlan()  {
	command := "plan"

	runTerraform(command)
}

func runTerraform(command string)  {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	imageName := "hashicorp/terraform"
	imageTag := viper.Get("TERRAFORM_VERSION")

	//TODO: Add Auto Pull Docker image
	//TODO: Check if such container exists to use fixed name
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: fmt.Sprintf("%v:%v", imageName, imageTag),
			Tty: true,
			Cmd: strings.Split(command, " "),
			AttachStdin: true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin: true,
			WorkingDir: fmt.Sprintf("%v", viper.Get("ENV_DIR")),
			Env: []string{
				fmt.Sprintf("ENV=%v", viper.Get("ENV")),
				fmt.Sprintf("AWS_PROFILE=%v", viper.Get("AWS_PROFILE")),
				fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
				fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
			},
		},

		&container.HostConfig{
			AutoRemove: false,
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
					Source: fmt.Sprintf("%v", viper.Get("HOME")),
					Target: fmt.Sprintf("%v", viper.Get("HOME")),
				},
			},

		}, nil, nil, "")

	var out io.ReadCloser

	out, err = cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, out)
	cobra.CheckErr(err)


	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	out, err = cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: false})
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, out)
	cobra.CheckErr(err)

}
