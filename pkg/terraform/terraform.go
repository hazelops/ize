package terraform

import "github.com/hazelops/ize/pkg/terminal"

type Terraform interface {
	Run() error
	RunUI(ui terminal.UI) error
	Prepare() error
	NewCmd(cmd []string)
	SetOutput(path string)
}
