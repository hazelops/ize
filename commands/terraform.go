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
	"github.com/hazelops/ize/pkg/gomplate"
	"github.com/moby/term"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			cc.log.Debug("Init Run Terrafrom Init")
			err := cc.Init()

			if err != nil {
				return err
			}

			err = gomplate.RunGomplate(gomplate.GomplateOptions{
				OutputFileDir:  viper.GetString("ENV_DIR"),
				InputFileDir:   fmt.Sprintf("%v/terraform/template/", viper.GetString("INFRA_DIR")),
				InputFileName:  "backend.tf.gotmpl",
				OutputFileName: "backend.tf",
				Env: []string{
					fmt.Sprintf("ENV=%v", viper.Get("ENV")),
					fmt.Sprintf("AWS_PROFILE=%v", viper.Get("AWS_PROFILE")),
					fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
					fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
				},
			}, cc.log)
			if err != nil {
				return err
			}

			err = gomplate.RunGomplate(gomplate.GomplateOptions{
				OutputFileDir:  viper.GetString("ENV_DIR"),
				InputFileDir:   fmt.Sprintf("%v/terraform/template/", viper.GetString("INFRA_DIR")),
				InputFileName:  "terraform.tfvars.gotmpl",
				OutputFileName: "terraform.tfvars",
				Env: []string{
					fmt.Sprintf("ENV=%v", viper.Get("ENV")),
					fmt.Sprintf("AWS_PROFILE=%v", viper.Get("AWS_PROFILE")),
					fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
					fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
				},
			}, cc.log)
			if err != nil {
				return err
			}

			runTerraform(cc)
			return nil
		},
	})

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func runTerraform(cc *terraformCmd) error {
	pterm.Success.Println("Init docker client")
	cc.log.Debug("Init docker client")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	imageName := "hashicorp/terraform"
	imageTag := cc.cfg.TerraformVersion
	termFd, _ := term.GetFdInfo(os.Stderr)

	pterm.Success.Printfln("Started pull terraform image %v:%v", imageName, imageTag)

	out, err := cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		return err
	}

	err = jsonmessage.DisplayJSONMessagesStream(out, &cc.log, termFd, true, nil)

	if err != nil {
		return err
	}

	pterm.Success.Printfln("Finished pulling terraform image %v:%v", imageName, imageTag)

	pterm.Success.Printfln("Start creating terraform container from image %v:%v", imageName, imageTag)

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
		}, nil, nil, "terraform")

	pterm.Success.Printfln("Finished creating terraform container from image %v:%v", imageName, imageTag)

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	pterm.Success.Printfln("Terraform container started!", cont.ID)

	return nil
}
