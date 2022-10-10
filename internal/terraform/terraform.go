package terraform

import (
	"github.com/cirruslabs/echelon"
	"io"
)

type Terraform interface {
	Run() error
	RunUI(ui *echelon.Logger) error
	Prepare() error
	NewCmd(cmd []string)
	SetOut(out io.Writer)
}
