package apps

import (
	"github.com/hazelops/ize/pkg/terminal"
)

type alias struct {
	Name string
}

func NewAliasApp(name string) *alias {
	return &alias{
		Name: name,
	}
}

func (a *alias) Deploy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deployment completed!", a.Name)
	s.Done()

	return nil
}

func (a *alias) Destroy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: destroy completed!", a.Name)
	s.Done()

	return nil
}

func (a *alias) Push(ui terminal.UI) error {
	return nil
}

func (a *alias) Build(ui terminal.UI) error {
	return nil
}

func (a *alias) Redeploy(ui terminal.UI) error {
	return nil
}
