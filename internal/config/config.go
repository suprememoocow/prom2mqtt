package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

// Query represents a single query in a group in the configuration.
type Query struct {
	Expr  string `yaml:"expr"`
	Topic string `yaml:"topic"`
}

// Group represents a set of queries which will be run at a fixed interval.
type Group struct {
	Name     string        `yaml:"name"`
	Interval time.Duration `yaml:"interval"`
	Queries  []*Query      `yaml:"queries"`
}

// Prom2MQTTConfig represents the full config file.
type Prom2MQTTConfig struct {
	Groups []*Group `yaml:"groups"`
}

// Load will load a config file from disk and return it as a struct.
func Load(filename string) (*Prom2MQTTConfig, error) {
	result := &Prom2MQTTConfig{}

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	err = yaml.Unmarshal(f, &result)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}

	return result, nil
}
