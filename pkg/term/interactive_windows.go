//go:build windows
// +build windows

package term

import (
	"os"
	"os/exec"
	"os/signal"
)

// InteractiveRun runs the input command that starts a child process.
func (r Runner) InteractiveRun(cmd *exec.Cmd) (err error) {
	sig := make(chan os.Signal, 1)
	// See https://golang.org/pkg/os/signal/#hdr-Windows
	signal.Notify(sig, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	cmd.Dir = r.dir
	if r.stdin != nil {
		cmd.Stdin = r.stdin
	}
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
