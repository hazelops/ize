package serverless

import (
	"fmt"
	"github.com/cirruslabs/echelon"
	"github.com/hazelops/ize/internal/config"
	"os"
	"path/filepath"
)

type Manager struct {
	Project *config.Project
	App     *config.Serverless
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

type Writer struct {
	logger *echelon.Logger
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.logger.Infof(string(p))
	return len(p), err
}

func (w *Writer) WriteHeader(status int) {
	return
}

func (sls *Manager) Deploy(ui *echelon.Logger) error {
	sls.prepare()

	sg := ui.Scoped(fmt.Sprintf("%s: deploying app...", sls.App.Name))

	switch sls.Project.PreferRuntime {
	case "native":
		s := sg.Scoped(fmt.Sprintf("%s: deploying app [run nvm use]...", sls.App.Name))
		err := sls.runNvm(&Writer{s})
		if err != nil {
			return fmt.Errorf("can't run nvm: %w", err)
		}
		s.Finish(true)

		s = sg.Scoped(fmt.Sprintf("%s: deploying app [run npm install]...", sls.App.Name))
		err = sls.runNpmInstall(&Writer{s})
		if err != nil {
			return fmt.Errorf("can't run npm install: %w", err)
		}
		s.Finish(true)

		if sls.App.CreateDomain {
			s = sg.Scoped(fmt.Sprintf("%s: deploying app [run serverless create_domain]...", sls.App.Name))
			err = sls.runCreateDomain(&Writer{s})
			if err != nil {
				return fmt.Errorf("can't run serverless create_domain: %w", err)
			}
			s.Finish(true)
		}

		s = sg.Scoped(fmt.Sprintf("%s: deploying app [run serverless deploy]...", sls.App.Name))
		err = sls.runDeploy(&Writer{s})
		if err != nil {
			return fmt.Errorf("can't run serverless deploy: %w", err)
		}
		s.Finish(true)
	case "docker":
		err := sls.deployWithDocker(sg)
		if err != nil {
			return err
		}
	}

	sg.Finish(true)

	return nil
}

func (sls *Manager) Destroy(ui *echelon.Logger) error {
	sls.prepare()

	sg := ui.Scoped(fmt.Sprintf("%s: destroying app...", sls.App.Name))

	switch sls.Project.PreferRuntime {
	case "native":
		s := sg.Scoped(fmt.Sprintf("%s: destroying app [run nvm use]...", sls.App.Name))

		err := sls.runNvm(&Writer{logger: s})
		if err != nil {
			return fmt.Errorf("can't run nvm: %w", err)
		}
		s.Finish(true)

		s = sg.Scoped(fmt.Sprintf("%s: destroying app [run npm install]...", sls.App.Name))
		err = sls.runNpmInstall(&Writer{logger: s})
		if err != nil {
			return fmt.Errorf("can't run npm install: %w", err)
		}
		s.Finish(true)

		if sls.App.CreateDomain {
			s = sg.Scoped(fmt.Sprintf("%s: destroying app [run serverless delete_domain]...", sls.App.Name))
			err = sls.runRemoveDomain(&Writer{logger: s})
			if err != nil {
				return fmt.Errorf("can't run serverless delete_domain: %w", err)
			}
			s.Finish(true)
		}

		s = sg.Scoped(fmt.Sprintf("%s: destroying app [run serverless remove]...", sls.App.Name))
		err = sls.runRemove(&Writer{logger: s})
		if err != nil {
			return fmt.Errorf("can't run serverless deploy: %w", err)
		}
		s.Finish(true)
	case "docker":
		err := sls.removeWithDocker(sg)
		if err != nil {
			return err
		}
	}

	sg.Finish(true)

	return nil
}

func (sls *Manager) Push(ui *echelon.Logger) error {
	return nil
}

func (sls *Manager) Build(ui *echelon.Logger) error {
	return nil
}

func (sls *Manager) Redeploy(ui *echelon.Logger) error {
	return nil
}
