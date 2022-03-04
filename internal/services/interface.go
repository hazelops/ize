package services

import "github.com/hazelops/ize/pkg/terminal"

type Deployment interface {
	Deploy(sg terminal.StepGroup, ui terminal.UI) error
	Destroy(sg terminal.StepGroup, ui terminal.UI) error
}
