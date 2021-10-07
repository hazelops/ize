package config

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const Filename = "ize.hcl"

type Config struct {
	hclConfig
	pathData map[string]string
}

type hclConfig struct {
	AwsConfig        string `hcl:"aws_config,optional"`
	TerraformVersion string `hcl:"terraform_version"`
	Env              string `hcl:"env"`
	AwsRegion        string `hcl:"aws_region"`
	AwsProfile       string `hcl:"aws_profile"`
	Namespace        string `hcl:"namespace"`

	Service []*struct {
		Type string   `hcl:"type,label"`
		Name string   `hcl:"name,label"`
		Body hcl.Body `hcl:",remain"`
	} `hcl:"service,block"`
}

type Service struct {
	Use    *Use     `hcl:"use,block"`
	Body   hcl.Body `hcl:",body"`
	Remain hcl.Body `hcl:",remain"`
}

type Use struct {
	Type string   `hcl:",label"`
	Body hcl.Body `hcl:",remain"`
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
	var ctx *hcl.EvalContext

	// We require an absolute path for the path so we can set the path vars
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	// Decode
	var cfg hclConfig
	if err := hclsimple.DecodeFile(path, ctx, &cfg); err != nil {
		return nil, err
	}

	return &Config{
		hclConfig: cfg,
	}, nil
}

type Ecs struct {
	TerraformStateBucketName string `hcl:"terraform_state_bucket_name"`
}
