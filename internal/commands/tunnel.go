package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/elliotchance/sshtunnel"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/hazelops/ize/pkg/ssmsession.go"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

type tunnelCmd struct {
	*baseBuilderCmd
}

func (b *commandsBuilder) newTunnelCmd() *tunnelCmd {
	cc := &tunnelCmd{}

	cmd := &cobra.Command{
		Use:   "tunnel",
		Short: "Tunnel management.",
		Long:  "",
		RunE:  nil,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Open tunnel.",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cc.Init()
			if err != nil {
				return err
			}

			pterm.DefaultSection.Printfln("Running SSH Tunnel Up")

			err = cc.BastionSSHTunnelUp()
			if err != nil {
				return err
			}

			return nil
		},
	},
		&cobra.Command{
			Use:   "ssh-key",
			Short: "Send ssh key to remote server.",
			Long:  "",
			RunE: func(cmd *cobra.Command, args []string) error {
				err := cc.Init()
				if err != nil {
					return err
				}

				pterm.DefaultSection.Printfln("Passing SSH Key")

				err = cc.SSHKeyEnsurePresent()
				if err != nil {
					return err
				}

				return nil
			},
		},
	)

	cc.baseBuilderCmd = b.newBuilderBasicCdm(cmd)

	return cc
}

func (c *tunnelCmd) SSHKeyEnsurePresent() error {
	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  c.config.AwsRegion,
		Profile: c.config.AwsProfile,
	})
	if err != nil {
		pterm.Error.Printfln("Getting AWS session")
		return err
	}

	pterm.Success.Printfln("Getting AWS session")

	ssmSvc := ssm.New(sess)

	out, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", c.config.Env)),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		pterm.Error.Printfln("Getting Bastion Instance ID")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*out.Parameter.Value)
	if err != nil {
		pterm.Error.Printfln("Getting Bastion Instance ID")
		return err
	}

	var tOut terraformOutput

	err = json.Unmarshal(value, &tOut)
	if err != nil {
		pterm.Error.Printfln("Getting Bastion Instance ID")
		return err
	}

	pterm.Success.Printfln("Getting Bastion Instance ID")

	key, err := getPublicKey()
	if err != nil {
		pterm.Error.Printfln("Reading user SSH public key")
		return err
	}

	pterm.Success.Printfln("Reading user SSH public key")

	command := fmt.Sprintf(
		`grep -qR "%s" /home/ubuntu/.ssh/authorized_keys || echo "%s" >> /home/ubuntu/.ssh/authorized_keys`,
		string(key), string(key),
	)

	_, err = ssmSvc.SendCommand(&ssm.SendCommandInput{
		InstanceIds:  []*string{&tOut.BastionInstanceID.Value},
		DocumentName: aws.String("AWS-RunShellScript"),
		Comment:      aws.String("Add an SSH public key to authorized_keys"),
		Parameters: map[string][]*string{
			"commands": {&command},
		},
	})

	if err != nil {
		pterm.Error.Printfln("Sending user SSH public Key")
		return err
	}

	pterm.Success.Printfln("Sending user SSH public Key")

	return nil
}

func privateKeyPath() string {
	return os.Getenv("HOME") + "/.ssh/id_rsa"
}

func (c *tunnelCmd) BastionSSHTunnelUp() error {
	err := c.Init()
	if err != nil {
		return err
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  c.config.AwsRegion,
		Profile: c.config.AwsProfile,
	})
	if err != nil {
		pterm.Error.Printfln("Getting AWS session")
		fmt.Println(err)
	}

	pterm.Success.Printfln("Getting AWS session")

	svc := ssm.New(sess)

	resp, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", "dev")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		pterm.Error.Printfln("Getting SSH forward config")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*resp.Parameter.Value)
	if err != nil {
		pterm.Error.Printfln("Getting SSH forward config")
		return err
	}

	var config terraformOutput

	err = json.Unmarshal(value, &config)
	if err != nil {
		pterm.Error.Printfln("Getting SSH forward config")
		return err
	}

	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		pterm.Error.Printfln("Getting SSH forward config")
		return err
	}

	hosts := re.FindAllStringSubmatch(
		strings.Join(config.SSHForwardConfig.Value, "\n"),
		-1,
	)

	if len(hosts) == 0 {
		pterm.Error.Printfln("Getting SSH forward config")
		return err
	}

	pterm.Success.Printfln("Getting SSH forward config")

	input := &ssm.StartSessionInput{
		DocumentName: aws.String("AWS-StartPortForwardingSession"),
		Parameters: map[string][]*string{
			"portNumber":      {aws.String(strconv.Itoa(22))},
			"localPortNumber": {aws.String("30022")},
		},
		Target: aws.String("i-03a1e68e16db39dea"),
	}

	out, err := svc.StartSession(input)
	if err != nil {
		pterm.Error.Printfln("Start session")
		return err
	}

	pterm.Success.Printfln("Start session")

	err = ssmsession.NewSSMPluginCommand(c.config.AwsRegion).Forward(out, input)
	if err != nil {
		fmt.Println(err)
		pterm.Error.Printfln("Forward server to localhost")
	}

	pterm.Success.Printfln("Forward server to localhost")

	for _, h := range hosts {
		destinationHost := h[2] + ":" + h[3]

		localPort := h[1]

		tunnel := sshtunnel.NewSSHTunnel(
			"ubuntu@localhost",
			sshtunnel.PrivateKeyFile(privateKeyPath()),
			destinationHost,
			localPort,
		)

		tunnel.Server.Port = 30022

		go func() {
			if err := tunnel.Start(); err != nil {
				pterm.Error.Printfln("Forward destination hosts to localhost")
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

	pterm.Info.Println("Press Ctrl-C to close the connections.")
	<-done
	pterm.Success.Println("Сlosing connections")

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

func getPublicKey() (string, error) {
	var key string
	home, _ := os.UserHomeDir()
	file, err := os.Open(fmt.Sprintf("%s/.ssh/id_rsa.pub", home))
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
