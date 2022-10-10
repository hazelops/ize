package manager

import (
	"github.com/cirruslabs/echelon"
)

type Manager interface {
	Deploy(ui *echelon.Logger) error
	Destroy(ui *echelon.Logger) error
	Build(ui *echelon.Logger) error
	Push(ui *echelon.Logger) error
	Redeploy(ui *echelon.Logger) error
}
