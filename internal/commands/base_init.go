package commands

import (
	"errors"
	"fmt"

	"github.com/hazelops/ize/internal/config"
)

func (c *izeBuilderCommon) initConfig(filename string) (*config.Config, error) {
	path, err := c.initConfigPath(filename)
	if err != nil {
		return nil, err
	}

	if path == "" {
		return nil, errors.New("an Ize configuration file (ize.hcl) is required but wasn't found")
	}

	return c.initConfigLoad(path)
}

func (c *izeBuilderCommon) initConfigPath(filename string) (string, error) {
	path, err := config.FindPath(filename)
	if err != nil {
		return "", fmt.Errorf("error looking for a Ize configuration: %s", err)
	}

	return path, nil
}

func (c *izeBuilderCommon) initConfigLoad(path string) (*config.Config, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return nil, err
	}

	//TODO: Validate

	return cfg, nil
}
