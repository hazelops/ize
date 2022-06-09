package terraform

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hazelops/ize/pkg/term"
	"github.com/hazelops/ize/pkg/terminal"
	tfswitcher "github.com/psihachina/terraform-switcher/lib"
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
		defaultMirror = "https://releases.hashicorp.com/terraform"
		path          = ""
	)

	path, err := installVersion(l.version, &defaultMirror)
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

func (l *local) RunUI(ui terminal.UI) error {
	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Running terraform v%s...", l.version)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 100) }()

	stdout := s.TermOutput()
	if l.output != nil {
		stdout = l.output
	}

	t := term.New(
		term.WithStderr(s.TermOutput()),
		term.WithStdout(stdout),
		term.WithDir(viper.GetString("ENV_DIR")),
	)

	err := t.InteractiveRun(l.tfpath, l.command)
	if err != nil {
		return err
	}

	s.Done()

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
			Install(requestedVersion, *mirrorURL)
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
		os.Exit(0)
	}

	//if does not have slash - append slash
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
		fmt.Println(errDownload)
		os.Exit(1)
	}

	/* unzip the downloaded zipfile */
	_, errUnzip := tfswitcher.Unzip(zipFile, installLocation)
	if errUnzip != nil {
		fmt.Println("[Error] : Unable to unzip downloaded zip file")
		log.Fatal(errUnzip)
		os.Exit(1)
	}

	/* rename unzipped file to terraform version name - terraform_x.x.x */
	installFilePath := tfswitcher.ConvertExecutableExt(filepath.Join(installLocation, "terraform"))
	tfswitcher.RenameFile(installFilePath, installFileVersionPath)

	/* remove zipped file to clear clutter */
	tfswitcher.RemoveFiles(zipFile)

	return nil
}
