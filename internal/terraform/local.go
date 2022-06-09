package terraform

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hazelops/ize/internal/terraform/tfswitcher"
	"github.com/hazelops/ize/pkg/term"
	"github.com/hazelops/ize/pkg/terminal"
	"github.com/spf13/viper"
)

const (
	versionPrefix = "terraform_"
)

type local struct {
	version string
	command []string
	env     []string
	output  io.Writer
	tfpath  string
}

func NewLocalTerraform(version string, command []string, env []string, out io.Writer) *local {
	return &local{
		version: version,
		command: command,
		env:     env,
		output:  out,
	}
}

func (l *local) Run() error {
	err := term.New(term.WithDir(viper.GetString("ENV_DIR"))).InteractiveRun(l.tfpath, l.command)
	if err != nil {
		return err
	}

	return nil
}

func (l *local) Prepare() error {
	var (
		tfpath        = "/usr/local/bin/terraform"
		defaultMirror = "https://releases.hashicorp.com/terraform"
		path          = ""
	)

	path, err := installVersion(l.version, &tfpath, &defaultMirror)
	if err != nil {
		return err
	}

	l.tfpath = path

	return nil
}

func (l *local) NewCmd(cmd []string) {
	l.command = cmd
}

func (l *local) SetOut(out io.Writer) {
	l.output = out
}

func printOutput(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		w.Write([]byte(scanner.Text() + "\n"))
	}
}

func runCommand(cmd *exec.Cmd, out io.Writer) (stdout, stderr string, err error) {
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	go printOutput(outReader, out)
	go printOutput(errReader, out)

	err = cmd.Wait()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if s, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				err = fmt.Errorf("exit status: %d", s.ExitStatus())
			}
		}
	}
	return
}

func (l *local) RunUI(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Running terraform v%s...", l.version)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 100) }()

	stdout := s.TermOutput()
	if l.output != nil {
		stdout = l.output
	}

	cmd := exec.Command(l.tfpath, l.command...)
	cmd.Dir = viper.GetString("ENV_DIR")
	_, _, err := runCommand(cmd, stdout)

	if err != nil {
		return err
	}

	s.Done()

	return nil
}

func installVersion(version string, custBinPath *string, mirrorURL *string) (string, error) {
	var (
		installFileVersionPath string
		err                    error
	)

	if tfswitcher.ValidVersionFormat(version) {
		requestedVersion := version

		//check to see if the requested version has been downloaded before
		installLocation := tfswitcher.GetInstallLocation()
		installFileVersionPath = tfswitcher.ConvertExecutableExt(filepath.Join(installLocation, versionPrefix+requestedVersion))
		recentDownloadFile := tfswitcher.CheckFileExist(installFileVersionPath)
		if recentDownloadFile {
			return installFileVersionPath, nil
		}

		//if the requested version had not been downloaded before
		listAll := true                                            //set list all true - all versions including beta and rc will be displayed
		tflist, _ := tfswitcher.GetTFList(*mirrorURL, listAll)     //get list of versions
		exist := tfswitcher.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

		if exist {
			installFileVersionPath, err = tfswitcher.Install(requestedVersion, *custBinPath, *mirrorURL)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("provided terraform version does not exist")
		}

	} else {
		tfswitcher.PrintInvalidTFVersion()
		return "", fmt.Errorf("argument must be a valid terraform version")
	}

	return installFileVersionPath, nil
}
