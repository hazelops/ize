package requirements

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type requirements struct {
	configFile bool
	ssmplugin  bool
	structure  bool
	nvm        bool
}

func CheckRequirements(options ...Option) error {
	r := requirements{}
	for _, opt := range options {
		opt(&r)
	}

	if r.nvm {
		err := checkNVM()
		if err != nil {
			return err
		}
	}

	if r.structure {
		if !isStructured() {
			pterm.Warning.Println("is not an ize-structured directory. Please run ize init or cd into an ize-structured directory.")
		}
	}

	if r.ssmplugin {
		if err := checkSessionManagerPlugin(); err != nil {
			return err
		}
	}

	switch viper.GetString("prefer_runtime") {
	case "native":
		logrus.Debug("Using native runtime")
	case "docker":
		if err := checkDocker(); err != nil {
			return err
		}
		logrus.Debug("Using docker runtime (deprecated)")
	default:
		return fmt.Errorf("unknown runtime type: %s", viper.GetString("prefer_runtime"))
	}

	if len(viper.ConfigFileUsed()) == 0 && r.configFile {
		return fmt.Errorf("this command requires a config file. Please add ize.toml to %s", viper.GetString("env_dir"))
	}

	return nil
}

type Option func(*requirements)

func WithIzeStructure() Option {
	return func(r *requirements) {
		r.structure = true
	}
}

func WithConfigFile() Option {
	return func(r *requirements) {
		r.configFile = true
	}
}

func WithSSMPlugin() Option {
	return func(r *requirements) {
		r.ssmplugin = true
	}
}

func WithNVM() Option {
	return func(r *requirements) {
		r.nvm = true
	}
}

func checkNVM() error {
	if len(os.Getenv("NVM_DIR")) == 0 {
		return errors.New("nvm is not installed (visit https://github.com/nvm-sh/nvm)")
	}

	return nil
}

func checkDocker() error {
	exist, _ := CheckCommand("docker", []string{"info"})
	if !exist {
		return errors.New("docker is not running or is not installed (visit https://www.docker.com/get-started)")
	}

	return nil
}

func isStructured() bool {
	var isStructured = false

	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("can't initialize config: %w", err)
	}

	_, err = os.Stat(filepath.Join(cwd, ".ize"))
	if !os.IsNotExist(err) {
		isStructured = true
	}

	_, err = os.Stat(filepath.Join(cwd, ".infra"))
	if !os.IsNotExist(err) {
		isStructured = true
	}

	return isStructured
}

func checkSessionManagerPlugin() error {
	exist, _ := CheckCommand("session-manager-plugin", []string{})
	if !exist {
		pterm.Warning.Println("SSM Agent plugin is not installed. Trying to install SSM Agent plugin")

		var pyVersion string

		exist, pyVersion := CheckCommand("python3", []string{"--version"})
		if !exist {
			exist, pyVersion = CheckCommand("python", []string{"--version"})
			if !exist {
				return errors.New("python is not installed")
			}

			c, err := semver.NewConstraint("<= 2.6.5")
			if err != nil {
				return err
			}

			v, err := semver.NewVersion(strings.TrimSpace(strings.Split(pyVersion, " ")[1]))
			if err != nil {
				return err
			}

			if c.Check(v) {
				return fmt.Errorf("python version %s below required %s", v.String(), "2.6.5")
			}
			return errors.New("python is not installed")
		}

		c, err := semver.NewConstraint("<= 3.3.0")
		if err != nil {
			return err
		}

		v, err := semver.NewVersion(strings.TrimSpace(strings.Split(pyVersion, " ")[1]))
		if err != nil {
			return err
		}

		if c.Check(v) {
			return fmt.Errorf("python version %s below required %s", v.String(), "3.3.0")
		}

		pterm.DefaultSection.Println("Installing SSM Agent plugin")

		err = downloadSSMAgentPlugin()
		if err != nil {
			return fmt.Errorf("download SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Downloading SSM Agent plugin")

		err = installSSMAgent()
		if err != nil {
			return fmt.Errorf("install SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Installing SSM Agent plugin")

		err = cleanupSSMAgent()
		if err != nil {
			return fmt.Errorf("cleanup SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Cleanup Session Manager plugin installation package")

		exist, _ = CheckCommand("session-manager-plugin", []string{})
		if !exist {
			return fmt.Errorf("check SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}
	}

	return nil
}
