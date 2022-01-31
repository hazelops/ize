package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Filename = "ize.hcl"

type Config struct {
	AwsRegion  string `mapstructure:"aws_region"`
	AwsProfile string `mapstructure:"aws_profile"`
	Namespace  string `mapstructure:"namespace"`
	Env        string `mapstructure:"env"`
}

func readConfigFile(path string, required bool) (*Config, error) {
	viper.SetConfigName("ize")
	viper.SetConfigType("toml")
	viper.AddConfigPath(strings.TrimRight(path, "/ize.toml"))
	viper.AddConfigPath(".")
	viper.AddConfigPath(viper.GetString("ENV_DIR"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if required {
				return nil, fmt.Errorf("this command required config file")
			}
			logrus.Warn("config file not found")
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type requiments struct {
	configFile bool
	docker     bool
	smplugin   bool
}

type Option func(*requiments)

func WithConfigFile() Option {
	return func(r *requiments) {
		r.configFile = true
	}
}

func WithDocker() Option {
	return func(r *requiments) {
		r.docker = true
	}
}

func WithSSMPlugin() Option {
	return func(r *requiments) {
		r.smplugin = true
	}
}

func InitializeConfig(options ...Option) (*Config, error) {
	cfg := &Config{}
	var err error

	r := requiments{}
	for _, opt := range options {
		opt(&r)
	}

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
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.FatalLevel)
	}

	if r.smplugin {
		if err := checkSessionManagerPlugin(); err != nil {
			return nil, err
		}
	}

	if r.docker {
		if err := checkDocker(); err != nil {
			return nil, err
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("can't initialize config: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't initialize config: %w", err)
	}

	viper.SetDefault("ENV", os.Getenv("ENV"))
	viper.SetDefault("AWS_PROFILE", os.Getenv("AWS_PROFILE"))
	viper.SetDefault("AWS_REGION", os.Getenv("AWS_REGION"))
	viper.SetDefault("NAMESPACE", os.Getenv("NAMESPACE"))
	// Default paths
	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.GetString("ENV")))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))

	// TODO: those static defaults should probably go to a separate package and/or function. Also would include image names and such.
	viper.SetDefault("TERRAFORM_VERSION", "1.1.3")

	cfg, err = readConfigFile(viper.GetString("config-file"), r.configFile)
	if err != nil {
		return nil, fmt.Errorf("can't initialize config: %w", err)
	}

	if len(cfg.AwsProfile) == 0 {
		return nil, fmt.Errorf("AWS profile must be specified using flags or config file")
	}

	if len(cfg.AwsRegion) == 0 {
		return nil, fmt.Errorf("AWS region must be specified using flags or config file")
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  viper.GetString("aws-region"),
		Profile: viper.GetString("aws-profile"),
	})
	if err != nil {
		return nil, err
	}

	resp, err := sts.New(sess).GetCallerIdentity(
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return nil, err
	}

	tag := ""
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		tag = viper.GetString("ENV")
		pterm.Warning.Printfln("could not run git rev-parse, the default tag was set: %s", tag)
	} else {
		tag = string(out)
	}

	viper.SetDefault("DOCKER_REGISTRY", fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, viper.GetString("aws-region")))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))
	viper.SetDefault("TAG", string(tag))
	// Reset env directory to default because env may change
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.GetString("ENV")))

	return cfg, nil
}

func checkDocker() error {
	_, err := CheckCommand("docker", []string{"info"})
	if err != nil {
		return errors.New("docker is not running or is not installed (visit https://www.docker.com/get-started)")
	}

	return nil
}

func checkSessionManagerPlugin() error {
	_, err := CheckCommand("session-manager-plugin", []string{})
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

		_, err = CheckCommand("session-manager-plugin", []string{})
		if err != nil {
			return fmt.Errorf("check SSM Agent plugin error: %v (visit https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)", err)
		}
	}

	return nil
}
