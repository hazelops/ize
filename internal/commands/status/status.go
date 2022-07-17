package status

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

func NewDebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show debug information",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			c, err := config.GetConfig()
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("can't load options for a command: %w", err)
			}

			dt := pterm.DefaultTable

			pterm.DefaultSection.Println("IZE Info")

			dt.WithData(pterm.TableData{
				{"ENV", c.Env},
				{"NAMESPACE", c.Namespace},
				{"TAG", c.Tag},
				{"INFRA DIR", c.InfraDir},
				{"PWD", cwd},
				{"IZE VERSION", version.Version},
				{"GIT REVISION", version.GitCommit},
				{"ENV DIR", c.EnvDir},
				{"PREFER_RUNTIME", c.PreferRuntime},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("Terraform Info")
			dt.WithData(pterm.TableData{
				{"TERRAFORM_VERSION", c.TerraformVersion},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("System Info")

			dt.WithData(pterm.TableData{
				{"OS", runtime.GOOS},
				{"ARCH", runtime.GOARCH},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("AWS Environment Info")

			if len(c.AwsProfile) > 0 {
				resp, err := sts.New(c.Session).GetCallerIdentity(
					&sts.GetCallerIdentityInput{},
				)
				if err != nil {
					return err
				}

				guo, err := iam.New(c.Session).GetUser(&iam.GetUserInput{})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with aws_profile=%s: username %s is not found in account %s", c.AwsProfile, *guo.User.UserName, *resp.Account)
					default:
						return err
					}
				}

				luto, err := iam.New(c.Session).ListUserTags(&iam.ListUserTagsInput{
					UserName: guo.User.UserName,
				})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with aws_profile=%s: username %s is not found in account %s", c.AwsProfile, *guo.User.UserName, *resp.Account)
					default:
						return err
					}
				}

				devEnvName := ""

				for _, k := range luto.Tags {
					if *k.Key == "devEnvironmentName" {
						devEnvName = *k.Value
					}
				}

				dt.WithData(pterm.TableData{
					{"AWS PROFILE", c.AwsProfile},
					{"AWS USER", *guo.User.UserName},
					{"AWS ACCOUNT", *resp.Account},
				}).WithLeftAlignment().Render()

				if len(devEnvName) > 0 {
					dt.WithData(pterm.TableData{
						{"AWS_DEV_ENV_NAME", devEnvName},
					}).WithLeftAlignment().Render()
				}
			} else {
				pterm.Println("No AWS profile credentials detected. Parameters used:")
				dt.WithData(pterm.TableData{
					{"AWS PROFILE", c.AwsProfile},
				}).WithLeftAlignment().Render()
			}

			return nil
		},
	}

	return cmd
}
