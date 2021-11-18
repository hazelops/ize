package commands

import (
	"fmt"

	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/spf13/cobra"
)

type mfaCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newMfaCmd() *mfaCmd {
	cc := &mfaCmd{}

	cmd := &cobra.Command{
		Use:              "mfa",
		Short:            "MFA management.",
		RunE:             nil,
		TraverseChildren: true,
	}

	mfaCmd := &cobra.Command{
		Use:   "export",
		Short: "Print mfa credentials.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			sess, err := utils.GetSession(&utils.SessionConfig{
				Region:  cc.config.AwsRegion,
				Profile: cc.config.AwsProfile,
			})
			if err != nil {
				return err
			}

			v, err := sess.Config.Credentials.Get()
			if err != nil {
				return err
			}

			fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s && \nexport AWS_SESSION_TOKEN=%s && \nexport AWS_ACCESS_KEY_ID=%s",
				v.SecretAccessKey, v.SessionToken, v.AccessKeyID,
			)

			return nil
		},
	}

	cmd.AddCommand(mfaCmd)

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}
