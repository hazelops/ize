package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"gopkg.in/ini.v1"
)

const (
	path = "/.aws/credentials-mfa"
)

type SessionConfig struct {
	Region  string
	Profile string
}

func GetSession(c *SessionConfig) (*session.Session, error) {
	upd := false

	config := aws.NewConfig().WithRegion(c.Region)

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  *config,
		Profile: c.Profile,
	})
	if err != nil {
		return nil, err
	}

	devices, err := iam.New(sess).ListMFADevices(&iam.ListMFADevicesInput{})
	if err != nil {
		return nil, err
	}

	if len(devices.MFADevices) == 0 {
		return sess, nil
	}

	home, _ := os.UserHomeDir()
	filePath := home + path

	credFile, err := ini.Load(filePath)
	if err != nil {
		credFile = ini.Empty(ini.LoadOptions{})
		upd = true
	}

	var sect *ini.Section
	var exp *ini.Key

	if !upd {
		sect, err = credFile.GetSection(fmt.Sprintf("%s-mfa", c.Profile))
		if err != nil {
			upd = true
		}
	}

	if !upd {
		if len(sect.KeyStrings()) != 4 {
			upd = true
		}
	}

	if !upd {
		exp, err = sect.GetKey("token_expiration")
		if err != nil {
			upd = true
		}
	}

	if !upd {
		timeExp, err := time.Parse("2006-01-02T15:04:05Z07:00", exp.String())
		if err != nil {
			upd = true
		}

		if timeExp.Before(time.Now().UTC()) {
			upd = true
		}
	}

	if upd {
		cred, err := getNewToken(sess, devices.MFADevices[0].SerialNumber)
		if err != nil {
			return nil, err
		}

		err = writeCredsToFile(cred, credFile, filePath, c.Profile)
		if err != nil {
			return nil, err
		}
	}

	sess, err = session.NewSessionWithOptions(
		session.Options{
			Config:            *aws.NewConfig().WithRegion(c.Region),
			Profile:           fmt.Sprintf("%s-mfa", c.Profile),
			SharedConfigFiles: []string{filePath},
		},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil

}

func getNewToken(sess *session.Session, sn *string) (*sts.Credentials, error) {
	stsSvc := sts.New(sess)

	token, err := stscreds.StdinTokenProvider()
	if err != nil {
		return nil, err
	}

	out, err := stsSvc.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber: sn,
		TokenCode:    &token,
	})

	if err != nil {
		return nil, err
	}

	return out.Credentials, nil
}

func writeCredsToFile(creds *sts.Credentials, f *ini.File, filepath, profile string) error {
	sect, err := f.NewSection(fmt.Sprintf("%s-mfa", profile))
	if err != nil {
		return err
	}

	sect.NewKey("aws_access_key_id", *creds.AccessKeyId)
	sect.NewKey("aws_secret_access_key", *creds.SecretAccessKey)
	sect.NewKey("aws_session_token", *creds.SessionToken)
	sect.NewKey("token_expiration", creds.Expiration.Format("2006-01-02T15:04:05Z07:00"))

	err = f.SaveTo(filepath)
	if err != nil {
		return err
	}

	return nil
}
