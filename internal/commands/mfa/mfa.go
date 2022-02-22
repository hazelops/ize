package mfa

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hazelops/ize/internal/config"
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
		Short: "return export command with AWS MFA credentials",
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
	devices, err := iam.New(o.Config.Session).ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		return err
	}

	if len(devices.MFADevices) == 0 {
		return fmt.Errorf("MFA hasn't configured\n")
	}

	v, err := o.Config.Session.Config.Credentials.Get()
	if err != nil {
		return err
	}

	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s && \\ \nexport AWS_SESSION_TOKEN=%s && \\ \nexport AWS_ACCESS_KEY_ID=%s\n",
		v.SecretAccessKey, v.SessionToken, v.AccessKeyID,
	)

	return nil
}
