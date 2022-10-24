package serverless

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hazelops/ize/internal/config"

	"github.com/hazelops/ize/pkg/terminal"
)

type Manager struct {
	Project *config.Project
	App     *config.Serverless
}

func (sls *Manager) Nvm(ui terminal.UI, command []string) error {
	sls.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: running '%s'...", sls.App.Name, strings.Join(command, " "))
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 200) }()

	err := sls.nvm(s.TermOutput(), strings.Join(command, " "))
	if err != nil {
		return fmt.Errorf("can't run nvm: %w", err)
	}

	s.Done()
	s = sg.Add("%s: running '%s' completed!", sls.App.Name, strings.Join(command, " "))
	s.Done()

	return nil
}

func (sls *Manager) prepare() {
	if sls.App.Path == "" {
		appsPath := sls.Project.AppsPath
		if !filepath.IsAbs(appsPath) {
			appsPath = filepath.Join(os.Getenv("PWD"), appsPath)
		}

		sls.App.Path = filepath.Join(appsPath, sls.App.Name)
	} else {
		rootDir := sls.Project.RootDir

		if !filepath.IsAbs(sls.App.Path) {
			sls.App.Path = filepath.Join(rootDir, sls.App.Path)
		}
	}

	if len(sls.App.File) == 0 {
		_, err := os.Stat(filepath.Join(sls.App.Path, "serverless.ts"))
		if os.IsNotExist(err) {
			sls.App.File = "serverless.yml"
		} else {
			sls.App.File = "serverless.ts"
		}
	}

	if len(sls.App.ServerlessVersion) == 0 {
		sls.App.ServerlessVersion = "2"
	}

	if len(sls.App.SLSNodeModuleCacheMount) == 0 {
		sls.App.SLSNodeModuleCacheMount = fmt.Sprintf("%s-node-modules", sls.App.Name)
	}

	if len(sls.App.AwsProfile) == 0 {
		sls.App.AwsProfile = sls.Project.AwsProfile
	}

	if len(sls.App.AwsRegion) == 0 {
		sls.App.AwsRegion = sls.Project.AwsRegion
	}

	sls.App.Env = append(sls.App.Env, "SLS_DEBUG=*")
}

func (sls *Manager) Deploy(ui terminal.UI) error {
	sls.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: deploying app...", sls.App.Name)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 200) }()

	switch sls.Project.PreferRuntime {
	case "native":
		s.Update("%s: deploying app [run nvm use]...", sls.App.Name)

		err := sls.runNvm(s.TermOutput())
		if err != nil {
			return fmt.Errorf("can't run nvm: %w", err)
		}

		s.Done()
		s = sg.Add("%s: deploying app [run npm install]...", sls.App.Name)
		err = sls.runNpmInstall(s.TermOutput())
		if err != nil {
			return fmt.Errorf("can't run npm install: %w", err)
		}

		if sls.App.CreateDomain {
			s.Done()
			s = sg.Add("%s: deploying app [run serverless create_domain]...", sls.App.Name)
			err = sls.runCreateDomain(s.TermOutput())
			if err != nil {
				return fmt.Errorf("can't run serverless create_domain: %w", err)
			}
		}

		s.Done()
		s = sg.Add("%s: deploying app [run serverless deploy]...", sls.App.Name)
		err = sls.runDeploy(s.TermOutput())
		if err != nil {
			return fmt.Errorf("can't run serverless deploy: %w", err)
		}
	case "docker":
		err := sls.deployWithDocker(s)
		if err != nil {
			return err
		}
	}

	s.Done()
	s = sg.Add("%s: deployment completed!", sls.App.Name)
	s.Done()

	return nil
}

func (sls *Manager) Destroy(ui terminal.UI) error {
	sls.prepare()

	sg := ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("%s: destroying app...", sls.App.Name)
	defer func() { s.Abort(); time.Sleep(time.Millisecond * 200) }()

	switch sls.Project.PreferRuntime {
	case "native":
		s.Update("%s: destroying app [run nvm use]...", sls.App.Name)

		err := sls.runNvm(s.TermOutput())
		if err != nil {
			return fmt.Errorf("can't run nvm: %w", err)
		}

		s.Done()
		s = sg.Add("%s: destroying app [run npm install]...", sls.App.Name)
		err = sls.runNpmInstall(s.TermOutput())
		if err != nil {
			return fmt.Errorf("can't run npm install: %w", err)
		}

		if sls.App.CreateDomain {
			s.Done()
			s = sg.Add("%s: destroying app [run serverless delete_domain]...", sls.App.Name)
			err = sls.runRemoveDomain(s.TermOutput())
			if err != nil {
				return fmt.Errorf("can't run serverless delete_domain: %w", err)
			}
		}

		s.Done()
		s = sg.Add("%s: destroying app [run serverless remove]...", sls.App.Name)
		err = sls.runRemove(s.TermOutput())
		if err != nil {
			return fmt.Errorf("can't run serverless deploy: %w", err)
		}
	case "docker":
		err := sls.removeWithDocker(s)
		if err != nil {
			return err
		}
	}

	s.Done()
	s = sg.Add("%s: destroy completed!", sls.App.Name)
	s.Done()

	return nil
}

func (sls *Manager) Push(ui terminal.UI) error {
	return nil
}

func (sls *Manager) Build(ui terminal.UI) error {
	return nil
}

func (sls *Manager) Redeploy(ui terminal.UI) error {
	return nil
}
