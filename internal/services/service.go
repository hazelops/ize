package services

import (
	"fmt"

	"github.com/hazelops/ize/pkg/terminal"
)

type Service struct {
	Name      string
	Type      string
	Path      string
	DependsOn []string               `mapstructure:"depends_on"`
	Body      map[string]interface{} `mapstructure:",remain"`
}

func (svs *Service) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	var deployment Deployment

	switch svs.Type {
	case "ecs":
		deployment = NewECSDeployment(*svs)
	case "serverless":
		deployment = NewServerlessDeployment(*svs)
	default:
		return fmt.Errorf("services type of %s not supported", svs.Type)
	}

	err := deployment.Deploy(sg, ui)
	if err != nil {
		return err
	}

	return nil
}
