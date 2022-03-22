package terraform

import (
	"fmt"
	"path/filepath"
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
	version    string
	command    []string
	env        []string
	outputPath string
}

func NewLocalTerraform(version string, command []string, env []string, out string) *local {
	return &local{
		version:    version,
		command:    command,
		env:        env,
		outputPath: out,
	}
}

func (l *local) Run() error {
	err := term.New().InteractiveRun("terraform", l.command)
	if err != nil {
		return err
	}

	return nil
}

func (l *local) Prepare() error {
	var (
		tfpath        = "/usr/local/bin/terraform"
		defaultMirror = "https://releases.hashicorp.com/terraform"
	)

	err := installVersion(l.version, &tfpath, &defaultMirror)
	if err != nil {
		return err
	}

	return nil
}

func (l *local) NewCmd(cmd []string) {
	l.command = cmd
}

func (l *local) SetOutput(path string) {
	l.outputPath = path
}

func (l *local) RunUI(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("running terraform v%s...", l.version)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 100) }()

	t := term.New(term.WithStderr(s.TermOutput()), term.WithStdout(s.TermOutput()), term.WithDir(viper.GetString("ENV_DIR")))

	err := t.InteractiveRun("terraform", l.command)
	if err != nil {
		return err
	}

	s.Done()

	return nil
}

func installVersion(version string, custBinPath *string, mirrorURL *string) error {
	if tfswitcher.ValidVersionFormat(version) {
		requestedVersion := version

		//check to see if the requested version has been downloaded before
		installLocation := tfswitcher.GetInstallLocation()
		installFileVersionPath := tfswitcher.ConvertExecutableExt(filepath.Join(installLocation, versionPrefix+requestedVersion))
		recentDownloadFile := tfswitcher.CheckFileExist(installFileVersionPath)
		if recentDownloadFile {
			tfswitcher.ChangeSymlink(installFileVersionPath, *custBinPath)
			fmt.Printf("Switched terraform to version %q \n", requestedVersion)
			tfswitcher.AddRecent(requestedVersion) //add to recent file for faster lookup
			return nil
		}

		//if the requested version had not been downloaded before
		listAll := true                                            //set list all true - all versions including beta and rc will be displayed
		tflist, _ := tfswitcher.GetTFList(*mirrorURL, listAll)     //get list of versions
		exist := tfswitcher.VersionExist(requestedVersion, tflist) //check if version exist before downloading it

		if exist {
			tfswitcher.Install(requestedVersion, *custBinPath, *mirrorURL)
		} else {
			return fmt.Errorf("the provided terraform version does not exist")
		}

	} else {
		tfswitcher.PrintInvalidTFVersion()
		return fmt.Errorf("args must be a valid terraform version")
	}

	return nil
}
