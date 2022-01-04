package tunnel

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelSSHKeyOptions struct {
	Env           string
	Region        string
	Profile       string
	PublicKeyFile string
}

func NewSSHKeyFlags() *TunnelSSHKeyOptions {
	return &TunnelSSHKeyOptions{}
}

func NewCmdSSHKey() *cobra.Command {
	o := NewSSHKeyFlags()

	cmd := &cobra.Command{
		Use:   "ssh-key",
		Short: "send ssh key to remote server",
		Long:  "Send ssh key to remote server.",
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

	cmd.Flags().StringVar(&o.PublicKeyFile, "ssh-public-key", "", "set ssh key public path")

	return cmd
}

func (o *TunnelSSHKeyOptions) Complete(cmd *cobra.Command, args []string) error {
	err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Env = viper.GetString("env")
	o.Profile = viper.GetString("aws-profile")
	o.Region = viper.GetString("aws-region")

	if o.Profile == "" {
		o.Profile = viper.GetString("aws_profile")
	}

	if o.Region == "" {
		o.Region = viper.GetString("aws_region")
	}

	if o.PublicKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PublicKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	return nil
}

func (o *TunnelSSHKeyOptions) Validate() error {
	if len(o.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	if len(o.Profile) == 0 {
		return fmt.Errorf("AWS profile must be specified")
	}

	if len(o.Region) == 0 {
		return fmt.Errorf("AWS region must be specified")
	}
	return nil
}

func (o *TunnelSSHKeyOptions) Run() error {
	pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
	})
	if err != nil {
		pterm.Error.Printfln("Getting AWS session")
		logrus.Error("getting AWS session")
		return err
	}

	logrus.Debug("getting AWS session")

	ssmSvc := ssm.New(sess)

	out, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", o.Env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		logrus.Error("getting bastion instance ID")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*out.Parameter.Value)
	if err != nil {
		logrus.Error("getting bastion instance ID")
		return err
	}

	var to terraformOutput

	err = json.Unmarshal(value, &to)
	if err != nil {
		logrus.Error("getting bastion instance ID")
		return err
	}

	logrus.Debug("getting bastion instance ID")

	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	key, err := getPublicKey(o.PublicKeyFile)
	if err != nil {
		logrus.Error("reading user SSH public key")
		return err
	}

	logrus.Debug("reading user SSH public key")

	command := fmt.Sprintf(
		`grep -qR "%s" /home/ubuntu/.ssh/authorized_keys || echo "%s" >> /home/ubuntu/.ssh/authorized_keys`,
		string(key), string(key),
	)

	_, err = ssmSvc.SendCommand(&ssm.SendCommandInput{
		InstanceIds:  []*string{&to.BastionInstanceID.Value},
		DocumentName: aws.String("AWS-RunShellScript"),
		Comment:      aws.String("Add an SSH public key to authorized_keys"),
		Parameters: map[string][]*string{
			"commands": {&command},
		},
	})

	if err != nil {
		logrus.Error("sending user SSH public key")
		return err
	}

	logrus.Debug("sending user SSH public key")

	return nil
}

func getPublicKey(path string) (string, error) {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return "", err
		}
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("%s does not exist", path)
	}

	var key string
	file, err := os.Open(path)
	if err != nil {
		return key, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return key, nil
}
