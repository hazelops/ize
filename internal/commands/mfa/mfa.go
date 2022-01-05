package mfa

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type MfaOptions struct {
	Profile string
	Region  string
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
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
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
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Profile = viper.GetString("aws_profile")
	o.Region = viper.GetString("aws_region")

	if o.Region == "" {
		o.Region = viper.GetString("aws-region")
	}

	if o.Profile == "" {
		o.Profile = viper.GetString("aws-profile")
	}

	return nil
}

func (o *MfaOptions) Validate() error {
	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("AWS region must be specified")
	}
	return nil
}

func (o *MfaOptions) Run() error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
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
