package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var (
	izeBinary = GetFromEnv("IZE_BINARY", "ize")

	examplesRootDir = GetFromEnv("IZE_EXAMPLES_PATH", "examples/simple-monorepo")
)

// A struct representation of the IZE binary
type binary struct {
	t          *testing.T
	binaryPath string
	workingDir string
}

func NewBinary(t *testing.T, binaryPath string, workingDir string) *binary {
	return &binary{
		t:          t,
		binaryPath: binaryPath,
		workingDir: workingDir,
	}
}

func GetFromEnv(key string, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		result = defaultValue
	}
	return result
}

// Runs a command with the arguments specified
func (b *binary) RunRaw(args ...string) (stdout, stderr string, err error) {
	cmd := b.NewCmd(args...)
	cmd.Stdin = nil
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}
	err = cmd.Run()
	stdout = cmd.Stdout.(*bytes.Buffer).String()
	stderr = cmd.Stderr.(*bytes.Buffer).String()
	return
}

// Builds a generic execer for running waypoint commands
func (b *binary) NewCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(b.binaryPath, args...)
	cmd.Dir = b.workingDir
	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, "CHECKPOINT_DISABLE=1")
	return cmd
}

// Runs the command, fails the test on errors
func (b *binary) Run(args string) (stdout string) {
	fmt.Printf("running %s ...\n", args)
	stdout, stderr, err := b.RunRaw(splitArgs(args)...)
	if err != nil {
		b.t.Fatalf("unexpected error running %q inside %q\nERROR:\n%s\n\nSTDERR:\n%s\n\nSTDOUT:\n%s", args, b.workingDir, err, stderr, stdout)
	}
	if stderr != "" {
		b.t.Fatalf("unexpected stderr output running %s:\n%s", args, stderr)
	}
	return stdout
}

func splitArgs(args string) []string {
	return strings.Split(args, " ")
}
