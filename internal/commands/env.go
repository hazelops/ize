package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/template"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type envCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newEnvCmd() *envCmd {
	cc := &envCmd{}

	cmd := &cobra.Command{
		Use:              "env",
		Short:            "Manage environment.",
		Long:             "",
		RunE:             nil,
		TraverseChildren: true,
	}

	envCmd := &cobra.Command{
		Use:   "terraform",
		Short: "Generate terraform files.",
		Long:  "This command generate terraform files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			pterm.DefaultSection.Printfln("Starting generate terrafrom files")

			err = template.GenerateBackendTf(template.BackendOpts{
				ENV:                            cc.config.Env,
				LOCALSTACK_ENDPOINT:            "",
				TERRAFORM_STATE_BUCKET_NAME:    fmt.Sprintf("%s-tf-state", cc.config.Namespace),
				TERRAFORM_STATE_KEY:            fmt.Sprintf("%v/terraform.tfstate", cc.config.Env),
				TERRAFORM_STATE_REGION:         cc.config.AwsRegion,
				TERRAFORM_STATE_PROFILE:        cc.config.AwsProfile,
				TERRAFORM_STATE_DYNAMODB_TABLE: "tf-state-lock", // So? // TODO: cc.config.TERRAFORM_STATE_DYNAMODB_TABLE
				TERRAFORM_AWS_PROVIDER_VERSION: "",
			},
				viper.GetString("ENV_DIR"),
			)

			if err != nil {
				pterm.DefaultSection.Println("Generate terrafrom file not completed")
				return err
			}

			pterm.Success.Println("backend.tf generated")

			sess, err := utils.GetSession(&utils.SessionConfig{
				Region:  cc.config.AwsRegion,
				Profile: cc.config.AwsProfile,
			})
			if err != nil {
				pterm.DefaultSection.Println("Generate terrafrom file not completed")
				return err
			}

			pterm.Success.Printfln("Read SSH public key")
			cc.log.Debug("Read SSH public key")

			home, _ := os.UserHomeDir()
			key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", home))
			if err != nil {
				pterm.DefaultSection.Println("Generate terrafrom file not completed")
				return err
			}

			stsSvc := sts.New(sess)

			resp, err := stsSvc.GetCallerIdentity(
				&sts.GetCallerIdentityInput{},
			)

			if err != nil {
				pterm.DefaultSection.Println("Generate terrafrom file not completed")
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
				pterm.DefaultSection.Println("Generate terrafrom file not completed")
				return err
			}

			pterm.Success.Println("terraform.tfvars generated")

			if err != nil {
				pterm.DefaultSection.Println("Generate terrafrom file not completed")
				return err
			}

			pterm.DefaultSection.Printfln("Generate terrafrom files completed")

			return nil
		},
	}

	cmd.AddCommand(envCmd)

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}
