package ssmsession

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/hazelops/ize/pkg/term"
)

const (
	ssmPluginBinaryName = "session-manager-plugin"
	startSessionAction  = "StartSession"
)

type SSMPlugingRunner interface {
	InteractiveRun(name string, args []string) error
	BackgroundRun(name string, args []string) error
}

type SSMPluginCommand struct {
	runner SSMPlugingRunner
	region string
}

func NewSSMPluginCommand(region string) SSMPluginCommand {
	return SSMPluginCommand{
		runner: term.New(),
		region: region,
	}
}

func (s SSMPluginCommand) Start(ssmSession *ecs.Session) error {
	response, err := json.Marshal(ssmSession)
	if err != nil {
		return fmt.Errorf("marshal session response: %w", err)
	}
	if err := s.runner.InteractiveRun(ssmPluginBinaryName,
		[]string{string(response), s.region, startSessionAction}); err != nil {
		return fmt.Errorf("start session: %w", err)
	}
	return nil
}

func (s SSMPluginCommand) Forward(ssmSession *ssm.StartSessionOutput, sessionInput *ssm.StartSessionInput) error {
	response, err := json.Marshal(ssmSession)
	if err != nil {
		return fmt.Errorf("marshal session response: %w", err)
	}

	params, err := json.Marshal(sessionInput)
	if err != nil {
		return fmt.Errorf("marshal session response: %w", err)
	}

	err = s.runner.BackgroundRun(ssmPluginBinaryName, []string{string(response), s.region, startSessionAction, "", string(params), *ssmSession.StreamUrl})
	if err != nil {
		return err
	}

	return nil
}
