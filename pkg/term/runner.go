package term

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type Runner struct {
	stdout io.Writer
	stderr io.Writer
	dir    string
}

type RunnerOption func(*Runner)

func WithStdout(stdout io.Writer) RunnerOption {
	return func(r *Runner) {
		r.stdout = stdout
	}
}

func WithStderr(stderr io.Writer) RunnerOption {
	return func(r *Runner) {
		r.stderr = stderr
	}
}

func WithDir(path string) RunnerOption {
	return func(r *Runner) {
		r.dir = path
	}
}

func New(opts ...RunnerOption) *Runner {
	r := &Runner{
		stderr: os.Stderr,
		stdout: os.Stdout,
		dir:    ".",
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

type Option func(cmd *exec.Cmd)

func (r Runner) Run(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {
	cmd.Stdin = os.Stdin
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
