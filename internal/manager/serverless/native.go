package serverless

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hazelops/ize/pkg/term"
	"github.com/sirupsen/logrus"
)

func (sls *Manager) runNpmInstall(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf("source %s/nvm.sh && nvm use %s && npm install --save-dev", nvmDir, sls.App.NodeVersion)

	if sls.App.UseYarn {
		command = npmToYarn(command)
	}

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (sls *Manager) nvm(w io.Writer, command string) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}
	err := sls.readNvmrc()
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c",
		fmt.Sprintf("source %s/nvm.sh && nvm install %s && %s", nvmDir, sls.App.NodeVersion, command),
	)

	return term.New(
		term.WithDir(sls.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func (sls *Manager) readNvmrc() error {
	_, err := os.Stat(filepath.Join(sls.App.Path, ".nvmrc"))
	if os.IsNotExist(err) {
	} else {
		file, err := os.ReadFile(filepath.Join(sls.App.Path, ".nvmrc"))
		if err != nil {
			return fmt.Errorf("can't read .nvmrc: %w", err)
		}
		sls.App.NodeVersion = strings.TrimSpace(string(file))
	}
	return nil
}

func (sls *Manager) runNvm(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	err := sls.readNvmrc()
	if err != nil {
		return err
	}

	command := fmt.Sprintf("source %s/nvm.sh && nvm install %s", nvmDir, sls.App.NodeVersion)

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	return term.New(
		term.WithDir(sls.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func (sls *Manager) runDeploy(w io.Writer) error {

	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}
	var command string

	// SLS v3 has breaking changes in syntax
	if sls.App.ServerlessVersion == "3" {
		command = fmt.Sprintf(
			`source %s/nvm.sh && 
				nvm use %s &&
				npx serverless deploy \
				--config=%s \
				--param="service=%s" \
				--region=%s \
				--aws-profile=%s \
				--stage=%s \
				--verbose`,
			nvmDir, sls.App.NodeVersion, sls.App.File,
			sls.App.Name, sls.App.AwsRegion,
			sls.App.AwsProfile, sls.Project.Env)
	} else {
		command = fmt.Sprintf(
			`source %s/nvm.sh && 
				nvm use %s &&
				npx serverless deploy \
				--config %s \
				--service %s \
				--verbose \
				--region %s \
				--aws-profile %s \
				--stage %s`,
			nvmDir, sls.App.NodeVersion, sls.App.File,
			sls.App.Name, sls.App.AwsRegion,
			sls.App.AwsProfile, sls.Project.Env)
	}

	if sls.App.UseYarn {
		command = npmToYarn(command)
	}

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	return term.New(
		term.WithDir(sls.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func (sls *Manager) runRemove(w io.Writer) error {

	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	var command string

	// SLS v3 has breaking changes in syntax
	if sls.App.ServerlessVersion == "3" {
		command = fmt.Sprintf(
			`source %s/nvm.sh && \
				nvm use %s && \
				npx serverless remove \
				--config=%s \
				--param="service=%s" \
				--region=%s \
				--aws-profile=%s \
				--stage=%s \
				--verbose`,
			nvmDir, sls.App.NodeVersion, sls.App.File,
			sls.App.Name, sls.App.AwsRegion,
			sls.App.AwsProfile, sls.Project.Env)
	} else {
		command = fmt.Sprintf(
			`source %s/nvm.sh && \
				nvm use %s && \
				npx serverless remove \
				--config %s \
				--service %s \
				--verbose \
				--region %s \
				--aws-profile %s \
				--stage %s`,
			nvmDir, sls.App.NodeVersion, sls.App.File,
			sls.App.Name, sls.App.AwsRegion,
			sls.App.AwsProfile, sls.Project.Env)
	}

	if sls.App.UseYarn {
		command = npmToYarn(command)
	}

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	return term.New(
		term.WithDir(sls.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func (sls *Manager) runCreateDomain(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf(
		`source %s/nvm.sh && \
				nvm use %s && \
				npx serverless create_domain \
				--verbose \
				--region %s \
				--aws-profile %s \
				--stage %s`,
		nvmDir, sls.App.NodeVersion, sls.App.AwsRegion,
		sls.App.AwsProfile, sls.Project.Env)

	if sls.App.UseYarn {
		command = npmToYarn(command)
	}

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	return term.New(
		term.WithDir(sls.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func (sls *Manager) runRemoveDomain(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf(
		`source %s/nvm.sh && \
				nvm use %s && \
				npx serverless delete_domain \
				--verbose \
				--region %s \
				--aws-profile %s \
				--stage %s`,
		nvmDir, sls.App.NodeVersion, sls.App.AwsRegion,
		sls.App.AwsProfile, sls.Project.Env)
		

	if sls.App.UseYarn {
		command = npmToYarn(command)
	}

	logrus.SetOutput(w)
	logrus.Debugf("command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	return term.New(
		term.WithDir(sls.App.Path),
		term.WithStdout(w),
		term.WithStderr(w),
	).InteractiveRun(cmd)
}

func npmToYarn(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "npm", "yarn")
	return strings.ReplaceAll(cmd, "npx", "yarn")
}
