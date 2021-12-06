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
		Use:   "terraform",
		Short: "Terraform management.",
		Long:  "This command contains subcommands for work with terraform.",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Run terraform init.",
		Long:  `This command run terraform init command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
				ContainerName: "terraform-init",
				Cmd:           []string{"init"},
			}

			opts.Cmd = append(opts.Cmd, args...)

			cc.log.Debug("starting terraform init")

			err = runTerraform(cc, opts)
			if err != nil {
				cc.log.Error("terraform init not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform init completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Run terraform apply.",
		Long: `This command run terraform apply command. Terraform apply
		command executes the actions proposed in a Terraform plan.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
				ContainerName: "terraform-apply",
				Cmd:           []string{"apply"},
			}

			opts.Cmd = append(opts.Cmd, args...)

			cc.log.Debug("starting terraform apply")

			err = runTerraform(cc, opts)

			if err != nil {
				cc.log.Error("terraform apply not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform apply completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "Run terraform plan.",
		Long: `This command run terraform plan command.
		The terraform plan command creates an execution plan.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
				ContainerName: "terraform-plan",
				Cmd:           []string{"plan"},
			}

			opts.Cmd = append(opts.Cmd, args...)

			cc.log.Debug("starting terraform plan")

			err = runTerraform(cc, opts)

			if err != nil {
				cc.log.Error("terraform plan not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform plan completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Run terraform destroy.",
		Long:  `This command run terraform destroy command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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
				ContainerName: "terraform-destroy",
				Cmd:           []string{"destroy"},
			}

			opts.Cmd = append(opts.Cmd, args...)

			cc.log.Debug("starting terraform destroy")

			err = runTerraform(cc, opts)

			if err != nil {
				cc.log.Error("terraform destroy not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform destroy completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	cmd.AddCommand(initCmd, applyCmd, destroyCmd, planCmd)
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

	cc.log.Debugf("container config: %s", contConfig)
	cc.log.Debugf("container host config: %s", contHostConfig)

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
		cc.log.Errorf("terraform container started:", cont.ID)
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

	if cc.log.GetLevel() >= 4 {
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
			fmt.Println(scanner.Text())
		}
	}

	if err != nil {
		return err
	}

	cc.log.Debugf("terraform container started:", cont.ID)

	return nil
}
