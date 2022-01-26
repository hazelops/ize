package tunnel

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelOptions struct {
	Config         *config.Config
	PrivateKeyFile string
	PublicKeyFile  string
}

func NewTunnelFlags() *TunnelOptions {
	return &TunnelOptions{}
}

func NewCmdTunnel() *cobra.Command {
	o := NewTunnelFlags()

	cmd := &cobra.Command{
		Use:              "tunnel",
		Short:            "tunnel management",
		Long:             "Tunnel management.",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run(cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(
		NewCmdSSHKey(),
		NewCmdTunnelUp(),
		NewCmdTunnelDown(),
		NewCmdTunnelStatus(),
	)

	return cmd
}

func (o *TunnelOptions) Complete(cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig()
	if err != nil {
		return err
	}

	o.Config = cfg

	if o.PrivateKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PrivateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", home)
	}

	if o.PublicKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PublicKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	return nil
}

func (o *TunnelOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	return nil
}

func (o *TunnelOptions) Run(cmd *cobra.Command) error {
	c := exec.Command(
		"ssh", "-S", "bastion.sock", "-O", "check", "",
	)
	out := &bytes.Buffer{}
	c.Stdout = out
	c.Stderr = out

	err := c.Run()
	if err == nil {
		sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))
		sshConfig, err := getSSHConfig(sshConfigPath)
		if err != nil {
			return fmt.Errorf("can't run tunnel: %w", err)
		}
		hosts := getHosts(sshConfig)
		pterm.Info.Printfln("tunnel is already up. Forwarded ports:")
		for _, h := range hosts {
			pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
		}

		return nil
	}

	pterm.Success.Println("checking for an existing tunnel")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Config.AwsRegion,
		Profile: o.Config.AwsProfile,
	})
	if err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}

	logrus.Debug("getting AWS session")

	to, err := getTerraformOutput(sess, o.Config.Env)
	if err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}

	logrus.Debug("getting bastion instance ID")

	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	err = sendSSHPublicKey(to.BastionInstanceID.Value, getPublicKey(o.PublicKeyFile), sess)
	if err != nil {
		return fmt.Errorf("can't run tunnel: %s", err)
	}

	pterm.Success.Println("sending user SSH public key")

	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))

	f, err := os.Create(sshConfigPath)
	if err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}

	sshConfig := strings.Join(to.SSHForwardConfig.Value, "\n")
	_, err = io.WriteString(f, sshConfig)
	if err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}
	if err = f.Close(); err != nil {
		return err
	}

	hosts := getHosts(sshConfig)
	if len(hosts) == 0 {
		return fmt.Errorf("can't tunnel: forward config is not valid")
	}
	logrus.Debugf("hosts: %s", hosts)

	c = exec.Command(
		"ssh", "-M", "-S", "bastion.sock", "-fNT",
		fmt.Sprintf("ubuntu@%s", to.BastionInstanceID.Value),
		"-F", sshConfigPath,
	)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}

	pterm.Success.Printfln("tunnel is up. Forwarding config:")
	for _, h := range hosts {
		pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
	}

	return nil
}

func getTerraformOutput(sess *session.Session, env string) (terraformOutput, error) {
	resp, err := ssm.New(sess).GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return terraformOutput{}, fmt.Errorf("can't get terraform output: %w", err)
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		return terraformOutput{}, fmt.Errorf("can't get terraform output: %w", err)
	}

	var config terraformOutput

	err = json.Unmarshal(value, &config)
	if err != nil {
		return terraformOutput{}, fmt.Errorf("can't get terraform output: %w", err)
	}

	logrus.Debugf("output: %s", config)

	return config, nil
}

type terraformOutput struct {
	BastionInstanceID struct {
		Value string `json:"value,omitempty"`
	} `json:"bastion_instance_id,omitempty"`
	SSHForwardConfig struct {
		Value []string `json:"value,omitempty"`
	} `json:"ssh_forward_config,omitempty"`
}

func sendSSHPublicKey(bastionID string, key string, sess *session.Session) error {
	command := fmt.Sprintf(
		`grep -qR "%s" /home/ubuntu/.ssh/authorized_keys || echo "%s" >> /home/ubuntu/.ssh/authorized_keys`,
		string(key), string(key),
	)

	_, err := ssm.New(sess).SendCommand(&ssm.SendCommandInput{
		InstanceIds:  []*string{&bastionID},
		DocumentName: aws.String("AWS-RunShellScript"),
		Comment:      aws.String("Add an SSH public key to authorized_keys"),
		Parameters: map[string][]*string{
			"commands": {&command},
		},
	})

	if err != nil {
		return fmt.Errorf("can't send SSH public key: %w", err)
	}

	return nil
}

func getPublicKey(path string) string {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(path); err != nil {
		log.Fatalf("%s does not exist", path)
	}

	var key string
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		key = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return key
}

func getHosts(sshconfig string) [][]string {
	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		log.Fatal(fmt.Errorf("can't get forwaed config: %w", err))
	}

	hosts := re.FindAllStringSubmatch(
		sshconfig,
		-1,
	)

	return hosts
}

func getSSHConfig(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("can't get ssh config: %w", err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("can't get ssh config: %w", err)
	}

	return string(b), nil
}
