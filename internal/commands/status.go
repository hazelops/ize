package commands

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/config"

	"github.com/hazelops/ize/internal/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func NewDebugCmd(project *config.Project) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show debug information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("can't load options for a command: %w", err)
			}

			dt := pterm.DefaultTable

			pterm.DefaultSection.Println("IZE Info")

			_ = dt.WithData(pterm.TableData{
				{"ENV", project.Env},
				{"NAMESPACE", project.Namespace},
				{"TAG", project.Tag},
				{"INFRA DIR", project.InfraDir},
				{"PWD", cwd},
				{"IZE VERSION", version.FullVersionNumber()},
				{"GIT REVISION", version.GitCommit},
				{"ENV DIR", project.EnvDir},
				{"PREFER_RUNTIME", project.PreferRuntime},
			}).WithLeftAlignment().Render()

			v := project.TerraformVersion
			if project.Terraform != nil {
				if i, ok := project.Terraform["infra"]; ok {
					if len(i.Version) != 0 {
						v = i.Version
					}
				}
			}

			pterm.DefaultSection.Println("Terraform Info")
			_ = dt.WithData(pterm.TableData{
				{"TERRAFORM_VERSION", v},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("System Info")

			_ = dt.WithData(pterm.TableData{
				{"OS", runtime.GOOS},
				{"ARCH", runtime.GOARCH},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("AWS Environment Info")

			if len(project.AwsProfile) > 0 {
				resp, err := project.AWSClient.STSClient.GetCallerIdentity(
					&sts.GetCallerIdentityInput{},
				)
				if err != nil {
					return err
				}

				guo, err := project.AWSClient.IAMClient.GetUser(&iam.GetUserInput{})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with aws_profile=%s: username %s is not found in account %s", project.AwsProfile, *guo.User.UserName, *resp.Account)
					default:
						return err
					}
				}

				tags, err := project.AWSClient.IAMClient.ListUserTags(&iam.ListUserTagsInput{
					UserName: guo.User.UserName,
				})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with aws_profile=%s: username %s is not found in account %s", project.AwsProfile, *guo.User.UserName, *resp.Account)
					default:
						return err
					}
				}

				devEnvName := ""

				for _, k := range tags.Tags {
					if *k.Key == "devEnvironmentName" {
						devEnvName = *k.Value
					}
				}

				_ = dt.WithData(pterm.TableData{
					{"AWS PROFILE", project.AwsProfile},
					{"AWS USER", *guo.User.UserName},
					{"AWS ACCOUNT", *resp.Account},
				}).WithLeftAlignment().Render()

				if len(devEnvName) > 0 {
					_ = dt.WithData(pterm.TableData{
						{"AWS_DEV_ENV_NAME", devEnvName},
					}).WithLeftAlignment().Render()
				}
			} else {
				pterm.Println("No AWS profile credentials detected. Parameters used:")
				_ = dt.WithData(pterm.TableData{
					{"AWS PROFILE", project.AwsProfile},
				}).WithLeftAlignment().Render()
			}

			return nil
		},
	}

	return cmd
}
