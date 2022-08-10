package term

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
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

func (r Runner) Run(name string, args []string, options ...Option) error {
	cmd := exec.Command(name, args...)

	cmd.Wait()
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	cmd.Dir = r.dir

	for _, opt := range options {
		opt(cmd)
	}
	return cmd.Run()
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
