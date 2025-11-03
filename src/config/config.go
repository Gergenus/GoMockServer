package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoEndpoints        = errors.New("no endpoints found in config")
	ErrEmptyPath          = errors.New("no path found in endpoint")
	ErrDuplicateEndpoints = errors.New("duplicate endpoints")
	ErrXMLfileNotFound    = errors.New("xml file not found")
)

type Config struct {
	Port      string     `yaml:"port"`
	LogLevel  string     `yaml:"log_level"`
	LogFormat string     `yaml:"log_format"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Type         string `yaml:"type"`
	Method       string `yaml:"method"`
	Status       int    `yaml:"status"`
	Path         string `yaml:"path"`
	ResponsePath string `yaml:"response_path"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Endpoints) == 0 {
		return ErrNoEndpoints
	}

	endpointMap := make(map[string]bool)

	for _, ep := range c.Endpoints {
		if ep.Path == "" {
			return fmt.Errorf("%w: emptry path in endpoint", ErrEmptyPath)
		}

		pathMethod := ep.Path + ":" + ep.Method
		if _, ok := endpointMap[pathMethod]; ok {
			return fmt.Errorf("%w %s %s", ErrDuplicateEndpoints, ep.Path, ep.Method)
		}

		endpointMap[pathMethod] = true

		if ep.ResponsePath != "" {
			if _, err := os.Stat(ep.ResponsePath); os.IsNotExist(err) {
				return fmt.Errorf("%w %s for %s %s", ErrXMLfileNotFound, ep.ResponsePath, ep.Path, ep.Method)
			}
		}
	}
	return nil
}
