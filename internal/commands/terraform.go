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

			sess, err := utils.GetSession(&utils.SessionConfig{
				Region: cc.cfg.AwsRegion,
			})
			if err != nil {
				return err
			}

			err = gomplate.RunGomplate(gomplate.GomplateOptions{
				OutputFileDir:  viper.GetString("ENV_DIR"),
				InputFileDir:   fmt.Sprintf("%v/terraform/template/", viper.GetString("INFRA_DIR")),
				InputFileName:  "backend.tf.gotmpl",
				OutputFileName: "backend.tf",
				Env: []string{
					fmt.Sprintf("TERRAFORM_STATE_KEY=%v/terraform.tfstate", cc.cfg.Env),
					fmt.Sprintf("TERRAFORM_STATE_BUCKET_NAME=%v-tf-state", cc.cfg.Namespace),
					fmt.Sprintf("TERRAFORM_STATE_REGION=%v", cc.cfg.AwsRegion),
					fmt.Sprintf("TERRAFORM_STATE_PROFILE=%v", cc.cfg.AwsProfile),
					fmt.Sprintf("TERRAFORM_STATE_DYNAMODB_TABLE=%v", "tf-state-lock"),
					fmt.Sprintf("ENV=%v", viper.Get("ENV")),
					fmt.Sprintf("AWS_PROFILE=%v", viper.Get("AWS_PROFILE")),
					fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
					fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
				},
			}, cc.log)
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

			key, err := ioutil.ReadFile("/home/psih/.ssh/id_rsa.pub")
			if err != nil {
				return err
			}

			err = gomplate.RunGomplate(gomplate.GomplateOptions{
				OutputFileDir:  viper.GetString("ENV_DIR"),
				InputFileDir:   fmt.Sprintf("%v/terraform/template/", viper.GetString("INFRA_DIR")),
				InputFileName:  "terraform.tfvars.gotmpl",
				OutputFileName: "terraform.tfvars",
				Env: []string{
					fmt.Sprintf("ENV=%v", cc.cfg.Env),
					fmt.Sprintf("AWS_PROFILE=%v", cc.cfg.AwsProfile),
					fmt.Sprintf("AWS_REGION=%v", cc.cfg.AwsRegion),
					fmt.Sprintf("EC2_KEY_PAIR_NAME=%v-%v", cc.cfg.Env, cc.cfg.Namespace),
					fmt.Sprintf("DOCKER_REGISTRY=%v.dkr.ecr.%v.amazonaws.com", *resp.Account, cc.cfg.AwsRegion),
					fmt.Sprintf("TAG=%v", cc.cfg.Env),
					fmt.Sprintf("SSH_PUBLIC_KEY=%s", string(key)[:len(string(key))-1]),
					fmt.Sprintf("NAMESPACE=%v", cc.cfg.Namespace),
					fmt.Sprintf("TF_LOG=%v", viper.Get("TF_LOG")),
					fmt.Sprintf("TF_LOG_PATH=%v", viper.Get("TF_LOG_PATH")),
				},
			}, cc.log)
			if err != nil {
				return err
			}

			opts := TerraformRunOption{
				ContainerName: "terraform-init",
				Cmd:           []string{"init", "-input=true"},
			}

			runTerraform(cc, opts)
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
