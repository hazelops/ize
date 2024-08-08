package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/hazelops/ize/internal/schema"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"

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
	AwsRegion  string `mapstructure:"aws_region"`
	AwsProfile string `mapstructure:"aws_profile"`
	Namespace  string `mapstructure:"namespace"`
	Env        string `mapstructure:"env"`
	LocalStack bool   `mapstructure:"localstack"`

	Session         *session.Session
	IsGlobal        bool
	IsDockerRuntime bool
	IsPlainText     bool
}

func (p *Project) GetConfig() error {

	err := MigrateAppsConfig()
	if err != nil {
		return err
	}

	err = MigrateInfraConfig()
	if err != nil {
		return err
	}

	err = MigrateTunnelConfig()
	if err != nil {
		return err
	}

	err = SetTag()
	if err != nil {
		return fmt.Errorf("can't set tag: %w. \nIs it a git repo?", err)
	}

	err = schema.Validate(viper.AllSettings())
	if err != nil {
		return err
	}

	err = viper.Unmarshal(p)
	if err != nil {
		return err
	}

	err = findDuplicates(p)
	if err != nil {
		return err
	}

	if p.LocalStack {
		// Set default Endpoint URL for localstack if it's enabled
		if len(p.EndpointUrl) == 0 {
			p.EndpointUrl = "http://127.0.0.1:4566"
		}
	}

	if len(p.SshPublicKey) == 0 {
		// Read id_rsa if it's not set
		home, _ := os.UserHomeDir()
		key, err := ioutil.ReadFile(fmt.Sprintf("%s/.ssh/id_rsa.pub", home))
		if err != nil {
			return fmt.Errorf("can't read public ssh key: %s", err)

		}

		p.SshPublicKey = string(key)[:len(string(key))-1]
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:      p.AwsRegion,
		Profile:     p.AwsProfile,
		EndpointUrl: p.EndpointUrl,
	})
	if err != nil {
		return err
	}

	p.Session = sess

	resp, err := sts.New(sess).GetCallerIdentity(
		&sts.GetCallerIdentityInput{},
	)
	if err != nil {
		return err
	}

	if len(p.DockerRegistry) == 0 {
		if p.LocalStack {
			p.DockerRegistry = fmt.Sprintf("%v.dkr.ecr.%v.localhost.localstack.cloud:4512", *resp.Account, p.AwsRegion)
		} else {
			p.DockerRegistry = fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, p.AwsRegion)
		}
		logrus.Debugf("Setting Docker Registry to %s", p.DockerRegistry)
	}
	// Reset env directory to default because env may change
	if len(p.DockerRegistry) == 0 {
		p.TFLogPath = fmt.Sprintf("%v/tflog.txt", p.EnvDir)
	}

	if viper.GetString("PREFER_RUNTIME") == "docker" {
		pterm.Warning.Println("Docker runtime is being deprecated. Please switch to native.")
	}

	if p.PlainText {
		pterm.DisableStyling()
	}

	return nil
}

