package commands

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type terraformCmd struct {
	*baseBuilderCmd
}

const ansi = `\x1B(?:[@-Z\\-_]|\[[0-?]*[-\]*[@-~])`

func (b *commandsBuilder) newTerraformCmd() *terraformCmd {
	cc := &terraformCmd{}

	cmd := &cobra.Command{
		Use:                "terraform <terraform command> [terraform flags]",
		Short:              "terraform management",
		Long:               "This command contains subcommands for work with terraform.",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			if len(args) != 0 {
				if args[0] == "-h" || args[0] == "--help" {
					return cmd.Help()
				}
			}

			err := cc.Init()
			if err != nil {
				return err
			}

			opts := TerraformRunOption{
				ContainerName: "terraform",
				Cmd:           args,
			}

			cc.log.Debug("starting terraform")

			err = runTerraform(cc, opts)
			if err != nil {
				cc.log.Errorf("terraform %s not completed", args[0])
				return err
			}

			pterm.DefaultSection.Printfln("Terraform %s completed", args[0])

			return nil
		},
	}

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

type TerraformRunOption struct {
	ContainerName string
	Cmd           []string
}

func cleanupOldContainers(cli *client.Client, opts TerraformRunOption) error {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if strings.Contains(container.Names[0], opts.ContainerName) {
			err = cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func runTerraform(cc *terraformCmd, opts TerraformRunOption) error {
	cc.log.Debugf("terrafrom run options: %s", opts)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		cc.log.Error("docker сlient initialization")
		return err
	}

	cc.log.Debug("docker сlient initialization")

	err = cleanupOldContainers(cli, opts)
	if err != nil {
		return err
	}

	cc.log.Debug("cleanup old containers successfully")

	imageName := "hashicorp/terraform"
	imageTag := cc.config.TerraformVersion

	cc.log.Infof("image name: %s, image tag: %s", imageName, imageTag)

	out, err := cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		cc.log.Errorf("pulling terraform image %v:%v", imageName, imageTag)
		return err
	}

	wr := ioutil.Discard
	if cc.log.GetLevel() >= 4 {
		wr = os.Stdout
	}

	var termFd uintptr

	err = jsonmessage.DisplayJSONMessagesStream(
		out,
		wr,
		termFd,
		true,
		nil,
	)
	if err != nil {
		cc.log.Errorf("pulling terraform image %v:%v", imageName, imageTag)
		return err
	}

	cc.log.Debugf("pulling terraform image %v:%v", imageName, imageTag)

	contConfig := &container.Config{
		Image:        fmt.Sprintf("%v:%v", imageName, imageTag),
		Tty:          true,
		Cmd:          opts.Cmd,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		WorkingDir:   fmt.Sprintf("%v", viper.Get("ENV_DIR")),
		Env: []string{
			fmt.Sprintf("ENV=%v", cc.config.Env),
			fmt.Sprintf("AWS_PROFILE=%v", cc.config.AwsProfile),
			fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
			fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
		},
	}

	contHostConfig := &container.HostConfig{
		AutoRemove: true,
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
				Source: fmt.Sprintf("%v/.aws", viper.Get("HOME")),
				Target: "/root/.aws",
			},
		},
	}

	cont, err := cli.ContainerCreate(
		context.Background(),
		contConfig,
		contHostConfig,
		nil,
		nil,
		opts.ContainerName,
	)

	if err != nil {
		cc.log.Errorf("creating terraform container from image %v:%v", imageName, imageTag)
		return err
	}

	cc.log.Debugf("creating terraform container from image %v:%v", imageName, imageTag)

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		cc.log.Errorf("terraform container started: %s", cont.ID)
		return err
	}

	reader, err := cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Error: ") {
			r := regexp.MustCompile(ansi)
			strErr := r.ReplaceAllString(scanner.Text(), "")
			strErr = strErr[strings.LastIndex(strErr, "Error: "):]
			strErr = strings.TrimPrefix(strErr, "Error: ")
			strErr = strings.ToLower(string(strErr[0])) + strErr[1:]
			err = fmt.Errorf(strErr)
		}
		if cc.log.GetLevel() >= 4 {
			fmt.Println(scanner.Text())
		}
	}

	if err != nil {
		return err
	}

	cc.log.Debugf("terraform container started: %s", cont.ID)

	return nil
}
