package commands

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

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

			pterm.DefaultSection.Println("Generating terraform files")

			opts := TerraformRunOption{
				ContainerName: "terraform-init",
				Cmd:           []string{"init", "-input=true"},
			}

			pterm.DefaultSection.Println("Starting Terraform init")

			err = runTerraform(cc, opts)
			if err != nil {
				pterm.DefaultSection.Println("Terraform init not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform init completed")

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "apply",
		Short: "Run terraform apply",
		Long: `This command run terraform apply command. Terraform apply
		command executes the actions proposed in a Terraform plan`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()

			if err != nil {
				return err
			}

			opts := TerraformRunOption{
				ContainerName: "terraform-apply",
				Cmd:           []string{"apply", "-input=false", fmt.Sprintf("%v/.terraform/tfplan", viper.Get("ENV_DIR"))},
			}

			pterm.DefaultSection.Println("Starting Terraform apply")

			err = runTerraform(cc, opts)

			if err != nil {
				pterm.DefaultSection.Println("Terraform apply not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform apply completed")

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "plan",
		Short: "Run terraform plan",
		Long: `This command run terraform plan command.
		The terraform plan command creates an execution plan.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()

			if err != nil {
				return err
			}

			opts := TerraformRunOption{
				ContainerName: "terraform-plan",
				Cmd:           []string{"plan", fmt.Sprintf("-out=%v/.terraform/tfplan", viper.Get("ENV_DIR")), "-input=false"},
			}

			pterm.DefaultSection.Println("Starting Terraform plan")

			err = runTerraform(cc, opts)

			if err != nil {
				pterm.DefaultSection.Println("Terraform plan not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform plan completed")

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "destroy",
		Short: "Run terraform destroy",
		Long:  `This command run terraform destroy command.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()

			if err != nil {
				return err
			}

			opts := TerraformRunOption{
				ContainerName: "terraform-destroy",
				Cmd:           []string{"destroy"},
			}

			pterm.DefaultSection.Println("Starting Terraform destroy")

			err = runTerraform(cc, opts)

			if err != nil {
				pterm.DefaultSection.Println("Terraform destroy not completed")
				return err
			}

			pterm.DefaultSection.Println("Terraform destroy completed")

			return nil
		},
	})

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
	termFd, _ := term.GetFdInfo(os.Stderr)

	out, err := cli.ImagePull(context.Background(), fmt.Sprintf("%v:%v", imageName, imageTag), types.ImagePullOptions{})
	if err != nil {
		pterm.Error.Printfln("Pulling terraform image %v:%v/n", imageName, imageTag)
		return err
	}

	pterm.Success.Printfln("Pulling terraform image %v:%v/n", imageName, imageTag)

	err = jsonmessage.DisplayJSONMessagesStream(out, &cc.log, termFd, true, nil)

	if err != nil {
		return err
	}

	cc.log.Debugf(fmt.Sprintf("%v", viper.Get("ENV_DIR")))

	//TODO: Add Auto Pull Docker image
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
		pterm.Error.Printfln("Terraform container started:", cont.ID)
		return err
	}

	statusCh, errCh := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionNextExit)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			out, err = cli.ContainerLogs(context.Background(), cont.ID, types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
			})
			if err != nil {
				return err
			}

			defer out.Close()
			content, _ := ioutil.ReadAll(out)
			pterm.Error.Printfln("Terraform container started: %s", cont.ID)

			return errors.New(string(content))
		}
	}

	pterm.Success.Printfln("Terraform container started: %s", cont.ID)

	return nil
}
