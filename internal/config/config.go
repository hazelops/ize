package config

import (
	"fmt"
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
	AwsRegion       string `mapstructure:"aws_region"`
	AwsProfile      string `mapstructure:"aws_profile"`
	Namespace       string `mapstructure:"namespace"`
	Env             string `mapstructure:"env"`
	Session         *session.Session
	IsGlobal        bool
	IsDockerRuntime bool
	IsPlainText     bool
}

func (p *Project) GetConfig() error {
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

	err := ConvertApps()
	if err != nil {
		return err
	}

	err = ConvertInfra()
	if err != nil {
		return err
	}

	err = ConvertTunnel()
	if err != nil {
		return err
	}

	SetTag()
	if err != nil {
		return fmt.Errorf("can't set tag: %w", err)
	}

	err = schema.Validate(viper.AllSettings())
	if err != nil {
		return err
	}

	logrus.Debug("config file used:", viper.ConfigFileUsed())

	err = viper.Unmarshal(p)
	if err != nil {
		return err
	}

	err = findDuplicates(p)
	if err != nil {
		return err
	}

	sess, err := utils.GetSession(&utils.SessionConfig{
		Region:  p.AwsRegion,
		Profile: p.AwsProfile,
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
		p.DockerRegistry = fmt.Sprintf("%v.dkr.ecr.%v.amazonaws.com", *resp.Account, p.AwsRegion)
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

	err := ConvertApps()
	if err != nil {
		return err
	}

	err = ConvertInfra()
	if err != nil {
		return err
	}

	err = ConvertTunnel()
	if err != nil {
		return err
	}

	SetTag()
	if err != nil {
		return fmt.Errorf("can't set tag: %w", err)
	}

	err = schema.Validate(viper.AllSettings())
	if err != nil {
		return err
	}

	logrus.Debug("config file used:", viper.ConfigFileUsed())

	err = viper.Unmarshal(p)
	if err != nil {
		return err
	}

	sess, err := utils.GetTestSession(&utils.SessionConfig{
		Region:  p.AwsRegion,
		Profile: p.AwsProfile,
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

	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	_ = viper.BindEnv("ENV", "ENV")
	_ = viper.BindEnv("TAG", "TAG")
	_ = viper.BindEnv("AWS_PROFILE", "AWS_PROFILE")
	_ = viper.BindEnv("AWS_REGION", "AWS_REGION")
	_ = viper.BindEnv("NAMESPACE", "NAMESPACE")

	// TODO: those static defaults should probably go to a separate package and/or function. Also would include image names and such.
	viper.SetDefault("TERRAFORM_VERSION", "1.1.3")
	viper.SetDefault("PREFER_RUNTIME", "native")
	viper.SetDefault("CUSTOM_PROMPT", false)
	viper.SetDefault("PLAIN_TEXT_OUTPUT", false)

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
	viper.SetDefault("IZE_DIR", fmt.Sprintf("%v/.ize", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.ize/env/%v", cwd, viper.GetString("ENV")))
	_, err := os.Stat(viper.GetString("IZE_DIR"))
	if err != nil {
		viper.SetDefault("IZE_DIR", fmt.Sprintf("%v/.infra", cwd))
		viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.GetString("ENV")))
	}
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

func (p *Project) Generate(tmpl string) error {
	t := template.New("template")
	t, err := t.Parse(tmpl)
	if err != nil {
		return err
	}

	err = t.Execute(os.Stdout, p)
	if err != nil {
		return err
	}

	return nil
}
