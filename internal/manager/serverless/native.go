package serverless

import (
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
)

func (sls *Manager) deployWithNative(w io.Writer) error {
	err := sls.runNpmInstall(w)
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	err = sls.runNpmInstall(w)
	if err != nil {
		return fmt.Errorf("can't run npm install: %w", err)
	}

	if sls.App.CreateDomain {
		err = sls.runCreateDomain(w)
		if err != nil {
			return fmt.Errorf("can't run serverless create_domain: %w", err)
		}
	}

	err = sls.runDeploy(w)
	if err != nil {
		return fmt.Errorf("can't run serverless deploy: %w", err)
	}

	return nil
}

func (sls *Manager) runNpmInstall(w io.Writer) error {
	//npmInstallCommand := "npm install --save-dev"

	cmd := exec.Command("npm", "install", "--save-dev")
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	return nil
}

func (sls *Manager) runNvm(w io.Writer) error {
	nvmCommand := fmt.Sprintf("source $HOME/.nvm/nvm.sh && nvm install %s && nvm use %s", sls.App.NodeVersion, sls.App.NodeVersion)

	cmd := exec.Command("bash", "-c", nvmCommand)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	return nil
}

func (sls *Manager) runDeploy(w io.Writer) error {
	cmd := exec.Command(
		"npx",
		"serverless",
		"deploy",
		"--config", sls.App.File,
		"--service", sls.App.Name,
		"--verbose",
		"--region", sls.Project.AwsRegion,
		"--profile", sls.Project.AwsProfile,
		"--env", sls.Project.Env,
	)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	return nil
}

func (sls *Manager) runRemove(w io.Writer) error {
	cmd := exec.Command(
		"npx",
		"serverless",
		"remove",
		"--config", sls.App.File,
		"--service", sls.App.Name,
		"--verbose",
		"--region", sls.Project.AwsRegion,
		"--profile", sls.Project.AwsProfile,
		"--env", sls.Project.Env,
	)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	return nil
}

func (sls *Manager) runCreateDomain(w io.Writer) error {
	cmd := exec.Command(
		"npx",
		"serverless",
		"create_domain",
		"--verbose",
		"--region", sls.Project.AwsRegion,
		"--profile", sls.Project.AwsProfile,
		"--env", sls.Project.Env,
	)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	return nil
}

func (sls *Manager) runRemoveDomain(w io.Writer) error {
	cmd := exec.Command(
		"npx",
		"serverless",
		"remove_domain",
		"--verbose",
		"--region", sls.Project.AwsRegion,
		"--profile", sls.Project.AwsProfile,
		"--env", sls.Project.Env,
	)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Dir = filepath.Join(sls.App.Path)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	return nil
}
