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
				{"IZE VERSION", version.Version},
				{"GIT REVISION", version.GitCommit},
				{"ENV DIR", project.EnvDir},
				{"PREFER_RUNTIME", project.PreferRuntime},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("Terraform Info")
			_ = dt.WithData(pterm.TableData{
				{"TERRAFORM_VERSION", project.TerraformVersion},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("System Info")

			_ = dt.WithData(pterm.TableData{
				{"OS", runtime.GOOS},
				{"ARCH", runtime.GOARCH},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("AWS Environment Info")

			if len(project.AwsProfile) > 0 {
				resp, err := sts.New(project.Session).GetCallerIdentity(
					&sts.GetCallerIdentityInput{},
				)
				if err != nil {
					return err
				}

				guo, err := iam.New(project.Session).GetUser(&iam.GetUserInput{})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with aws_profile=%s: username %s is not found in account %s", project.AwsProfile, *guo.User.UserName, *resp.Account)
					default:
						return err
					}
				}

				tags, err := iam.New(project.Session).ListUserTags(&iam.ListUserTagsInput{
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