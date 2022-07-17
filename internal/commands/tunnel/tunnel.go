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
	"syscall"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/config"
	"github.com/pterm/pterm"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const sshConfig = `# SSH over Session Manager
host i-* mi-*
ServerAliveInterval 180
ProxyCommand sh -c "aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'"

{{range $k :=  .}}LocalForward {{$k}}
{{end}}
`

type Options struct {
	Config         *config.Project
	PrivateKeyFile string
	PublicKeyFile  string
	BastionHostID  string
	ForwardHost    []string
}

func NewTunnelFlags() *Options {
	return &Options{}
}

func NewCmdTunnel() *cobra.Command {
	o := NewTunnelFlags()

	cmd := &cobra.Command{
		Use:              "tunnel",
		Short:            "Tunnel management",
		Long:             "Tunnel management.",
		Args:             cobra.NoArgs,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			err := o.Complete()
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

	cmd.AddCommand(
		NewCmdSSHKey(),
		NewCmdTunnelUp(),
		NewCmdTunnelDown(),
		NewCmdTunnelStatus(),
	)

	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")
	cmd.Flags().StringVar(&o.PublicKeyFile, "ssh-public-key", "", "set ssh key public path")
	cmd.Flags().StringVar(&o.BastionHostID, "bastion-instance-id", "", "set bastion host instance id (i-xxxxxxxxxxxxxxxxx)")
	cmd.Flags().StringSliceVar(&o.ForwardHost, "forward-host", nil, "set forward host for redirect with next format: <remote-host>:<remote-port>. In this case a free local port will be selected automatically.  It's possible to set local manually using <remote-host>:<remote-port>:<local-port>")

	return cmd
}

func (o *Options) Complete() error {
	if err := config.CheckRequirements(config.WithSSMPlugin()); err != nil {
		return err
	}
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("can't configure tunnel: %w", err)
	}

	o.Config = cfg

	isUp, err := checkTunnel(o.Config.EnvDir)
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
		return fmt.Errorf("can't load options for a command: --forward-host parameter requires --bastion-instance-id")
	}

	if len(o.ForwardHost) == 0 && len(o.BastionHostID) != 0 {
		return fmt.Errorf("can't load options for a command: --bastion-instance-id requires --forward-host parameter")
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		if o.Config.Tunnel != nil {
			o.ForwardHost = o.Config.Tunnel.ForwardHost
			o.BastionHostID = o.Config.Tunnel.BastionInstanceID
		}
	}

	if len(o.BastionHostID) == 0 && len(o.ForwardHost) == 0 {
		bastionHostID, forwardHost, err := writeSSHConfigFromSSM(o.Config.Session, o.Config.Env, o.Config.EnvDir)
		if err != nil {
			return err
		}

		o.BastionHostID = bastionHostID
		o.ForwardHost = forwardHost
		pterm.Success.Println("Tunnel forwarding configuration obtained from SSM")
	} else {
		err := writeSSHConfigFromConfig(o.ForwardHost, o.Config.EnvDir)
		if err != nil {
			return err
		}
		pterm.Success.Println("Tunnel forwarding configuration obtained from the config file")
	}

	return nil
}

func (o *Options) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	for _, h := range o.ForwardHost {
		p, _ := strconv.Atoi(strings.Split(h, ":")[2])
		if err := checkPort(p); err != nil {
			return fmt.Errorf("tunnel forwarding config validation failed: %w", err)
		}
	}

	return nil
}

