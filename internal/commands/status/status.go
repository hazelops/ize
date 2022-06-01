package status

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
	"os"
	"runtime"

	"github.com/hazelops/ize/internal/version"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDebugCmd() *cobra.Command {
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

			dt.WithData(pterm.TableData{
				{"ENV", viper.GetString("env")},
				{"NAMESPACE", viper.GetString("namepsace")},
				{"TAG", viper.GetString("tag")},
				{"INFRA DIR", viper.GetString("infra_dir")},
				{"PWD", cwd},
				{"IZE VERSION", version.Version},
				{"GIT REVISION", version.GitCommit},
				{"ENV DIR", viper.GetString("env_dir")},
				{"PREFER_RUNTIME", viper.GetString("prefer_runtime")},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("Terrafor Info")
			dt.WithData(pterm.TableData{
				{"TERRAFORM_VERSION", viper.GetString("terraform_version")},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("System Info")

			dt.WithData(pterm.TableData{
				{"OS", runtime.GOOS},
				{"ARCH", runtime.GOARCH},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("AWS Environment Info")

			if len(viper.GetString("aws_profile")) > 1 {
				sess, err := utils.GetSession(&utils.SessionConfig{
					Region:  viper.GetString("aws_region"),
					Profile: viper.GetString("aws_profile"),
				})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoCredentialProviders":
						return fmt.Errorf("Error estabilishing a session with AWS. Please make sure your credentials are valid. Using aws_profile: %s", viper.GetString("aws_profile"))
					default:
						return err
					}
				}

				resp, err := sts.New(sess).GetCallerIdentity(
					&sts.GetCallerIdentityInput{},
				)
				if err != nil {
					return err
				}

				guo, err := iam.New(sess).GetUser(&iam.GetUserInput{})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with aws_profile=%s: username %s is not found in account %s", viper.GetString("aws_profile"), *guo.User.UserName, *resp.Account)
					default:
						return err
					}
				}

				luto, err := iam.New(sess).ListUserTags(&iam.ListUserTagsInput{
					UserName: guo.User.UserName,
				})
				if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case "NoSuchEntity":
						return fmt.Errorf("error obtaining AWS user with %s aws profile: %s is not found in account %s", viper.GetString("aws_profile"), *guo.User.UserName, *resp.Account)
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
					{"AWS PROFILE", viper.GetString("aws_profile")},
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
					{"AWS PROFILE", viper.GetString("aws_profile")},
				}).WithLeftAlignment().Render()
			}

			return nil
		},
	}

	return cmd
}
