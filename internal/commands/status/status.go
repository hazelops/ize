package status

import (
	"fmt"
	"os"
	"runtime"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
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

			sess, err := utils.GetSession(&utils.SessionConfig{
				Region:  viper.GetString("aws_region"),
				Profile: viper.GetString("aws_profile"),
			})
			if err != nil {
				return err
			}

			resp, err := sts.New(sess).GetCallerIdentity(
				&sts.GetCallerIdentityInput{},
			)
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
				{"ENV", viper.GetString("env")},
				{"TAG", viper.GetString("tag")},
				{"INFRA DIR", viper.GetString("infra_dir")},
				{"PWD", cwd},
				{"IZE VERSION", version.Version},
				{"IZE GIT REVISION", version.GitCommit},
				{"ENV DIR", viper.GetString("env_dir")},
			}).WithLeftAlignment().Render()

			pterm.DefaultSection.Println("System Info")

			dt.WithData(pterm.TableData{
				{"OS", runtime.GOOS},
				{"ARCH", runtime.GOARCH},
			}).WithLeftAlignment().Render()

			guo, err := iam.New(sess).GetUser(&iam.GetUserInput{})
			if err != nil {
				return err
			}

			luto, err := iam.New(sess).ListUserTags(&iam.ListUserTagsInput{
				UserName: guo.User.UserName,
			})
			if err != nil {
				return err
			}

			devEnvName := ""

			for _, k := range luto.Tags {
				if *k.Key == "devEnvironmentName" {
					devEnvName = *k.Value
				}
			}

			pterm.DefaultSection.Println("AWS Environment Info")

			dt.WithData(pterm.TableData{
				{"AWS_DEV_ENV_NAME", devEnvName},
				{"AWS PROFILE", fmt.Sprintf("%s-%s", viper.GetString("env"), viper.GetString("namespace"))},
				{"AWS USER", *guo.User.UserName},
				{"AWS ACCOUNT", *resp.Account},
			}).WithLeftAlignment().Render()

			return nil
		},
	}

	return cmd
}
