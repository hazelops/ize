package services

import (
	"fmt"

	"github.com/hazelops/ize/pkg/terminal"
)

type App struct {
	Name      string
	Type      string
	Path      string
	DependsOn []string               `mapstructure:"depends_on"`
	Body      map[string]interface{} `mapstructure:",remain"`
}

func (a *App) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	var deployment Deployment

	switch a.Type {
	case "ecs":
		deployment = NewECSDeployment(*a)
	case "serverless":
		deployment = NewServerlessDeployment(*a)
	case "alias":
		deployment = NewAliasDeployment(*a)
	default:
		return fmt.Errorf("services type of %s not supported", a.Type)
	}

	err := deployment.Deploy(sg, ui)
	if err != nil {
		return err
	}

	return nil
}
