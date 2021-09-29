package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/pkg/logger"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type terraformCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newTerraformCmd() *terraformCmd {
	cc := &terraformCmd{}

	cmd := &cobra.Command{
		Use:   "terraform",
		Short: "Terraform management.",
		Long:  `This command contains subcommands for work with terraform.`,
		RunE:  nil,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Download terraform docker image",
		Long:  `This command download terraform docker image of the specified version.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			runTerraform()
			return nil
		},
	})

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func runTerraform() {
	l := logger.NewSugaredLogger(zapcore.Level(0))

	fmt.Println("Init docker client")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	imageName := "hashicorp/terraform"
	imageTag := viper.Get("TERRAFORM_VERSION")
	termFd, _ := term.GetFdInfo(os.Stderr)

	fmt.Printf("Started pull terraform image %v:%v", imageName, imageTag)

	out, err := cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	err = jsonmessage.DisplayJSONMessagesStream(out, &l, termFd, true, nil)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Finished pulling terraform image %v:%v \n", imageName, imageTag)

	fmt.Printf("Start creating terraform container from image %v:%v \n", imageName, imageTag)

	//TODO: Add Auto Pull Docker image
	//TODO: Check if such container exists to use fixed name
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
			Tty:          true,
			Cmd:          strings.Split("init", " "),
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
			WorkingDir:   fmt.Sprintf("%v", viper.Get("ENV_DIR")),
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

	fmt.Println("Finished creating terraform container from image", fmt.Sprintf("%v:%v", imageName, imageTag))

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Terraform container started!", cont.ID)
}
