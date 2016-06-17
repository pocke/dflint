package main

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	IgnoreRules []string `yaml:"ignore_rules"`
}

func ParseConfig() (*Config, error) {
	path, err := homedir.Expand("~/.config/dflint.yaml")
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return new(Config), nil
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
