package terraform

import (
	"io"

	"github.com/hazelops/ize/pkg/terminal"
)

type Terraform interface {
	Run() error
	RunUI(ui terminal.UI) error
	Prepare() error
	NewCmd(cmd []string)
	SetOut(out io.Writer)
}
