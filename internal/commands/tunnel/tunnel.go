package tunnel

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const sshConfig = `# SSH over Session Manager
host i-* mi-*
ServerAliveInterval 180
ProxyCommand sh -c "aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'"

{{range $k :=  .}}LocalForward {{$k}}
{{end}}
`

type TunnelOptions struct {
	Config         *config.Config
	PrivateKeyFile string
	PublicKeyFile  string
	BastionHostID  string
	ForwardHost    []string
}

func NewTunnelFlags() *TunnelOptions {
	return &TunnelOptions{}
}

func NewCmdTunnel(ui terminal.UI) *cobra.Command {
	o := NewTunnelFlags()

	cmd := &cobra.Command{
		Use:              "tunnel",
		Short:            "tunnel management",
		Long:             "Tunnel management.",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			sg := ui.StepGroup()
			defer sg.Wait()

			err := o.Complete(ui, sg, cmd, args)
			if err != nil {
				return err
			}

			err = o.Validate()
			if err != nil {
				return err
			}

			err = o.Run(ui, sg, cmd)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(
		NewCmdSSHKey(ui),
		NewCmdTunnelUp(ui),
		NewCmdTunnelDown(ui),
		NewCmdTunnelStatus(ui),
	)

	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")
	cmd.Flags().StringVar(&o.PublicKeyFile, "ssh-public-key", "", "set ssh key public path")
	cmd.Flags().StringVar(&o.BastionHostID, "bastion-instance-id", "", "set bastion host instance id (i-xxxxxxxxxxxxxxxxx)")
	cmd.Flags().StringSliceVar(&o.ForwardHost, "forward-host", nil, "set forward host for redirect with next format: <remote-host>:<remote-port>. In this case a free local port will be selected automatically.  It's possible to set local manually using <remote-host>:<remote-port>:<local-port>")

	return cmd
}

func (o *TunnelOptions) Complete(ui terminal.UI, sg terminal.StepGroup, cmd *cobra.Command, args []string) error {
	cfg, err := config.InitializeConfig(config.WithSSMPlugin())
	if err != nil {
		return fmt.Errorf("can't complete options: %w", err)
	}

	o.Config = cfg

	isUp, err := checkTunnel(ui, sg)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}
	if isUp {
		os.Exit(0)
	}

	if o.PrivateKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PrivateKeyFile = fmt.Sprintf("%s/.ssh/id_rsa", home)
	}

	if o.PublicKeyFile == "" {
		home, _ := os.UserHomeDir()
		o.PublicKeyFile = fmt.Sprintf("%s/.ssh/id_rsa.pub", home)
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) != 0 {
		return fmt.Errorf("cat't complete options: --forward-host parameter requires --bastion-instance-id\n")
	}

	if len(o.ForwardHost) == 0 && len(o.BastionHostID) != 0 {
		return fmt.Errorf("cat't complete options: --bastion-instance-id requires --forward-host parameter\n")
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		viper.UnmarshalKey("infra.tunnel.bastion_instance_id", &o.BastionHostID)
		viper.UnmarshalKey("infra.tunnel.forward_host", &o.ForwardHost)
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		s := sg.Add("writing SSH config from SSM...")
		bastionHostID, forwardHost, err := writeSSHConfigFromSSM(o.Config.Session, o.Config.Env)
		if err != nil {
			return err
		}

		o.BastionHostID = bastionHostID
		o.ForwardHost = forwardHost
		s.Done()
	} else {
		s := sg.Add("writing SSH config from config...", terminal.WithInfoStyle())
		err := writeSSHConfigFromConfig(o.ForwardHost)
		if err != nil {
			return err
		}
		s.Done()
	}

	return nil
}

func (o *TunnelOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified\n")
	}

	return nil
}

