package ssmsession

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Netflix/go-expect"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hazelops/ize/pkg/term"
)

const (
	ssmPluginBinaryName = "session-manager-plugin"
	startSessionAction  = "StartSession"
)

type SSMPluginRunner interface {
	Run(cmd *exec.Cmd) (stdout string, stderr string, exitCode int, err error)
	InteractiveRun(cmd *exec.Cmd) (err error)
}

type SSMPluginCommand struct {
	runner SSMPluginRunner
	region string
}

func NewSSMPluginCommand(region string) SSMPluginCommand {
	return SSMPluginCommand{
		runner: term.New(term.WithStdin(os.Stdin)),
		region: region,
	}
}

func (s SSMPluginCommand) Start(ssmSession *ecs.Session) error {
	var output bytes.Buffer
	response, err := json.Marshal(ssmSession)
	if err != nil {
		return fmt.Errorf("marshal session response: %w", err)
	}

	c, err := expect.NewConsole(expect.WithStdout(&output), expect.WithStdin(os.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cmd := exec.Command(ssmPluginBinaryName, []string{string(response), s.region, startSessionAction}...)
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	go func() {
		c.ExpectEOF()
	}()

	_, _, _, err = s.Run(cmd)
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}
	fmt.Println(strings.TrimSpace(output.String()))
	if strings.Contains(output.String(), "ERROR") {
		return fmt.Errorf("exit status: 1")
	}

	return nil
}

func (s SSMPluginCommand) StartInteractive(ssmSession *ecs.Session) error {
	response, err := json.Marshal(ssmSession)
	if err != nil {
		return fmt.Errorf("marshal session response: %w", err)
	}
	cmd := exec.Command(ssmPluginBinaryName, []string{string(response), s.region, startSessionAction}...)
	err = s.runner.InteractiveRun(cmd)
	if err != nil {
		return fmt.Errorf("start session: %w", err)
	}

	return nil
}

func (s SSMPluginCommand) Run(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {

	if err = cmd.Start(); err != nil {
		return
	}

	err = cmd.Wait()

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
