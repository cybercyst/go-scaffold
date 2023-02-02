package template

import (
	"encoding/json"
	"errors"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type TemplateConfig struct {
	Title       string      `json:"title" yaml:"title"`
	Version     string      `json:"version" yaml:"version"`
	Description string      `json:"description" yaml:"description"`
	Tags        []string    `json:"tags" yaml:"tags"`
	Icon        string      `json:"icon" yaml:"icon"`
	RawSchema   interface{} `json:"schema" yaml:"schema"`
}

type ConfigType int8

const (
	JSON ConfigType = iota
	YAML
)

func readConfigFile(fs afero.Fs, fname string, configType ConfigType) (*TemplateConfig, error) {
	var config TemplateConfig

	bytes, err := afero.ReadFile(fs, fname)
	if err != nil {
		return nil, err
	}

	switch configType {
	case JSON:
		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return nil, err
		}
	case YAML:
		err = yaml.Unmarshal(bytes, &config)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func loadConfig(fs afero.Fs) (*TemplateConfig, error) {
	exists, _ := afero.Exists(fs, "template.yaml")
	if exists {
		return readConfigFile(fs, "template.yaml", YAML)
	}

	exists, _ = afero.Exists(fs, "template.yml")
	if exists {
		return readConfigFile(fs, "template.yml", YAML)
	}

	exists, _ = afero.Exists(fs, "template.json")
	if exists {
		return readConfigFile(fs, "template.json", JSON)
	}

	return nil, errors.New("no valid template definition found")
}

func LoadConfig(fs afero.Fs) (*TemplateConfig, error) {
	config, err := loadConfig(fs)
	if err != nil {
		return nil, err
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (config *TemplateConfig) Validate() error {
	if config.Title == "" {
		return errors.New("required property title missing from template config")
	}

	if config.Version == "" {
		return errors.New("required property version missing from template config")
	}

	return nil
}