func (o *TunnelOptions) Run(ui terminal.UI, sg terminal.StepGroup, cmd *cobra.Command) error {
	s := sg.Add("sending user SSH public key...")
	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	err := sendSSHPublicKey(o.BastionHostID, getPublicKey(o.PublicKeyFile), o.Config.Session)
	if err != nil {
		return fmt.Errorf("can't run tunnel: %s", err)
	}

	s.Done()
	s = sg.Add("upping tunnel...")

	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))

	if err := setAWSCredentials(o.Config.Session); err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}

	c := exec.Command(
		"ssh", "-M", "-S", "bastion.sock", "-fNT",
		fmt.Sprintf("ubuntu@%s", o.BastionHostID),
		"-F", sshConfigPath,
		"-i", getPrivateKey(o.PrivateKeyFile),
	)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = viper.GetString("ENV_DIR")
	if err := c.Run(); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	s.Done()
	ui.Output("tunnel is up! Forwarded ports:", terminal.WithSuccessStyle())

	var fconfig string
	for _, h := range o.ForwardHost {
		ss := strings.Split(h, ":")
		fconfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", ss[0], ss[1], ss[2])
	}
	ui.Output(fconfig)

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
			logrus.Fatal(err)
		}
	}

	if _, err := os.Stat(path); err != nil {
		logrus.Fatalf("%s does not exist", path)
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
		log.Fatal(fmt.Errorf("can't get forward config: %w", err))
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

func getPrivateKey(path string) string {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	f, err := os.Stat(path)
	if err != nil {
		logrus.Fatalf("%s does not exist", path)
	}

	if f.IsDir() {
		logrus.Fatalf("%s is a directory", path)
	}

	return path
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

func writeSSHConfigFromSSM(sess *session.Session, env string) (string, []string, error) {
	var bastionHostID string
	var forwardHost []string

	to, err := getTerraformOutput(sess, env)
	if err != nil {
		return "", []string{}, fmt.Errorf("can't write SSH config: %w", err)
	}

	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))

	f, err := os.Create(sshConfigPath)
	if err != nil {
		return "", []string{}, fmt.Errorf("can't write SSH config: %w", err)
	}

	sshConfig := strings.Join(to.SSHForwardConfig.Value, "\n")
	_, err = io.WriteString(f, sshConfig)
	if err != nil {
		return "", []string{}, fmt.Errorf("can't write SSH config: %w", err)
	}
	if err = f.Close(); err != nil {
		return "", []string{}, fmt.Errorf("can't write SSH config: %w", err)
	}

	hosts := getHosts(sshConfig)
	if len(hosts) == 0 {
		return "", []string{}, fmt.Errorf("can't write SSH config: forwarding config is not valid")
	}

	bastionHostID = to.BastionInstanceID.Value

	for _, h := range hosts {
		forwardHost = append(forwardHost, fmt.Sprintf("%s:%s:%s", h[2], h[3], h[1]))
	}

	return bastionHostID, forwardHost, nil
}

func writeSSHConfigFromConfig(forwardHost []string) error {
	sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))
	f, err := os.Create(sshConfigPath)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	tmplData := []string{}
	for k, v := range forwardHost {
		ss := strings.Split(v, ":")
		if len(ss) < 2 || len(ss) > 3 {
			return fmt.Errorf("can't complete options: invalid format for forward host (should be host:port:localport)\n")
		}
		if len(ss) == 2 {
			p, err := getFreePort()
			if err != nil {
				return fmt.Errorf("can't complete options: %w", err)
			}
			forwardHost[k] = forwardHost[k] + ":" + strconv.Itoa(p)
			ss = append(ss, strconv.Itoa(p))
		} else if len(ss[2]) == 0 {
			return fmt.Errorf("can't complete options: invalid format for forward host (should be host:port:localport)\n")
		}
		tmplData = append(tmplData, fmt.Sprintf("%s %s:%s", ss[2], ss[0], ss[1]))
	}
	t := template.New("sshConfig")
	t, err = t.Parse(sshConfig)
	if err != nil {
		return err
	}
	err = t.Execute(f, tmplData)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	return nil
}

func checkTunnel(ui terminal.UI, sg terminal.StepGroup) (bool, error) {
	s := sg.Add("checking for an existing tunnel...")

	c := exec.Command(
		"ssh", "-S", "bastion.sock", "-O", "check", "",
	)
	out := &bytes.Buffer{}
	c.Stdout = out
	c.Stderr = out
	c.Dir = viper.GetString("ENV_DIR")

	err := c.Run()
	s.Done()
	if err == nil {
		sshConfigPath := fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR"))
		sshConfig, err := getSSHConfig(sshConfigPath)
		if err != nil {
			return false, fmt.Errorf("can't check tunnel: %w", err)
		}

		ui.Output("tunnel is up. Forwarding config:", terminal.WithSuccessStyle())

		hosts := getHosts(sshConfig)
		var fconfig string
		for _, h := range hosts {
			fconfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", h[2], h[3], h[1])
		}
		ui.Output(fconfig)

		return true, nil
	}

	return false, nil
}

func setAWSCredentials(sess *session.Session) error {
	v, err := sess.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	os.Setenv("AWS_SECRET_ACCESS_KEY", v.SecretAccessKey)
	os.Setenv("AWS_ACCESS_KEY_ID", v.AccessKeyID)
	os.Setenv("AWS_SESSION_TOKEN", v.SessionToken)

	return nil
}
