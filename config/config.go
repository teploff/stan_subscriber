package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config main application config.
type Config struct {
	Stan StanConfig `yaml:"stan"`
}

// StanConfig configuration for the stan (nats-streaming).
//
// ClusterID - represented conn to stan. It cans contain only alphanumeric and `-` or `_` characters.
//
// ClientID - unique client identificator.
//
// Addr - Bind to host address
type StanConfig struct {
	ClusterID string `yaml:"cluster_id"`
	ClientID  string `yaml:"client_id"`
	Addr      string `yaml:"addr"`
	Subject   string `yaml:"subject"`
}

// LoadFromFile create configuration from file.
func LoadFromFile(fileName string) (Config, error) {
	cfg := Config{}

	configBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