func (p *Project) GetTestConfig() error {
	switch viper.GetString("log_level") {
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

	err := MigrateAppsConfig()
	if err != nil {
		return err
	}

	err = MigrateInfraConfig()
	if err != nil {
		return err
	}

	err = MigrateTunnelConfig()
	if err != nil {
		return err
	}

	SetTag()
	if err != nil {
		return fmt.Errorf("can't set tag: %w. \nIs it a git repo?", err)
	}

	err = schema.Validate(viper.AllSettings())
	if err != nil {
		return err
	}

	logrus.Debug("Config file used: ", viper.ConfigFileUsed())

	err = viper.Unmarshal(p)
	if err != nil {
		return err
	}

	sess, err := utils.GetTestSession(&utils.SessionConfig{
		Region:      p.AwsRegion,
		Profile:     p.AwsProfile,
		EndpointUrl: p.EndpointUrl,
	})
	if err != nil {
		return err
	}

	p.Session = sess

	if len(p.DockerRegistry) == 0 {
		p.DockerRegistry = fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", 0, p.AwsRegion)
	}
	// Reset env directory to default because env may change
	if len(p.DockerRegistry) == 0 {
		p.TFLogPath = fmt.Sprintf("%v/tflog.txt", p.EnvDir)
	}

	if viper.GetString("PREFER_RUNTIME") == "docker" {
		pterm.Warning.Println("Docker runtime is being deprecated. Please switch to native.")
	}

	if p.PlainText {
		pterm.DisableStyling()
	}

	return nil
}

func SetTag() error {
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
	return err
}

func InitConfig() {
	viper.SetEnvPrefix("IZE")

	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	// Set Log Level
	switch viper.GetString("log_level") {
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

	// Variables that would be read even without IZE_ prefix (in addition to IZE_ENV, IZE_TAG, etc)
	_ = viper.BindEnv("ENV", "ENV")
	_ = viper.BindEnv("TAG", "TAG")
	_ = viper.BindEnv("AWS_PROFILE", "AWS_PROFILE")
	_ = viper.BindEnv("AWS_REGION", "AWS_REGION")
	_ = viper.BindEnv("NAMESPACE", "NAMESPACE")
	_ = viper.BindEnv("SSH_PUBLIC_KEY", "SSH_PUBLIC_KEY")

	// TODO: those static defaults should probably go to a separate package and/or function. Also would include image names and such.
	viper.SetDefault("TERRAFORM_VERSION", "1.1.3")
	viper.SetDefault("NVM_VERSION", "0.39.7")
	viper.SetDefault("PREFER_RUNTIME", "native")
	viper.SetDefault("CUSTOM_PROMPT", false)
	viper.SetDefault("PLAIN_TEXT_OUTPUT", false)
	viper.SetDefault("LOCALSTACK", false)
	viper.SetDefault("apps_provider", "ecs")

	home, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatalln("Can't initialize config: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("Can't initialize config: %w", err)
	}

	// Set default apps folder - first option is `projects`
	appsPath := filepath.Join(cwd, "projects")
	_, err = os.Stat(appsPath)
	if os.IsNotExist(err) {
		// Second option `apps`
		appsPath = filepath.Join(cwd, "apps")

		// TODO: Add multi-repo support (cwd is the app directory).
		// Maybe use ../../cwd? so the repo directory would be the app name? and ize.toml would be in the root of the repo.
		// For now try setting `export APPS_PATH=../
	}

	logrus.Debugf("Setting APPS_PATH to %v", appsPath)
	viper.SetDefault("APPS_PATH", appsPath)

	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))

	setDefaultInfraDir(cwd)

	configFileLocation := viper.GetString("config_file")
	cfg, err := readConfigFile(configFileLocation)
	if err != nil {
		logrus.Fatal("can't initialize config: %w", err)
	}

	if cfg == nil {
		cfg, err = readGlobalConfigFile()
		if err != nil {
			logrus.Fatal("can't initialize config: %w", err)
		}
	}

	logrus.Debug("Config file used: ", viper.ConfigFileUsed())

	if cfg == nil {
		if err := viper.Unmarshal(&cfg); err != nil {
			logrus.Fatalln(err)
		}
	}

	viper.SetDefault("TF_LOG", "")

	// Global Config is an experimental feature where you can have one config in your home directory
	if cfg.IsGlobal {
		viper.SetDefault("ENV_DIR", fmt.Sprintf("%s/.ize/%s/%s", home, cfg.Namespace, cfg.Env))
		_, err := os.Stat(viper.GetString("ENV_DIR"))
		if os.IsNotExist(err) {
			err := os.MkdirAll(viper.GetString("ENV_DIR"), defaultPerm)
			if err != nil {
				logrus.Fatalln(err)
			}
		}
	}
}

