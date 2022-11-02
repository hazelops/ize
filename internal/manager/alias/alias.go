package alias

import (
	"time"

	"github.com/hazelops/ize/internal/config"
	"github.com/hazelops/ize/pkg/terminal"
)

type Manager struct {
	Project *config.Project
	App     *config.Alias
}

func (a *Manager) Deploy(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deployment completed!", a.App.Name)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 200) }()
	s.Done()

	time.Sleep(time.Millisecond * 200)

	return nil
}

func (a *Manager) Destroy(ui terminal.UI, autoApprove bool) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: destroy completed!", a.App.Name)
	s.Done()

	return nil
}

func (a *Manager) Push(ui terminal.UI) error {
	return nil
}

func (a *Manager) Build(ui terminal.UI) error {
	return nil
}

func (a *Manager) Redeploy(ui terminal.UI) error {
	return nil
}

func (e *Manager) Explain() error {
	return nil
}
