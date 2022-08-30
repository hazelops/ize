//go:build !windows
// +build !windows

package term

import (
	"os"
	"os/exec"
	"os/signal"
)

func (r Runner) InteractiveRun(cmd *exec.Cmd) (err error) {
	// Ignore interrupt signal otherwise the program exits.
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	cmd.Dir = r.dir
	if r.stdin != nil {
		cmd.Stdin = r.stdin
	}
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
