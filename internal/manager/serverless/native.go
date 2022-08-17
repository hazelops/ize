package serverless

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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

	command := fmt.Sprintf("source %s/nvm.sh && nvm install %s", nvmDir, sls.App.NodeVersion)

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

func (sls *Manager) runDeploy(w io.Writer) error {

	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf(
		`source %s/nvm.sh && 
				nvm use %s &&
				npx serverless deploy \
				--config %s \
				--service %s \
				--verbose \
				--region %s \
				--profile %s \
				--stage %s`,
		nvmDir, sls.App.NodeVersion, sls.App.File,
		sls.App.Name, sls.App.AwsRegion,
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

func (sls *Manager) runRemove(w io.Writer) error {

	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf(
		`source %s/nvm.sh && 
				nvm use %s &&
				npx serverless remove \
				--config %s \
				--service %s \
				--verbose \
				--region %s \
				--profile %s \
				--stage %s`,
		nvmDir, sls.App.NodeVersion, sls.App.File,
		sls.App.Name, sls.App.AwsRegion,
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

func (sls *Manager) runCreateDomain(w io.Writer) error {
	nvmDir := os.Getenv("NVM_DIR")
	if len(nvmDir) == 0 {
		nvmDir = "$HOME/.nvm"
	}

	command := fmt.Sprintf(
		`source %s/nvm.sh && 
				nvm use %s &&
				npx serverless create_domain \
				--verbose \
				--region %s \
				--profile %s \
				--stage %s`,
		nvmDir, sls.App.Name, sls.App.AwsRegion,
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
		`source %s/nvm.sh && 
				nvm use %s &&
				npx serverless remove_domain \
				--verbose \
				--region %s \
				--profile %s \
				--stage %s`,
		nvmDir, sls.App.Name, sls.App.AwsRegion,
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
