package apps

import (
	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
)

type AliasService struct {
	Project *config.Project
	App     *config.Alias
}

func (a *AliasService) Deploy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deployment completed!", a.App.Name)
	s.Done()

	return nil
}

func (a *AliasService) Destroy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: destroy completed!", a.App.Name)
	s.Done()

	return nil
}

func (a *AliasService) Push(ui terminal.UI) error {
	return nil
}

func (a *AliasService) Build(ui terminal.UI) error {
	return nil
}

func (a *AliasService) Redeploy(ui terminal.UI) error {
	return nil
}
