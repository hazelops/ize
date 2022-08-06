//go:build windows
// +build windows

package term

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// InteractiveRun runs the input command that starts a child process.
func (r Runner) InteractiveRun(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {
	sig := make(chan os.Signal, 1)
	// See https://golang.org/pkg/os/signal/#hdr-Windows
	signal.Notify(sig, os.Interrupt)
	defer signal.Reset(os.Interrupt)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	var bufOut, bufErr bytes.Buffer
	outReader2 := io.TeeReader(outReader, &bufOut)
	errReader2 := io.TeeReader(errReader, &bufErr)

	if err = cmd.Start(); err != nil {
		return
	}

	go r.printOutputWithHeader("", outReader2)
	go r.printOutputWithHeader("", errReader2)

	err = cmd.Wait()

	stdout = bufOut.String()
	stderr = bufErr.String()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}
	return
}
