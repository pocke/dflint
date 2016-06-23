package main

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	IgnoreRules []string `yaml:"ignore_rules"`
}

// ParseConfig parses config file.
// When file doesn't exists, the func returns an empty config(NOT return error).
func ParseConfig(path string) (*Config, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}

	if !fileExists(path) {
		return new(Config), nil
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := new(Config)
	yaml.Unmarshal(b, c)
	return c, nil
}

func (c *Config) IsEnabledRule(r string) bool {
	for _, cr := range c.IgnoreRules {
		if cr == r {
			return false
		}
	}
	return true
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
