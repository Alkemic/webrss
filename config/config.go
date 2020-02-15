package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
}

func LoadConfig(configFile string) (*Config, error) {
	var config Config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	return &config, err
}