func (o *Options) Run() error {
	logrus.Debugf("public key path: %s", o.PublicKeyFile)

	err := sendSSHPublicKey(o.BastionHostID, getPublicKey(o.PublicKeyFile), o.Config.Session)
	if err != nil {
		return fmt.Errorf("can't run tunnel: %s", err)
	}

	sshConfigPath := fmt.Sprintf("%s/ssh.config", o.Config.EnvDir)

	if err := setAWSCredentials(o.Config.Session); err != nil {
		return fmt.Errorf("can't run tunnel: %w", err)
	}

	args := []string{"-M", "-t", "-S", "bastion.sock", "-fN",
		"-o", "StrictHostKeyChecking=no",
		fmt.Sprintf("ubuntu@%s", o.BastionHostID),
		"-F", sshConfigPath}

	if _, err := os.Stat(o.PrivateKeyFile); !os.IsNotExist(err) {
		args = append(args, "-i", o.PrivateKeyFile)
	}

	c := exec.Command("ssh", args...)

	c.Dir = o.Config.EnvDir

	_, _, code, err := runCommand(c)
	if err != nil {
		return err
	}

	if code != 0 {
		return fmt.Errorf("exit status: %d", code)
	}

	pterm.Success.Println("Tunnel is up! Forwarded ports:")

	var fconfig string
	for _, h := range o.ForwardHost {
		ss := strings.Split(h, ":")
		fconfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", ss[0], ss[1], ss[2])
	}
	pterm.Println(fconfig)

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

	var output terraformOutput

	err = json.Unmarshal(value, &output)
	if err != nil {
		return terraformOutput{}, fmt.Errorf("can't get terraform output: %w", err)
	}

	logrus.Debugf("output: %s", output)

	return output, nil
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

func checkPort(port int) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("can't check address %s: %w", fmt.Sprintf("127.0.0.1:%d", port), err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("port %d is not available. Please make sure there is no other process that is using the port %d", port, port)
	}

	l.Close()

	return nil
}

func writeSSHConfigFromSSM(sess *session.Session, env string, dir string) (string, []string, error) {
	var bastionHostID string
	var forwardHost []string

	to, err := getTerraformOutput(sess, env)
	if err != nil {
		return "", []string{}, fmt.Errorf("can't write SSH config: %w", err)
	}

	sshConfigPath := fmt.Sprintf("%s/ssh.config", dir)

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

func writeSSHConfigFromConfig(forwardHost []string, dir string) error {
	sshConfigPath := fmt.Sprintf("%s/ssh.config", dir)
	f, err := os.Create(sshConfigPath)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}

	tmplData := []string{}
	for k, v := range forwardHost {
		ss := strings.Split(v, ":")
		if len(ss) < 2 || len(ss) > 3 {
			return fmt.Errorf("can't load options for a command: invalid format for forward host (should be host:port:localport)")
		}
		if len(ss) == 2 {
			p, err := getFreePort()
			if err != nil {
				return fmt.Errorf("can't load options for a command: %w", err)
			}
			forwardHost[k] = forwardHost[k] + ":" + strconv.Itoa(p)
			ss = append(ss, strconv.Itoa(p))
		} else if len(ss[2]) == 0 {
			return fmt.Errorf("can't load options for a command: invalid format for forward host (should be host:port:localport)")
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

func checkTunnel(dir string) (bool, error) {
	c := exec.Command(
		"ssh", "-S", "bastion.sock", "-O", "check", "",
	)
	out := &bytes.Buffer{}
	c.Stdout = out
	c.Stderr = out
	c.Dir = dir

	err := c.Run()
	if err == nil {
		sshConfigPath := fmt.Sprintf("%s/ssh.config", dir)
		sshConfig, err := getSSHConfig(sshConfigPath)
		if err != nil {
			return false, fmt.Errorf("can't check tunnel: %w", err)
		}

		pterm.Success.Println("Tunnel is up. Forwarding config:")
		hosts := getHosts(sshConfig)
		var fconfig string
		for _, h := range hosts {
			fconfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", h[2], h[3], h[1])
		}
		pterm.Println(fconfig)

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

func printOutputWithHeader(header string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("%s%s\n", header, scanner.Text())
	}
}

func runCommand(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	var bufout, buferr bytes.Buffer
	outReader2 := io.TeeReader(outReader, &bufout)
	errReader2 := io.TeeReader(errReader, &buferr)

	if err = cmd.Start(); err != nil {
		return
	}

	go printOutputWithHeader("", outReader2)
	go printOutputWithHeader("", errReader2)

	err = cmd.Wait()

	stdout = bufout.String()
	stderr = buferr.String()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}
	return
}
