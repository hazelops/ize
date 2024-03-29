package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/internal/requirements"
	"github.com/hazelops/ize/pkg/term"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const sshConfig = `# SSH over Session Manager
host i-* mi-*
ServerAliveInterval 180
ProxyCommand sh -c "aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'"

{{range $k :=  .}}LocalForward {{$k}}
{{end}}
`

var explainTunnelUpTmpl = `
# Set variables
SSH_CONFIG={{.EnvDir}}/ssh.config
SSH_PUBLIC_KEY=$(cat ~/.ssh/id_rsa.pub)

# Get bastion instance id
BASTION_INSTANCE_ID=$(aws ssm get-parameter --name "/{{.Env}}/terraform-output" --with-decryption | jq -r '.Parameter.Value' | base64 -d | jq -r '.bastion_instance_id.value'

# Get ssh config
aws ssm get-parameter --name "/{{.Env}}/terraform-output" --with-decryption | jq -r '.Parameter.Value' | base64 -d | jq -r '.ssh_forward_config.value[]' > $SSH_CONFIG

# Send ssh public key to instance
aws ssm send-command --instance-ids $BASTION_INSTANCE_ID --document-name AWS-RunShellScript --comment 'Add an SSH public key to authorized_keys' --parameters '{"commands": ["grep -qR \"$(SSH_PUBLIC_KEY)\" /home/ubuntu/.ssh/authorized_keys || echo \"$(SSH_PUBLIC_KEY)\" >> /home/ubuntu/.ssh/authorized_keys"]}' 1> /dev/null)

# Change to the dir and up tunnel
(cd {{.EnvDir}} && $(aws ssm get-parameter --name "/{{.Env}}/terraform-output" --with-decryption | jq -r '.Parameter.Value' | base64 -d | jq -r '.cmd.value.tunnel.up') -F $SSH_CONFIG)
`

type TunnelUpOptions struct {
	Config                *config.Project
	PrivateKeyFile        string
	PublicKeyFile         string
	BastionHostID         string
	ForwardHost           []string
	StrictHostKeyChecking bool
	Metadata              bool
	Explain               bool
}

func NewTunnelUpFlags(project *config.Project) *TunnelUpOptions {
	return &TunnelUpOptions{
		Config: project,
	}
}

func NewCmdTunnelUp(project *config.Project) *cobra.Command {
	o := NewTunnelUpFlags(project)

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Open tunnel with sending ssh key",
		Long:  "Open tunnel with sending ssh key to remote server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			if o.Explain {
				err := o.Config.Generate(explainTunnelUpTmpl, nil)
				if err != nil {
					return err
				}

				return nil
			}

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

	cmd.Flags().StringVar(&o.BastionHostID, "bastion-instance-id", "", "set bastion host instance id (i-xxxxxxxxxxxxxxxxx)")
	cmd.Flags().StringSliceVar(&o.ForwardHost, "forward-host", nil, "set forward hosts for redirect with next format: <remote-host>:<remote-port>, <remote-host>:<remote-port>, <remote-host>:<remote-port>. In this case a free local port will be selected automatically.  It's possible to set local manually using <remote-host>:<remote-port>:<local-port>")
	cmd.Flags().StringVar(&o.PublicKeyFile, "ssh-public-key", "", "set ssh key public path")
	cmd.Flags().StringVar(&o.PrivateKeyFile, "ssh-private-key", "", "set ssh key private path")
	cmd.PersistentFlags().BoolVar(&o.StrictHostKeyChecking, "strict-host-key-checking", false, "set strict host key checking")
	cmd.PersistentFlags().BoolVar(&o.Metadata, "use-ec2-metadata", false, "send ssh key to EC2 metadata (work only for Ubuntu versions > 20.0)")
	cmd.Flags().BoolVar(&o.Explain, "explain", false, "bash alternative shown")

	return cmd
}

