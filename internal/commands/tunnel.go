package commands

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

			err = cc.BastionSHHTunnelUp()
			if err != nil {
				return err
			}

			return nil
		},
	},
		&cobra.Command{
			Use:   "down",
			Short: "Close tunnel.",
			Long:  "",
			RunE: func(cmd *cobra.Command, args []string) error {
				err := cc.Init()
				if err != nil {
					return err
				}

				pterm.DefaultSection.Printfln("Running SSH Tunnel Down")

				err = cc.BastionSHHTunnelDown()
				if err != nil {
					return err
				}

				return nil
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show status tunnel.",
			Long:  "",
			RunE: func(cmd *cobra.Command, args []string) error {
				err := cc.Init()
				if err != nil {
					return err
				}

				pterm.DefaultSection.Printfln("Running SSH Tunnel Status")

				err = cc.BastionSHHTunnelStatus()
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
		&cobra.Command{
			Use:   "config",
			Short: "Create ssh config.",
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

				pterm.DefaultSection.Printfln("Getting SSH Tunnel config")

				err = cc.SSHTunnelConfigCreate()
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

func (c *tunnelCmd) SSHTunnelConfigCreate() error {
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
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", "dev")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		pterm.Error.Printfln("Getting SSH Tunnel config")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*out.Parameter.Value)
	if err != nil {
		pterm.Error.Printfln("Getting SSH Tunnel config")
		return err
	}

	var tOut terraformOutput

	err = json.Unmarshal(value, &tOut)
	if err != nil {
		pterm.Error.Printfln("Getting SSH Tunnel config")
		return err
	}

	pterm.Success.Printfln("Getting SSH Tunnel config")

	f, err := os.Create(fmt.Sprintf("%s/%s", viper.GetString("ENV_DIR"), "ssh.config"))
	if err != nil {
		pterm.Error.Printfln("Writing SSH Tunnel config to file")
		return err
	}

	re, err := regexp.Compile(`LocalForward\s(?P<localPort>\d+)\s(?P<remoteHost>.+):(?P<remotePort>\d+)`)
	if err != nil {
		return err
	}

	var config string

	defer f.Close()
	for _, v := range tOut.SSHForwardConfig.Value {
		config += v + "\n"
		_, err = fmt.Fprintln(f, v)
		if err != nil {
			pterm.Error.Printfln("Writing SSH Tunnel config to file")
			return err
		}
	}

	pterm.Success.Printfln("Writing SSH Tunnel config to file")

	res := re.FindAllStringSubmatch(config, -1)

	ports := "SSH Tunnel Available:\n"

	for _, v := range res {
		ports += fmt.Sprintf("%s:%s => localhost:%s\n", v[2], v[3], v[1])
	}

	pterm.Info.Print(ports)

	return nil
}

func (c *tunnelCmd) BastionSHHTunnelUp() error {
	var err error
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
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", "dev")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		pterm.Error.Printfln("Getting tunnel up command")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*out.Parameter.Value)
	if err != nil {
		pterm.Error.Printfln("Getting tunnel up command")
		return err
	}

	var tOut terraformOutput

	err = json.Unmarshal(value, &tOut)
	if err != nil {
		pterm.Error.Printfln("Getting tunnel up command")
		return err
	}

	command := strings.Split(tOut.Cmd.Value.Tunnel.Up, " ")
	command = command[:len(command)-1]
	command = append(command, "-F")
	command = append(command, fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR")))

	pterm.Success.Printfln("Getting tunnel up command")

	err = CallProcess(command[0], command[1:])
	if err != nil {
		pterm.Error.Printfln("Running tunnel up command")
		fmt.Println(err)
	}

	pterm.Success.Printfln("Running tunnel up command")

	return nil
}

func (c *tunnelCmd) BastionSHHTunnelStatus() error {
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
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", "dev")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		pterm.Error.Printfln("Getting tunnel status command")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*out.Parameter.Value)
	if err != nil {
		pterm.Error.Printfln("Getting tunnel status command")
		return err
	}

	var tOut terraformOutput

	err = json.Unmarshal(value, &tOut)
	if err != nil {
		pterm.Error.Printfln("Getting tunnel status command")
		return err
	}

	command := strings.Split(tOut.Cmd.Value.Tunnel.Status, " ")
	command = append(command, "-F")
	command = append(command, fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR")))

	pterm.Success.Printfln("Getting tunnel status command")

	err = CallProcess(command[0], command[1:])
	if err != nil {
		exiterr := err.(*exec.ExitError)
		status := exiterr.Sys().(syscall.WaitStatus)
		if status.ExitStatus() != 255 {
			pterm.Error.Printfln("Running tunnel down command")
			return err
		}
	}

	pterm.Success.Printfln("Running tunnel status command")

	return nil
}

func (c *tunnelCmd) BastionSHHTunnelDown() error {
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
		Name:           aws.String(fmt.Sprintf("/%s/terraform-output", "dev")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		pterm.Error.Printfln("Getting tunnel down command")
		return err
	}

	var value []byte

	value, err = base64.StdEncoding.DecodeString(*out.Parameter.Value)
	if err != nil {
		pterm.Error.Printfln("Getting tunnel down command")
		return err
	}

	var tOut terraformOutput

	err = json.Unmarshal(value, &tOut)
	if err != nil {
		pterm.Error.Printfln("Getting tunnel down command")
		return err
	}

	command := strings.Split(tOut.Cmd.Value.Tunnel.Down, " ")
	command = command[:len(command)-1]
	command = append(command, "-F")
	command = append(command, fmt.Sprintf("%s/ssh.config", viper.GetString("ENV_DIR")))

	pterm.Success.Printfln("Getting tunnel down command")

	err = CallProcess(command[0], command[1:])
	if err != nil {
		exiterr := err.(*exec.ExitError)
		status := exiterr.Sys().(syscall.WaitStatus)
		if status.ExitStatus() != 255 {
			pterm.Error.Printfln("Running tunnel down command")
			return err
		}
	}

	pterm.Success.Printfln("Running tunnel down command")

	return nil
}

type terraformOutput struct {
	BastionInstanceID struct {
		Value string `json:"value,omitempty"`
	} `json:"bastion_instance_id,omitempty"`
	Cmd struct {
		Value struct {
			Tunnel struct {
				Down   string `json:"down,omitempty"`
				Status string `json:"status,omitempty"`
				Up     string `json:"up,omitempty"`
			} `json:"tunnel,omitempty"`
		} `json:"value,omitempty"`
	} `json:"cmd,omitempty"`
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

func CallProcess(app string, flags []string) error {
	if app == "" {
		return errors.New("application is not specified")
	}

	errc := make(chan error, 1)

	cmd := exec.Command(app, flags...)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer close(errc)
		buf := bufio.NewReader(stderrPipe)
		for {
			line, err := buf.ReadString('\n')
			if len(line) > 0 {
				pterm.Info.Printfln(strings.TrimSuffix(line, "\n"))
			}
			if err != nil {
				return
			}
		}
	}()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
