package mfa

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type MfaOptions struct {
	Config *config.Config
}

func NewMfaFlags() *MfaOptions {
	return &MfaOptions{}
}

func NewCmdMfa() *cobra.Command {
	o := NewMfaFlags()

	cmd := &cobra.Command{
		Use:   "mfa",
		Short: "generate terraform files",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Run()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *MfaOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	return nil
}

func (o *MfaOptions) Run() error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		return err
	}

	devices, err := iam.New(sess).ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		return err
	}

	if len(devices.MFADevices) == 0 {
		logrus.Error("MFA hasn't configured")
		return fmt.Errorf("MFA hasn't configured")
	}

	v, err := sess.Config.Credentials.Get()
	if err != nil {
		return err
	}

	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s && \\ \nexport AWS_SESSION_TOKEN=%s && \\ \nexport AWS_ACCESS_KEY_ID=%s",
		v.SecretAccessKey, v.SessionToken, v.AccessKeyID,
	)

	return nil
}
