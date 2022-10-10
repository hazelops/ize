package alias

import (
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/hazelops/ize/internal/config"
	"time"
)

type Manager struct {
	Project *config.Project
	App     *config.Alias
}

func (a *Manager) Deploy(ui *echelon.Logger) error {
	s := ui.Scoped(fmt.Sprintf("%s: deployment completed!", a.App.Name))
	s.Finish(true)
	time.Sleep(time.Millisecond * 50)

	return nil
}

func (a *Manager) Destroy(ui *echelon.Logger) error {
	s := ui.Scoped(fmt.Sprintf("%s: destroy completed!", a.App.Name))
	s.Finish(true)
	time.Sleep(time.Millisecond * 50)

	return nil
}

func (a *Manager) Push(ui *echelon.Logger) error {
	return nil
}

func (a *Manager) Build(ui *echelon.Logger) error {
	return nil
}

func (a *Manager) Redeploy(ui *echelon.Logger) error {
	return nil
}
