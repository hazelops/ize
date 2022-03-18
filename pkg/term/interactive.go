//go:build !windows
// +build !windows

package term

import (
	"os"
	"os/exec"
	"os/signal"
)

func (r Runner) InteractiveRun(name string, args []string) error {
	// Ignore interrupt signal otherwise the program exits.
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	cmd := exec.Command(name, args...)
	cmd.Dir = r.dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
