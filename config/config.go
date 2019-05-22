package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config is a struct containing the current configuration
type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Epic     string `yaml:"epic"`
	Assignee string `yaml:"assignee"`
	Status   string `yaml:"status"`
}

// Read the configuration from disk
func Read() Config {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Failed to read config from disk: #%v ", err)
	}

	conf := Config{}
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return conf
}

// Write the configuration to disk
func (conf Config) Write() {
	yamlFile, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	ioutil.WriteFile("config.yaml", yamlFile, 0644)
}
