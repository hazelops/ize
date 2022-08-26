package term

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	expect "github.com/Netflix/go-expect"
)

type Runner struct {
	stdout io.Writer
	stderr io.Writer
	dir    string
	stdin  io.Reader
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

func WithStdin(stdin io.Reader) RunnerOption {
	return func(r *Runner) {
		r.stdin = stdin
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
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		return "", "", 0, err
	}
	defer c.Close()

	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	if r.stdin != nil {
		cmd.Stdin = r.stdin
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		stdout, _ = c.ExpectEOF()
	}()

	if err = cmd.Start(); err != nil {
		return
	}

	err = cmd.Wait()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}

	c.Tty().Close()
	<-done

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