func findDuplicates(cfg *Project) error {
	existingKeys := map[string]string{}
	duplicateKeys := map[string]map[string]string{}

	for k := range cfg.Terraform {
		if val, ok := existingKeys[k]; ok {
			if duplicateKeys[k] == nil {
				duplicateKeys[k] = map[string]string{}
			}
			duplicateKeys[k]["terraform"] = k
			if _, ok := duplicateKeys[k][val]; !ok {
				duplicateKeys[k][val] = k

			}
		}
		existingKeys[k] = "terraform"
	}

	for k := range cfg.Ecs {
		if val, ok := existingKeys[k]; ok {
			if duplicateKeys[k] == nil {
				duplicateKeys[k] = map[string]string{}
			}
			duplicateKeys[k]["ecs"] = k
			if _, ok := duplicateKeys[k][val]; !ok {
				duplicateKeys[k][val] = k

			}
		}
		existingKeys[k] = "ecs"
	}

	for k := range cfg.Serverless {
		if val, ok := existingKeys[k]; ok {
			if duplicateKeys[k] == nil {
				duplicateKeys[k] = map[string]string{}
			}
			duplicateKeys[k]["serverless"] = k
			if _, ok := duplicateKeys[k][val]; !ok {
				duplicateKeys[k][val] = k

			}
		}
		existingKeys[k] = "serverless"
	}

	for k := range cfg.Alias {
		if val, ok := existingKeys[k]; ok {
			if duplicateKeys[k] == nil {
				duplicateKeys[k] = map[string]string{}
			}
			duplicateKeys[k]["alias"] = k
			if _, ok := duplicateKeys[k][val]; !ok {
				duplicateKeys[k][val] = k

			}
		}
		existingKeys[k] = "alias"
	}

	errMsg := ""
	if len(duplicateKeys) != 0 {
		for name, v := range duplicateKeys {
			errMsg += fmt.Sprintf("\nOnly one section with the name \"%s\" is allowed. Please rename one of the following:\n", name)
			for k, v := range v {
				errMsg += fmt.Sprintf("- [%s.%s]\n", k, v)
			}
		}
	}

	if len(errMsg) != 0 {
		return fmt.Errorf(errMsg)
	}

	return nil
}

func setDefaultInfraDir(cwd string) {
	izeDir := fmt.Sprintf("%v/.ize", cwd)
	envDir := fmt.Sprintf("%v/env/%v", izeDir, viper.GetString("ENV"))

	_, err := os.Stat(izeDir) // Check if directory that we've set exists
	if err != nil {
		// izeDir doesn't exist, so setting the default to .infra
		logrus.Debugf("Tried %v, but not found.", izeDir)

		izeDir = fmt.Sprintf("%v/.infra", cwd)
		envDir = fmt.Sprintf("%v/env/%v", izeDir, viper.GetString("ENV"))

		_, err := os.Stat(izeDir) // Check if directory that we've set exists
		if err != nil {
			// izeDir doesn't exist, so setting the default to .infra
			logrus.Debugf("Tried %v, but not found.", izeDir)

			izeDir = fmt.Sprintf("%v", cwd)
			envDir = fmt.Sprintf("%v/env/%v", izeDir, viper.GetString("ENV"))
		}
	}

	logrus.Debug("Setting IZE_DIR to ", izeDir)
	viper.SetDefault("IZE_DIR", izeDir)

	logrus.Debug("Setting ENV_DIR to ", envDir)
	viper.SetDefault("ENV_DIR", envDir)

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
		// If path is defined use it to read config
		logrus.Debug("Reading config file:", path)
		viper.SetConfigFile(path)

	} else {
		// If path is undefined read using viper's ConfigPath
		logrus.Debug("Config file is not overriden via `config_file`")

		viper.SetConfigName("ize")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		envDir := filepath.Join(viper.GetString("ENV_DIR"))

		logrus.Debug(fmt.Sprintf("Adding config path to viper: %s", envDir))
		viper.AddConfigPath(envDir)

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

func MigrateAppsConfig() error {
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

func MigrateInfraConfig() error {
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

func MigrateTunnelConfig() error {
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

// GetApps returns a list of application names in the directory for shell completions
func GetApps(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	var apps []string

	dir, _ := os.ReadDir("./apps")

	if dir != nil {
		for _, entry := range dir {
			apps = append(apps, entry.Name())
		}
	}

	dir, _ = os.ReadDir("./projects")

	if dir != nil {
		for _, entry := range dir {
			apps = append(apps, entry.Name())
		}
	}

	return apps, cobra.ShellCompDirectiveNoFileComp
}

func (p *Project) Generate(tmpl string, funcs template.FuncMap) error {
	t := template.New("template")
	t.Funcs(funcs)
	t, err := t.Parse(tmpl)
	if err != nil {
		return err
	}

	data := struct {
		Project
		Data interface{}
	}{
		*p,
		tmpl,
	}

	err = t.Execute(os.Stdout, data)
	if err != nil {
		return err
	}

	return nil
}
