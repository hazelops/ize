package config

import (
	"os"
	"path"
	"path/filepath"

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
