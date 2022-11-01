package manager

import (
	"github.com/hazelops/ize/pkg/terminal"
)

type Manager interface {
	Deploy(ui terminal.UI) error
	Destroy(ui terminal.UI, autoApprove bool) error
	Build(ui terminal.UI) error
	Push(ui terminal.UI) error
	Redeploy(ui terminal.UI) error
	Explain() error
}
