package config

import (
	"errors"
	"fmt"
	"github.com/hazelops/ize/internal/schema"
	"github.com/mitchellh/mapstructure"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
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
	structure  bool
	nvm        bool
}

type Option func(*requiments)

func WithIzeStructure() Option {
	return func(r *requiments) {
		r.structure = true
	}
}

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

func WithNVM() Option {
	return func(r *requiments) {
		r.nvm = true
	}
}

func CheckRequirements(options ...Option) error {
	r := requiments{}
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

	if r.smplugin {
		if err := checkSessionManagerPlugin(); err != nil {
			return err
		}
	}

	switch viper.GetString("prefer_runtime") {
	case "native":
		logrus.Debug("use native runtime")
	case "docker":
		if err := checkDocker(); err != nil {
			return err
		}
		logrus.Debug("use docker runtime")
	default:
		return fmt.Errorf("unknown runtime type: %s", viper.GetString("prefer_runtime"))
	}

	if len(viper.ConfigFileUsed()) == 0 && r.configFile {
		return fmt.Errorf("this command required config file")
	}

	return nil
}

func checkNVM() error {
	if len(os.Getenv("NVM_DIR")) == 0 {
		return errors.New("nvm is not installed (visit https://github.com/nvm-sh/nvm)")
	}

	return nil
}

func GetConfig() (*Project, error) {
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

	err := ConvertApps()
	if err != nil {
		return nil, err
	}

	err = ConvertInfra()
	if err != nil {
		return nil, err
	}

	err = ConvertTunnel()
	if err != nil {
		return nil, err
	}

	SetTag()
	if err != nil {
		return nil, fmt.Errorf("can't set tag: %w", err)
	}

	err = schema.Validate(viper.AllSettings())
	if err != nil {
		return nil, err
	}

	logrus.Debug("config file used:", viper.ConfigFileUsed())

	cfg := &Project{}

	viper.Unmarshal(&cfg)

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

	if len(cfg.DockerRegistry) == 0 {
		cfg.DockerRegistry = fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, cfg.AwsRegion)
	}
	// Reset env directory to default because env may change
	if len(cfg.DockerRegistry) == 0 {
		cfg.TFLogPath = fmt.Sprintf("%v/tflog.txt", cfg.EnvDir)
	}

	if viper.GetString("PREFER_RUNTIME") == "docker" {
		pterm.Warning.Println("Docker runtime is being deprecated. Please switch to native")
	}

	if cfg.PlainText {
		pterm.DisableStyling()
	}

	return cfg, nil
}

func SetTag() {
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
}

func InitConfig() {
	viper.SetEnvPrefix("IZE")
	viper.AutomaticEnv()

	viper.BindEnv("ENV", "ENV")
	viper.BindEnv("TAG", "TAG")
	viper.BindEnv("AWS_PROFILE", "AWS_PROFILE")
	viper.BindEnv("AWS_REGION", "AWS_REGION")
	viper.BindEnv("NAMESPACE", "NAMESPACE")

	// TODO: those static defaults should probably go to a separate package and/or function. Also would include image names and such.
	viper.SetDefault("TERRAFORM_VERSION", "1.1.3")
	viper.SetDefault("PREFER_RUNTIME", "native")
	viper.SetDefault("CUSTOM_PROMPT", false)
	viper.SetDefault("PLAIN_TEXT", false)

	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalln("can't initialize config: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("can't initialize config: %w", err)
	}

	// set default apps folder
	_, err = os.Stat(filepath.Join(cwd, "projects"))
	if os.IsNotExist(err) {
		viper.SetDefault("APPS_PATH", filepath.Join(cwd, "apps"))
	} else {
		viper.SetDefault("APPS_PATH", filepath.Join(cwd, "projects"))
	}

	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	setDefaultInfraDir(cwd)

	cfg, err := readConfigFile(viper.GetString("config_file"))
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
	viper.SetDefault("IZE_DIR", fmt.Sprintf("%v/.ize", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.ize/env/%v", cwd, viper.GetString("ENV")))
	_, err := os.Stat(viper.GetString("IZE_DIR"))
	if err != nil {
		viper.SetDefault("IZE_DIR", fmt.Sprintf("%v/.infra", cwd))
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

func ConvertApps() error {
	ecs := map[string]interface{}{}
	serverless := map[string]interface{}{}

	apps := viper.GetStringMap("app")
	for name, app := range apps {
		body := app.(map[string]interface{})
		t := body["type"].(string)
		switch t {
		case "ecs":
			ecsApp := Ecs{}
			err := mapstructure.Decode(&body, &ecsApp)
			if err != nil {
				return err
			}

			ecs[name] = structToMap(ecsApp)
		case "serverless":
			slsApp := Serverless{}
			err := mapstructure.Decode(&body, &slsApp)
			if err != nil {
				return err
			}

			serverless[name] = structToMap(slsApp)
		default:
			return fmt.Errorf("does not support %s type", t)
		}

	}

	err := viper.MergeConfigMap(map[string]interface{}{
		"ecs":        ecs,
		"serverless": serverless,
	})
	if err != nil {
		return err
	}

	return nil
}

func ConvertInfra() error {
	tf := viper.GetStringMap("infra.terraform")

	version, ok := tf["terraform_version"]
	if ok {
		delete(tf, "terraform_version")
		tf["version"] = version
	}

	err := viper.MergeConfigMap(map[string]interface{}{
		"terraform.infra": tf,
	})
	if err != nil {
		return err
	}

	return nil
}

func ConvertTunnel() error {
	tunnel := viper.GetStringMap("infra.tunnel")

	err := viper.MergeConfigMap(map[string]interface{}{
		"tunnel": tunnel,
	})
	if err != nil {
		return err
	}

	return nil
}

func structToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}

	v := reflect.ValueOf(item)
	typeOfOpts := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			key := strings.Split(typeOfOpts.Field(i).Tag.Get("mapstructure"), ",")[0]
			if len(key) == 0 {
				key = strings.ToLower(typeOfOpts.Field(i).Name)
			}
			res[key] = v.Field(i).Interface()
		}
	}

	return res
}
