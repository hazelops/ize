package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const Filename = "ize.hcl"

type Config struct {
	hclConfig
}

type hclConfig struct {
	TerraformVersion string                                  `mapstructure:"terraform_version"`
	Env              string                                  `mapstructure:"env"`
	AwsRegion        string                                  `mapstructure:"aws_region"`
	AwsProfile       string                                  `mapstructure:"aws_profile"`
	Namespace        string                                  `mapstructure:"namespace"`
	Infra            map[string]map[string]map[string]string `mapstructure:"infra,block"`
}

func FindPath(filename string) (string, error) {
	var err error
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if filename == "" {
		filename = Filename
	}

	path := filepath.Join(wd, filename)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	} else {
		return "", err
	}
}

func Load(path string) (*Config, error) {

	// We require an absolute path for the path so we can set the path vars
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

	viper.SetEnvPrefix("ize")
	viper.AutomaticEnv()

	//Decode
	var cfg hclConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &Config{
		hclConfig: cfg,
	}, nil
}
