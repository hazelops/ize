package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

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

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory")
	}

	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AutomaticEnv()

	//TODO ensure values of the variables are checked for nil before passing down to docker.

	viper.SetDefault("ROOT_DIR", cwd)
	viper.SetDefault("INFRA_DIR", fmt.Sprintf("%v/.infra", cwd))
	viper.SetDefault("ENV_DIR", fmt.Sprintf("%v/.infra/env/%v", cwd, viper.GetString("ENV")))
	viper.SetDefault("HOME", fmt.Sprintf("%v", home))
	viper.SetDefault("TF_LOG", fmt.Sprintf(""))
	viper.SetDefault("TF_LOG_PATH", fmt.Sprintf("%v/tflog.txt", viper.Get("ENV_DIR")))

	//Decode
	var cfg hclConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &Config{
		hclConfig: cfg,
	}, nil
}
