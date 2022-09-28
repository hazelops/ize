package serverless

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (sls *Manager) runNpmInstall(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf("source %s/nvm.sh && nvm use %s && npm install --save-dev", nvmDir, sls.App.NodeVersion)

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

func (sls *Manager) runNvm(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}
	var command string
	_, err := os.Stat(filepath.Join(sls.App.Path, ".nvmrc"))
	if os.IsNotExist(err) {
		command = fmt.Sprintf("source %s/nvm.sh && nvm install %s", nvmDir, sls.App.NodeVersion)

	} else {
		file, err := os.ReadFile(filepath.Join(sls.App.Path, ".nvmrc"))
		if err != nil {
			return fmt.Errorf("can't read .nvmrc: %w", err)
		}
		sls.App.NodeVersion = strings.TrimSpace(string(file))
		command = fmt.Sprintf("source %s/nvm.sh && nvm install %s", nvmDir, sls.App.NodeVersion)
	}

	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
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
