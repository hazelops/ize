package tunnel

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/elliotchance/sshtunnel"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/ssmsession"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TunnelOptions struct {
	Env            string
	Region         string
	Profile        string
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

func (o *TunnelOptions) Run(cmd *cobra.Command) error {
	pterm.DefaultSection.Printfln("Sending SSH public key")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
	})
	if err != nil {
		return fmt.Errorf("can't get AWS session: %w", err)
	}

	logrus.Debug("getting AWS session")

	to, err := getTerraformOutput(sess, o.Env)
	if err != nil {
		return fmt.Errorf("can't get forward config: %w", err)
	}

	logrus.Debug("getting bastion instance ID")

	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	err = sendSSHPublicKey(to.BastionInstanceID.Value, getPublicKey(o.PublicKeyFile), sess)
	if err != nil {
		logrus.Error("sending user SSH public key")
		return err
	}

	logrus.Debug("sending user SSH public key")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

	hosts := getForwardConfig(to)

	logrus.Debug("getting SSH forward config")

	logrus.Debugf("hosts: %s", hosts)

	if len(hosts) == 0 {
		return fmt.Errorf("can't tunnel up: forward config is not valid")
	}

	logrus.Debug("forwarding config is valid")

	_ = killDaemon(daemonContext(ctx))

	p, err := daemonContext(ctx).Reborn()
	if err != nil {
		return fmt.Errorf("restarted tunnel process: %w", err)
	}
	if p != nil {
		logrus.Debugf("tunnel is up(pid: %d)", p.Pid)
		pterm.Success.Printfln("tunnel is up")
		pterm.Info.Printfln("forward config:")
		for _, h := range hosts {
			pterm.Info.Printfln("%s:%s ➡ localhost:%s", h[2], h[3], h[1])
		}
		return nil
	}
	defer daemonContext(ctx).Release()

	localport, sessionID, err := startPortForwardSession(to, o.Region, sess)
	if err != nil {
		return fmt.Errorf("can't tunnel up: %w", err)
	}

	logrus.Debugf("private key path: %s", o.PrivateKeyFile)

	for _, h := range hosts {
		destinationHost := h[2] + ":" + h[3]

		localPort := h[1]

		tunnel := sshtunnel.NewSSHTunnel(
			"ubuntu@localhost",
			sshtunnel.PrivateKeyFile(getPrivateKey(o.PrivateKeyFile)),
			destinationHost,
			localPort,
		)

		tunnel.Server.Port = localport

		go func() {
			if err := tunnel.Start(); err != nil {
				pterm.Error.Printfln("Forward destination host to localhost")
				os.Exit(1)
			}
		}()
		pterm.Info.Printfln("%s ➡ localhost:%s", destinationHost, localPort)
		time.Sleep(100 * time.Millisecond)
	}

	for {
		select {
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGTERM:
				ssm.New(sess).TerminateSession(&ssm.TerminateSessionInput{
					SessionId: &sessionID,
				})
				cancel()
				return nil
			}
		}
	}
}

func getTerraformOutput(sess *session.Session, env string) (terraformOutput, error) {
	resp, err := ssm.New(sess).GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		logrus.Error("getting SSH forward config")
		return terraformOutput{}, err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return terraformOutput{}, err
	}

	var config terraformOutput

	err = json.Unmarshal(value, &config)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return terraformOutput{}, err
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

func startPortForwardSession(to terraformOutput, region string, sess *session.Session) (int, string, error) {
	localport, err := getFreePort()
	if err != nil {
		return 0, "", fmt.Errorf("can't start port forward session: %w", err)
	}

	input := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartPortForwardingSession"),
		Parameters: map[string][]*string{
			"portNumber":      {aws.String(strconv.Itoa(22))},
			"localPortNumber": {aws.String(strconv.Itoa(localport))},
		},
		Target: &to.BastionInstanceID.Value,
	}

	out, err := ssm.New(sess).StartSession(input)
	if err != nil {
		return 0, "", fmt.Errorf("can't start port forward session: %w", err)
	}

	err = ssmsession.NewSSMPluginCommand(region).Forward(out, input)
	if err != nil {
		return 0, "", fmt.Errorf("can't start port forward session: %w", err)
	}

	return localport, *out.SessionId, nil
}

func getPrivateKey(path string) string {
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

	return path
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

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func getForwardConfig(to terraformOutput) [][]string {
	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		log.Fatal(fmt.Errorf("can't get forwaed config: %w", err))
	}

	hosts := re.FindAllStringSubmatch(
		strings.Join(to.SSHForwardConfig.Value, "\n"),
		-1,
	)

	return hosts
}
