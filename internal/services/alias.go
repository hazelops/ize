package services

import (
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/mitchellh/mapstructure"
)

type alias struct {
	Name string
}

func NewAliasDeployment(app App) *alias {
	var aliasConfig alias

	mapstructure.Decode(app, &aliasConfig)

	return &aliasConfig
}

func (a *alias) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s deployment completed!", a.Name)
	s.Done()

	return nil
}

func (a *alias) Destroy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s destroy completed!", a.Name)
	s.Done()

	return nil
}
