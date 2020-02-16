package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const defaultPerPage = 50

type configRun struct {
	Host string
	Port int
}

type configDB struct {
	Verbose bool
	Time    bool

	User     string
	Password string
	Database string

	Host string
	Port string
}

type Config struct {
	DB  configDB
	Run configRun

	PerPage int64
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("cannot parse config file: %w", err)
	}

	if config.PerPage == 0 {
		config.PerPage = defaultPerPage
	}

	return &config, nil
}
