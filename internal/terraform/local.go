package terraform

import (
	"bufio"
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/hazelops/ize/internal/config"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hazelops/ize/pkg/term"
	tfswitcher "github.com/psihachina/terraform-switcher/lib"
)

const (
	versionPrefix = "terraform_"
	defaultMirror = "https://releases.hashicorp.com/terraform"
)

type local struct {
	version string
	command []string
	env     []string
	output  io.Writer
	tfpath  string
	project *config.Project
	state   string
}

func NewLocalTerraform(state string, command []string, env []string, out io.Writer, project *config.Project) *local {
	if len(project.Terraform[state].Version) == 0 {
		project.Terraform[state].Version = project.TerraformVersion
	}

	return &local{
		state:   state,
		version: project.Terraform[state].Version,
		command: command,
		env:     env,
		output:  out,
		project: project,
	}
}

func (l *local) Run() error {
	if len(l.project.EnvDir) == 0 {
		l.project.EnvDir = "."
	}

	stateDir := filepath.Join(l.project.EnvDir, l.state)
	if l.state == "infra" {
		stateDir = l.project.EnvDir
	}

	cmd := exec.Command(l.tfpath, l.command...)
	cmd.Dir = stateDir

	err := term.New(term.WithDir(l.project.EnvDir), term.WithStdin(os.Stdin)).InteractiveRun(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (l *local) Prepare() error {
	var (
		mirror = defaultMirror
		path   = ""
	)

	path, err := installVersion(l.version, &mirror)
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

	cmd.Stdout = out
	cmd.Stderr = out

	if err = cmd.Start(); err != nil {
		return
	}

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

func (l *local) RunUI(ui *echelon.Logger) error {
	s := ui.Scoped(fmt.Sprintf("[%s][%s] terraform %s (v%s)", l.project.Env, l.state, strings.Join(l.command, " "), l.version))

	var stdout io.Writer
	stdout = s.AsWriter(echelon.InfoLevel)
	if l.output != nil {
		stdout = l.output
	}

	if len(l.project.EnvDir) == 0 {
		l.project.EnvDir = "."
	}

	stateDir := filepath.Join(l.project.EnvDir, l.state)
	if l.state == "infra" {
		stateDir = l.project.EnvDir
	}

	cmd := exec.Command(l.tfpath, l.command...)
	cmd.Dir = stateDir
	_, _, err := runCommand(cmd, stdout)

	if err != nil {
		return err
	}

	s.Finish(true)
	time.Sleep(time.Millisecond)

	return nil
}

func getInstallLocation(installPath string) string {
	/* get current user */
	usr, errCurr := user.Current()
	if errCurr != nil {
		log.Fatal(errCurr)
	}

	userCommon := usr.HomeDir

	/* set installation location */
	installLocation := filepath.Join(userCommon, installPath)

	/* Create local installation directory if it does not exist */
	tfswitcher.CreateDirIfNotExist(installLocation)

	return installLocation

}

func installVersion(version string, mirrorURL *string) (string, error) {
	var (
		installFileVersionPath string
	)

	if tfswitcher.ValidVersionFormat(version) {
		requestedVersion := version

		//check to see if the requested version has been downloaded before
		installLocation := getInstallLocation(".ize/versions/terraform/")
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
			err := Install(requestedVersion, *mirrorURL)
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

func Install(tfversion string, mirrorURL string) error {
	installLocation := getInstallLocation(".ize/versions/terraform/") //get installation location -  this is where we will put our terraform binary file

	goarch := runtime.GOARCH
	goos := runtime.GOOS

	// Terraform darwin arm64 comes with version 1.0.2 and above
	tfver, _ := version.NewVersion(tfversion)
	tf102, _ := version.NewVersion("1.0.2")
	if goos == "darwin" && goarch == "arm64" && tfver.LessThan(tf102) {
		goarch = "amd64"
	}

	/* check if selected version has already been downloaded */
	installFileVersionPath := tfswitcher.ConvertExecutableExt(filepath.Join(installLocation, versionPrefix+tfversion))
	fileExist := tfswitcher.CheckFileExist(installFileVersionPath)

	/* if selected version already exists */
	if fileExist {
		return nil
	}

	// if it does not have a slash - append it
	hasSlash := strings.HasSuffix(mirrorURL, "/")
	if !hasSlash {
		mirrorURL = fmt.Sprintf("%s/", mirrorURL)
	}

	/* if selected version already exist, */
	/* proceed to download it from the hashicorp release page */
	url := mirrorURL + tfversion + "/" + versionPrefix + tfversion + "_" + goos + "_" + goarch + ".zip"
	zipFile, errDownload := tfswitcher.DownloadFromURL(installLocation, url)

	/* If unable to download file from url, exit(1) immediately */
	if errDownload != nil {
		return errDownload
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := tfswitcher.Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("[Error] : Unable to unzip downloaded zip file")
		return errUnzip
	}

	/* rename unzipped file to terraform version name - terraform_x.x.x */
	installFilePath := tfswitcher.ConvertExecutableExt(filepath.Join(installLocation, "terraform"))
	tfswitcher.RenameFile(installFilePath, installFileVersionPath)

	/* remove zipped file to clear clutter */
	tfswitcher.RemoveFiles(zipFile)

	return nil
}
