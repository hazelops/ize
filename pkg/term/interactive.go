//go:build !windows
// +build !windows

package term

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

func (r Runner) InteractiveRun(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {
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

func (r Runner) printOutputWithHeader(header string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if r.stdout != nil {
			fmt.Fprintf(r.stdout, "%s%s\n", header, scanner.Text())
		} else {
			fmt.Printf("%s%s\n", header, scanner.Text())
		}
	}
}