func (o *TunnelUpOptions) Complete() error {
	if err := requirements.CheckRequirements(requirements.WithSSMPlugin()); err != nil {
		return err
	}

	isUp, err := checkTunnel(o.Config.EnvDir)
	if err != nil {
		return fmt.Errorf("can't run tunnel up: %w", err)
	}
	if isUp {
		os.Exit(0)
	}

	if o.PrivateKeyFile == "" && o.Config.Tunnel != nil {
		o.PrivateKeyFile = o.Config.Tunnel.SSHPrivateKey
	}

	if o.PublicKeyFile == "" && o.Config.Tunnel != nil {
		o.PublicKeyFile = o.Config.Tunnel.SSHPublicKey
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
		wr := new(SSMWrapper)
		wr.Api = ssm.New(o.Config.Session)
		bastionHostID, forwardHost, err := writeSSHConfigFromSSM(wr, o.Config.Env, o.Config.EnvDir)
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

func (o *TunnelUpOptions) Validate() error {
	if len(o.Config.Env) == 0 {
		return fmt.Errorf("env must be specified")
	}

	for _, h := range o.ForwardHost {
		p, _ := strconv.Atoi(strings.Split(h, ":")[2])
		if err := checkPort(p, o.Config.EnvDir); err != nil {
			return fmt.Errorf("tunnel forwarding config validation failed: %w", err)
		}
	}

	return nil
}

func (o *TunnelUpOptions) Run() error {
	logrus.Debugf("public key path: %s", o.PublicKeyFile)
	logrus.Debugf("private key path: %s", o.PrivateKeyFile)

	err := o.checkOsVersion()
	if err != nil {
		return err
	}

	pk, err := getPublicKey(o.PublicKeyFile)
	if err != nil {
		return fmt.Errorf("can't get public key: %s", err)
	}

	logrus.Debugf("public key:\n%s", pk)

	if o.Metadata {
		err = sendSSHPublicKey(o.BastionHostID, pk, o.Config.Session)
		if err != nil {
			return fmt.Errorf("can't run tunnel: %s", err)
		}
	} else {
		err = sendSSHPublicKeyLegacy(o.BastionHostID, pk, o.Config.Session)
		if err != nil {
			return fmt.Errorf("can't run tunnel: %s", err)
		}
	}

	forwardConfig, err := o.upTunnel()
	if err != nil {
		return err
	}

	pterm.Success.Println("Tunnel is up! Forwarded ports:")
	pterm.Println(forwardConfig)

	return nil
}

func (o *TunnelUpOptions) checkOsVersion() error {
	diio, err := o.Config.AWSClient.SSMClient.DescribeInstanceInformation(&ssm.DescribeInstanceInformationInput{
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("InstanceIds"),
				Values: aws.StringSlice([]string{o.BastionHostID}),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("can't get instance '%s' information: %s", o.BastionHostID, err)
	}

	if len(diio.InstanceInformationList) == 0 {
		return fmt.Errorf("can't get instance '%s' information", o.BastionHostID)
	}

	osName := *diio.InstanceInformationList[0].PlatformName
	osVersion := *diio.InstanceInformationList[0].PlatformVersion

	switch osName {
	case "Ubuntu":
		if semver.MustParse(osVersion).LessThan(semver.MustParse("20.04")) {
			pterm.Warning.Printfln("Your bastion host AMI is Ubuntu %s, Instance Connect is not installed by default on that version of OS. Please upgrade to at least v20.04", osVersion)
		}
	case "Amazon Linux AMI":
		if semver.MustParse(osVersion).LessThan(semver.MustParse("2.0.20190618")) {
			pterm.Warning.Printfln("Your bastion host AMI is Amazon Linux AMI %s, Instance Connect is not installed by default on that version of OS. Please upgrade your AMI to at least 2.0.20190618", osVersion)
		}
	}

	return nil
}

func (o *TunnelUpOptions) upTunnel() (string, error) {
	sshConfigPath := fmt.Sprintf("%s/ssh.config", o.Config.EnvDir)
	logrus.Debugf("ssh config path: %s", sshConfigPath)

	if err := setAWSCredentials(o.Config.Session); err != nil {
		return "", fmt.Errorf("can't run tunnel: %w", err)
	}

	args := o.getSSHCommandArgs(sshConfigPath)

	err := o.runSSH(args)
	if err != nil {
		return "", err
	}

	var forwardConfig string
	for _, h := range o.ForwardHost {
		ss := strings.Split(h, ":")
		forwardConfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", ss[0], ss[1], ss[2])
	}
	return forwardConfig, nil
}

func (o *TunnelUpOptions) runSSH(args []string) error {
	c := exec.Command("ssh", args...)

	c.Dir = o.Config.EnvDir
	os.Setenv("AWS_REGION", o.Config.AwsRegion)

	runner := term.New(term.WithStdin(os.Stdin))
	_, _, code, err := runner.Run(c)
	if err != nil {
		return err
	}

	if code != 0 {
		return fmt.Errorf("exit status: %d", code)
	}
	return nil
}

func (o *TunnelUpOptions) getSSHCommandArgs(sshConfigPath string) []string {
	args := []string{"-M", "-t", "-S", "bastion.sock", "-fN"}
	if !o.StrictHostKeyChecking {
		args = append(args, "-o", "StrictHostKeyChecking=no")
	}
	args = append(args, fmt.Sprintf("ubuntu@%s", o.BastionHostID))
	args = append(args, "-F", sshConfigPath)

	if _, err := os.Stat(o.PrivateKeyFile); !os.IsNotExist(err) {
		args = append(args, "-i", o.PrivateKeyFile)
	}

	if o.Config.LogLevel == "debug" {
		args = append(args, "-vvv")
	}

	return args
}

type SSMWrapper struct {
	Api ssmiface.SSMAPI
}

func getTerraformOutput(wr *SSMWrapper, env string) (terraformOutput, error) {
	resp, err := wr.Api.GetParameter(&ssm.GetParameterInput{
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

	logrus.Debugf("decoded terrafrom output: \n%s", value)

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
	_, err := ec2instanceconnect.New(sess).SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
		InstanceId:     aws.String(bastionID),
		InstanceOSUser: aws.String("ubuntu"),
		SSHPublicKey:   aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func sendSSHPublicKeyLegacy(bastionID string, key string, sess *session.Session) error {
	// This command is executed in the bastion host and it checks if our public key is present. If it's not it uploads it to _authorized_keys file.
	command := fmt.Sprintf(
		`grep -qR "%s" /home/ubuntu/.ssh/authorized_keys || echo "%s" >> /home/ubuntu/.ssh/authorized_keys`,
		strings.TrimSpace(key), strings.TrimSpace(key),
	)

	logrus.Debugf("send command: \n%s", command)

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

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	_, _, _, _, err = ssh.ParseAuthorizedKey(f)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

func getHosts(config string) [][]string {
	// This regexp reads ssh.conf configuration, so we can display it nicely in the UI
	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		log.Fatal(fmt.Errorf("can't get forward config: %w", err))
	}

	hosts := re.FindAllStringSubmatch(
		config,
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
	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}

func checkPort(port int, dir string) error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("can't check address %s: %w", fmt.Sprintf("127.0.0.1:%d", port), err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		command := fmt.Sprintf("lsof -i tcp:%d | grep LISTEN | awk '{print $1, $2}'", port)
		stdout, stderr, code, err := term.New(term.WithStdout(io.Discard), term.WithStderr(io.Discard)).Run(exec.Command("bash", "-c", command))
		if err != nil {
			return fmt.Errorf("can't run command '%s': %w", command, err)
		}
		if code == 0 {
			stdout = strings.TrimSpace(stdout)
			processName := strings.Split(stdout, " ")[0]
			processPid, err := strconv.Atoi(strings.Split(stdout, " ")[1])
			if err != nil {
				return fmt.Errorf("can't get pid: %w", err)
			}
			pterm.Info.Printfln("Can't start tunnel on port %d. It seems like it's take by a process '%s'.", port, processName)
			proc, err := os.FindProcess(processPid)
			if err != nil {
				return fmt.Errorf("can't find process: %w", err)
			}

			_, err = os.Stat(filepath.Join(dir, "bastion.sock"))
			if processName == "ssh" && os.IsNotExist(err) {
				return fmt.Errorf("it could be another ize tunnel, but we can't find a socket. Something went wrong. We suggest terminating it and starting it up again")
			}
			isContinue := false
			if terminal.IsTerminal(int(os.Stdout.Fd())) {
				isContinue, err = pterm.DefaultInteractiveConfirm.WithDefaultText("Would you like to terminate it?").Show()
				if err != nil {
					return err
				}
			} else {
				isContinue = true
			}

			if !isContinue {
				return fmt.Errorf("destroying was canceled")
			}
			err = proc.Kill()
			if err != nil {
				return fmt.Errorf("can't kill process: %w", err)
			}

			pterm.Info.Printfln("Process '%s' (pid %d) was killed", processName, processPid)

			return nil
		}
		return fmt.Errorf("error during run command: %s (exit code: %d, stderr: %s)", command, code, stderr)
	}

	err = l.Close()
	if err != nil {
		return err
	}

	return nil
}

func writeSSHConfigFromSSM(wr *SSMWrapper, env string, dir string) (string, []string, error) {
	var bastionHostID string
	var forwardHost []string

	to, err := getTerraformOutput(wr, env)
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
		errMsg := "can't write SSH config: forwarding config is not valid"
		if logrus.GetLevel() == logrus.DebugLevel {
			errMsg += fmt.Sprintf(". Config in SSM: \n%s", sshConfig)
		}
		return "", []string{}, fmt.Errorf(errMsg)
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

	var tmplData []string
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
	pathToSocket := filepath.Join(dir, "bastion.sock")
	if _, err := os.Stat(pathToSocket); !os.IsNotExist(err) {
		pterm.Info.Printfln("A socket file from another tunnel has been detected: %s", pathToSocket)
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
			var forwardConfig string
			for _, h := range hosts {
				forwardConfig += fmt.Sprintf("%s:%s ➡ localhost:%s\n", h[2], h[3], h[1])
			}
			pterm.Println(forwardConfig)

			return true, nil
		} else {
			pterm.Warning.Println("Tunnel socket file seems to be not useable. We have deleted it")
			err := os.Remove(pathToSocket)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}

	return false, nil
}

func setAWSCredentials(sess *session.Session) error {
	v, err := sess.Config.Credentials.Get()
	if err != nil {
		return fmt.Errorf("can't set AWS credentials: %w", err)
	}

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", v.SecretAccessKey)
	if err != nil {
		return err
	}
	err = os.Setenv("AWS_ACCESS_KEY_ID", v.AccessKeyID)
	if err != nil {
		return err
	}
	err = os.Setenv("AWS_SESSION_TOKEN", v.SessionToken)
	if err != nil {
		return err
	}

	return nil
}
