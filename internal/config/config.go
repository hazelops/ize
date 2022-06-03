package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hazelops/ize/internal/aws/utils"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	Filename    = "ize.hcl"
	defaultPerm = 0665
)

type Config struct {
	AwsRegion       string `mapstructure:"aws_region"`
	AwsProfile      string `mapstructure:"aws_profile"`
	Namespace       string `mapstructure:"namespace"`
	Env             string `mapstructure:"env"`
	Session         *session.Session
	IsGlobal        bool
	IsDockerRuntime bool
	IsPlainText     bool
}

type requiments struct {
	configFile bool
	smplugin   bool
}

type Option func(*requiments)

func WithConfigFile() Option {
	return func(r *requiments) {
		r.configFile = true
	}
}

func WithSSMPlugin() Option {
	return func(r *requiments) {
		r.smplugin = true
	}
}

func CheckRequirements(options ...Option) error {
	r := requiments{}
	for _, opt := range options {
		opt(&r)
	}

	if r.smplugin {
		if err := checkSessionManagerPlugin(); err != nil {
			return err
		}
	}

	switch viper.GetString("prefer-runtime") {
	case "native":
		logrus.Debug("use native runtime")
	case "docker":
		if err := checkDocker(); err != nil {
			return err
		}
		logrus.Debug("use docker runtime")
	default:
		return fmt.Errorf("unknown runtime type: %s", viper.GetString("prefer-runtime"))
	}

	if len(viper.ConfigFileUsed()) == 0 && r.configFile {
		return fmt.Errorf("this command required config file")
	}

	return nil
}

func GetConfig() (*Config, error) {
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

	logrus.Debug("config file used:", viper.ConfigFileUsed())

	cfg := &Config{}

	viper.Unmarshal(&cfg)

	if len(cfg.AwsProfile) == 0 {
		return nil, fmt.Errorf("AWS profile must be specified using flags, env or config file")
	}

	if len(cfg.AwsRegion) == 0 {
		return nil, fmt.Errorf("AWS region must be specified using flags, env or config file")
	}

	if len(cfg.Namespace) == 0 {
		return nil, fmt.Errorf("namespace must be specified using flags, env or config file")
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  cfg.AwsRegion,
		Profile: cfg.AwsProfile,
	})
	if err != nil {
		return nil, err
	}

	cfg.Session = sess

	resp, err := sts.New(sess).GetCallerIdentity(
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return nil, err
	}

	viper.SetDefault("DOCKER_REGISTRY", fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, viper.GetString("aws_region")))
	// Reset env directory to default because env may change
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))

	plainText := viper.GetBool("plain-text") || viper.GetBool("plain_text")

	if plainText {
		pterm.DisableStyling()
	}

	cfg.IsPlainText = plainText

	if viper.GetString("prefer-runtime") == "docker" {
		cfg.IsDockerRuntime = true
	}

	return cfg, nil
}

func InitConfig() {
	viper.SetEnvPrefix("IZE")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetDefault("ENV", os.Getenv("ENV"))
	viper.SetDefault("AWS_PROFILE", os.Getenv("AWS_PROFILE"))
	viper.SetDefault("AWS_REGION", os.Getenv("AWS_REGION"))
	viper.SetDefault("NAMESPACE", os.Getenv("NAMESPACE"))
	// TODO: those static defaults should probably go to a separate package and/or function. Also would include image names and such.
	viper.SetDefault("TERRAFORM_VERSION", "1.1.3")

	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalln("can't initialize config: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("can't initialize config: %w", err)
	}

	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("PROJECTS_PATH", filepath.Join(cwd, "projects"))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	setDefaultInfraDir(cwd)

	cfg, err := readConfigFile(viper.GetString("config-file"))
	if err != nil {
		logrus.Fatal("can't initialize config: %w", err)
	}

	if cfg == nil {
		cfg, err = readGlobalConfigFile()
		if err != nil {
			logrus.Fatal("can't initialize config: %w", err)
		}
	}

	logrus.Debug("config file used:", viper.ConfigFileUsed())

	if cfg == nil {
		if err := viper.Unmarshal(&cfg); err != nil {
			logrus.Fatalln(err)
		}
	}

	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		if viper.GetString("ENV") == "" {
			pterm.Warning.Printfln("Can't read ENV, please set the value via --env flag or env variable")
		} else {
			viper.SetDefault("TAG", viper.GetString("ENV"))
			pterm.Warning.Printfln("Could not run git rev-parse, the default tag was set: %s", viper.GetString("TAG"))
		}
	} else {
		viper.SetDefault("TAG", strings.Trim(string(out), "\n"))
	}

	viper.SetDefault("TF_LOG", "")

	if cfg.IsGlobal {
		viper.SetDefault("ENV_DIR", fmt.Sprintf("%s/.ize/%s/%s", home, cfg.Namespace, cfg.Env))
		_, err := os.Stat(viper.GetString("ENV_DIR"))
		if os.IsNotExist(err) {
			os.MkdirAll(viper.GetString("ENV_DIR"), defaultPerm)
		}
	}

}

func setDefaultInfraDir(cwd string) {
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.ize", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.ize/env/%v", cwd, viper.GetString("ENV")))
	_, err := os.Stat(viper.GetString("INFRA_DIR"))
	if err != nil {
		viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
		viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.GetString("ENV")))
	}
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

func readGlobalConfigFile() (*Config, error) {
	env := viper.GetString("env")
	namespace := viper.GetString("namespace")

	if len(env) == 0 || len(namespace) == 0 {
		return nil, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(fmt.Sprintf("%s/.ize", home))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil
		} else {
			return nil, fmt.Errorf("can't read config file %s: %w", viper.ConfigFileUsed(), err)
		}
	}

	var cfg Config
	if viper.IsSet(fmt.Sprintf("%s.%s", namespace, env)) {
		if err := viper.UnmarshalKey(fmt.Sprintf("%s.%s", namespace, env), &cfg); err != nil {
			return nil, err
		}
	} else {
		logrus.Debug(fmt.Sprintf("config for %s.%s not found", namespace, env))
	}

	cfg.Env = env
	cfg.Namespace = namespace
	cfg.IsGlobal = true

	return &cfg, nil
}

func readConfigFile(path string) (*Config, error) {
	if len(path) != 0 {
		viper.SetConfigFile(path)
	} else {
		viper.SetConfigName("ize")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(viper.GetString("ENV_DIR"))
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, nil
		} else {
			return nil, fmt.Errorf("can't read config file %s: %w", viper.ConfigFileUsed(), err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
