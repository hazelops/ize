package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const Filename = "ize.hcl"

type Config struct {
	hclConfig
}

type hclConfig struct {
	TerraformVersion string                            `mapstructure:"terraform_version"`
	Env              string                            `mapstructure:"env"`
	AwsRegion        string                            `mapstructure:"aws_region"`
	AwsProfile       string                            `mapstructure:"aws_profile"`
	Namespace        string                            `mapstructure:"namespace"`
	Infra            map[string]map[string]interface{} `mapstructure:"infra"`
	InfraDir         string                            `mapstructure:"infra_dir"`
	RootDir          string                            `mapstructure:"root_dir"`
	Tag              string                            `mapstructure:"tag"`
}

func FindPath(filename string) (string, error) {
	if !path.IsAbs(filename) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		if filename == "" {
			filename = Filename
		}

		filename = filepath.Join(wd, filename)
	}

	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	} else {
		return "", err
	}
}

func Load(path string) (*Config, error) {

	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	//Decode
	var cfg hclConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &Config{
		hclConfig: cfg,
	}, nil
}

func InitializeConfig() error {
	viper.SetEnvPrefix("IZE")
	viper.AutomaticEnv()

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		PadLevelText:     true,
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	switch viper.GetString("log-level") {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(0)
	}

	if err := CheckRequirements(); err != nil {
		return err
	}

	if viper.GetString("config-file") != "" {
		_, err := initConfig(viper.GetString("config-file"))
		if err != nil {
			return err
		}
	}

	if len(viper.GetString("aws-profile")) == 0 {
		viper.Set("aws-profile", viper.GetString("aws_profile"))
		if len(viper.GetString("aws-profile")) == 0 {
			return fmt.Errorf("AWS profile must be specified using flags or config file")
		}
	}

	if len(viper.GetString("aws-region")) == 0 {
		viper.Set("aws-region", viper.GetString("aws_region"))
		if len(viper.GetString("aws-region")) == 0 {
			return fmt.Errorf("AWS region must be specified using flags or config file")
		}
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  viper.GetString("aws-region"),
		Profile: viper.GetString("aws-profile"),
	})
	if err != nil {
		return err
	}

	resp, err := sts.New(sess).GetCallerIdentity(
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AutomaticEnv()

	//TODO ensure values of the variables are checked for nil before passing down to docker.

	tag := ""
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		tag = viper.GetString("ENV")
		pterm.Warning.Printfln("could not run git rev-parse, the default tag was set: %s", tag)
	} else {
		tag = string(out)
	}

	viper.SetDefault("DOCKER_REGISTRY", fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, viper.GetString("aws-region")))
	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.GetString("ENV")))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))
	viper.SetDefault("TAG", string(tag))

	return nil
}

func CheckRequirements() error {
	//Check Docker and SSM Agent
	_, err := CheckCommand("docker", []string{"info"})
	if err != nil {
		return errors.New("docker is not running or is not installed (visit https://www.docker.com/get-started)")
	}

	_, err = CheckCommand("session-manager-plugin", []string{})
	if err != nil {
		pterm.Warning.Println("SSM Agent plugin is not installed. Trying to install SSM Agent plugin")

		var pyVersion string

		pyVersion, err = CheckCommand("python3", []string{"--version"})
		if err != nil {
			pyVersion, err = CheckCommand("python", []string{"--version"})
			if err != nil {
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

		err = DownloadSSMAgentPlugin()
		if err != nil {
			return fmt.Errorf("download SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Downloading SSM Agent plugin")

		err = InstallSSMAgent()
		if err != nil {
			return fmt.Errorf("install SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Installing SSM Agent plugin")

		err = CleanupSSMAgent()
		if err != nil {
			return fmt.Errorf("cleanup SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}

		pterm.Success.Println("Cleanup Session Manager plugin installation package")

		_, err = CheckCommand("session-manager-plugin", []string{})
		if err != nil {
			return fmt.Errorf("check SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}
	}

	return nil
}

func initConfig(filename string) (*Config, error) {
	path, err := initConfigPath(filename)
	if err != nil {
		return nil, err
	}

	if path == "" {
		return nil, errors.New("an Ize configuration file (ize.hcl) is required but wasn't found")
	}

	return initConfigLoad(path)
}

func initConfigPath(filename string) (string, error) {
	path, err := FindPath(filename)
	if err != nil {
		return "", fmt.Errorf("error looking for a Ize configuration: %s", err)
	}

	return path, nil
}

func initConfigLoad(path string) (*Config, error) {
	cfg, err := Load(path)
	if err != nil {
		return nil, err
	}

	//TODO: Validate

	return cfg, nil
}
