package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/template"
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

			err = template.GenereateBackendTf(template.BackendOpts{
				ENV:                            cc.config.Env,
				LOCALSTACK_ENDPOINT:            "",
				TERRAFORM_STATE_BUCKET_NAME:    fmt.Sprintf("%s-tf-state", cc.config.Namespace),
				TERRAFORM_STATE_KEY:            fmt.Sprintf("%v/terraform.tfstate", cc.config.Env),
				TERRAFORM_STATE_REGION:         cc.config.AwsRegion,
				TERRAFORM_STATE_PROFILE:        cc.config.AwsProfile,
				TERRAFORM_STATE_DYNAMODB_TABLE: "tf-state-lock", // So?
				TERRAFORM_AWS_PROVIDER_VERSION: "",
			},
				viper.GetString("ENV_DIR"),
			)

			pterm.Success.Println("backend.tf generated")

			if err != nil {
				pterm.Error.Println("backend.tf not generated")
				return err
			}

			sess, err := utils.GetSession(&utils.SessionConfig{
				Region: cc.config.AwsRegion,
			})
			if err != nil {
				return err
			}

			pterm.Success.Printfln("Read SSH public key")
			cc.log.Debug("Read SSH public key")

			key, err := ioutil.ReadFile("/home/psih/.ssh/id_rsa.pub")
			if err != nil {
				return err
			}

			stsSvc := sts.New(sess)

			resp, err := stsSvc.GetCallerIdentity(
				&sts.GetCallerIdentityInput{},
			)

			if err != nil {
				return err
			}

			err = template.GenerateVarsTf(template.VarsOpts{
				ENV:               cc.config.Env,
				AWS_PROFILE:       cc.config.AwsProfile,
				AWS_REGION:        cc.config.AwsRegion,
				EC2_KEY_PAIR_NAME: fmt.Sprintf("%v-%v", cc.config.Env, cc.config.Namespace),
				TAG:               cc.config.Env,
				SSH_PUBLIC_KEY:    string(key)[:len(string(key))-1],
				DOCKER_REGISTRY:   fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, cc.config.AwsRegion),
				NAMESPACE:         cc.config.Namespace,
			},
				viper.GetString("ENV_DIR"),
			)

			if err != nil {
				pterm.Error.Println("terraform.tfvars not generated")
				return err
			}

			pterm.Success.Println("terraform.tfvars generated")

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

			err = runTerraform(cc, opts)

			if err != nil {
				return err
			}

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
	pterm.Success.Println("Init docker client")
	cc.log.Debug("Init docker client")

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	imageName := "hashicorp/terraform"
	pterm.Info.Printfln("cfg", *(cc.cfg))
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

	ps, err := pterm.DefaultSpinner.Start("Check existing container")
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: true, // include stopped containers
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: opts.ContainerName,
		}),
	})
	if err != nil {
		return err
	}

	if len(containers) > 0 {
		ps.Success("Start exist container")

		if err := cli.ContainerStart(context.Background(), containers[0].ID, types.ContainerStartOptions{}); err != nil {
			pterm.Error.Printfln("Container start:", err)
			return err
		}

		return nil
	}

	pterm.Success.Printfln("Finished pulling terraform image %v:%v", imageName, imageTag)

	pterm.Success.Printfln("Start creating terraform container from image %v:%v", imageName, imageTag)

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
				fmt.Sprintf("ENV=%v", cc.cfg.Env),
				fmt.Sprintf("AWS_PROFILE=%v", cc.cfg.AwsProfile),
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
		fmt.Println(err)
		return err
	}

	pterm.Success.Printfln("Finished creating terraform container from image %v:%v", imageName, imageTag)

	if err := cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{}); err != nil {
		pterm.Error.Printfln("Container start:", err)
		return err
	}

	pterm.Success.Printfln("Terraform container started!", cont.ID)

	return nil
}
