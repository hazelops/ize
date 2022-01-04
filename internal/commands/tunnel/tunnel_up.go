package tunnel

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

type TunnelUpOptions struct {
	Env            string
	Region         string
	Profile        string
	PrivateKeyFile string
}

func NewTunnelUpFlags() *TunnelUpOptions {
	return &TunnelUpOptions{}
}

func NewCmdTunnelUp() *cobra.Command {
	o := NewTunnelUpFlags()

	cmd := &cobra.Command{
		Use:   "up",
		Short: "open tunnel",
		Long:  "Open tunnel.",
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

	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")

	return cmd
}

func (o *TunnelUpOptions) Complete(cmd *cobra.Command, args []string) error {
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

	return nil
}

func (o *TunnelUpOptions) Validate() error {
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

func (o *TunnelUpOptions) Run() error {
	pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  o.Region,
		Profile: o.Profile,
	})
	if err != nil {
		logrus.Error("getting AWS session")
		fmt.Println(err)
	}

	logrus.Debug("getting AWS session")

	svc := ssm.New(sess)

	resp, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", o.Env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		logrus.Error("getting SSH forward config")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return err
	}

	var config terraformOutput

	err = json.Unmarshal(value, &config)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return err
	}

	logrus.Debugf("output: %s", config)

	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		logrus.Error("getting SSH forward config")
		return err
	}

	hosts := re.FindAllStringSubmatch(
		strings.Join(config.SSHForwardConfig.Value, "\n"),
		-1,
	)

	logrus.Debug("getting SSH forward config")

	logrus.Debugf("hosts: %s", hosts)

	if len(hosts) == 0 {
		logrus.Error("getting SSH forward config")
		return err
	}

	logrus.Debug("port forwarding config is valid")

	localport, err := getFreePort()
	if err != nil {
		logrus.Error("start session")
		return err
	}

	logrus.Debugf("localport: %d", localport)

	input := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartPortForwardingSession"),
		Parameters: map[string][]*string{
			"portNumber":      {aws.String(strconv.Itoa(22))},
			"localPortNumber": {aws.String(strconv.Itoa(localport))},
		},
		Target: &config.BastionInstanceID.Value,
	}

	out, err := svc.StartSession(input)
	if err != nil {
		logrus.Error("start session")
		return err
	}

	pterm.Success.Printfln("Start session")

	err = ssmsession.NewSSMPluginCommand(o.Region).Forward(out, input)
	if err != nil {
		fmt.Println(err)
		logrus.Error("forward server to localhost")
	}

	pterm.Success.Printfln("Forward server to localhost")

	if !filepath.IsAbs(o.PrivateKeyFile) {
		var err error
		o.PrivateKeyFile, err = filepath.Abs(o.PrivateKeyFile)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(o.PrivateKeyFile); err != nil {
		return fmt.Errorf("%s does not exist", o.PrivateKeyFile)
	}

	logrus.Debugf("private key path: %s", o.PrivateKeyFile)

	for _, h := range hosts {
		destinationHost := h[2] + ":" + h[3]

		localPort := h[1]

		tunnel := sshtunnel.NewSSHTunnel(
			"ubuntu@localhost",
			sshtunnel.PrivateKeyFile(o.PrivateKeyFile),
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
		pterm.Info.Printfln("%s:%s => localhost:%s", h[2], h[3], h[1])
		time.Sleep(100 * time.Millisecond)
	}

	pterm.Success.Printfln("Forward destination hosts to localhost")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println()
		done <- true
	}()

	pterm.Info.Println("Press Ctrl-C to to stop the tunnel.")
	<-done
	pterm.Success.Println("Closing connections")

	return err
}

type terraformOutput struct {
	BastionInstanceID struct {
		Value string `json:"value,omitempty"`
	} `json:"bastion_instance_id,omitempty"`
	SSHForwardConfig struct {
		Value []string `json:"value,omitempty"`
	} `json:"ssh_forward_config,omitempty"`
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
