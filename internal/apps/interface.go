package apps

import (
	"github.com/hazelops/ize/pkg/terminal"
)

type App interface {
	Deploy(ui terminal.UI) error
	Destroy(ui terminal.UI) error
	Build(ui terminal.UI) error
	Push(ui terminal.UI) error
	Redeploy(ui terminal.UI) error
}
