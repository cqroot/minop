package manager

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ModuleConfig struct {
	ModuleArgs []map[string]string `yaml:"modules"`
}

func ModulesFromYaml(filename string) (*ModuleConfig, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	moduleArgs := ModuleConfig{}
	err = yaml.Unmarshal(content, &moduleArgs)
	if err != nil {
		return nil, err
	}

	return &moduleArgs, nil
}
