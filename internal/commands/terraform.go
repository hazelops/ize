package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type terraformCmd struct {
	*baseBuilderCmd
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

func (b *commandsBuilder) newTerraformCmd() *terraformCmd {
	cc := &terraformCmd{}

	cmd := &cobra.Command{
		Use:   "terraform",
		Short: "Terraform management.",
		Long:  `This command contains subcommands for work with terraform.`,
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Download terraform docker image",
		Long:  `This command download terraform docker image of the specified version.`,
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

			pterm.DefaultSection.Println("Starting Terraform init")

			err = runTerraform(cc, opts)
			if err != nil {
				pterm.DefaultSection.Println("Terraform init not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform init completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Run terraform apply",
		Long: `This command run terraform apply command. Terraform apply
		command executes the actions proposed in a Terraform plan`,
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

			pterm.DefaultSection.Println("Starting Terraform apply")

			err = runTerraform(cc, opts)

			if err != nil {
				pterm.DefaultSection.Println("Terraform apply not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform apply completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "Run terraform plan",
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

			pterm.DefaultSection.Println("Starting Terraform plan")

			err = runTerraform(cc, opts)

			if err != nil {
				pterm.DefaultSection.Println("Terraform plan not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform plan completed")

			return nil
		},
		DisableFlagParsing: true,
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "Run terraform destroy",
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

			pterm.DefaultSection.Println("Starting Terraform destroy")

			err = runTerraform(cc, opts)

			if err != nil {
				pterm.DefaultSection.Println("Terraform destroy not completed")
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

func runTerraform(cc *terraformCmd, opts TerraformRunOption) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		pterm.Error.Println("Docker Clinet initialization")
		return err
	}

	pterm.Success.Println("Docker Clinet initialization")

	imageName := "hashicorp/terraform"
	imageTag := cc.config.TerraformVersion

	out, err := cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		pterm.Error.Printfln("Pulling terraform image %v:%v/n", imageName, imageTag)
		return err
	}

	pterm.Success.Printfln("Pulling terraform image %v:%v/n", imageName, imageTag)

	if cc.log.SugaredLogger != nil {
		termFd, _ := term.GetFdInfo(os.Stderr)
		err = jsonmessage.DisplayJSONMessagesStream(out, &cc.log, termFd, true, nil)
		if err != nil {
			return err
		}
	}

	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
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
		},

		&container.HostConfig{
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
		}, nil, nil, opts.ContainerName)

	if err != nil {
		pterm.Error.Printfln("Creating terraform container from image %v:%v", imageName, imageTag)
		return err
	}

	pterm.Success.Printfln("Creating terraform container from image %v:%v", imageName, imageTag)

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		pterm.Error.Println("Terraform container started:", cont.ID)
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
			strErr = strings.TrimRight(strErr, ".")
			strErr = strings.TrimPrefix(strErr, "Error: ")
			strErr = strings.ToLower(string(strErr[0])) + strErr[1:]
			err = fmt.Errorf(strErr)
		}
		fmt.Println(scanner.Text())
	}

	if err != nil {
		return err
	}

	pterm.Success.Printfln("Terraform container started: %s", cont.ID)

	return nil
}
