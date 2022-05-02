package apps

import (
	"github.com/hazelops/ize/pkg/terminal"
)

type alias struct {
	Name string
}

func NewAliasDeployment(name string) *alias {
	return &alias{
		Name: name,
	}
}

func (a *alias) Deploy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s: deployment completed!", a.Name)
	s.Done()

	return nil
}

func (a *alias) Destroy(sg terminal.StepGroup, ui terminal.UI) error {
	s := sg.Add("%s: destroy completed!", a.Name)
	s.Done()

	return nil
}
