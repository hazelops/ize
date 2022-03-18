//go:build windows
// +build windows

package term

import (
	"os"
	"os/exec"
	"os/signal"
)

// InteractiveRun runs the input command that starts a child process.
func (r Runner) InteractiveRun(name string, args []string) error {
	sig := make(chan os.Signal, 1)
	// See https://golang.org/pkg/os/signal/#hdr-Windows
	signal.Notify(sig, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	cmd := exec.Command(name, args...)
	cmd.Dir = r.dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
