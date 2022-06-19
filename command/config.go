package command

import (
	"fmt"
	"io/ioutil"

	yml "gopkg.in/yaml.v3"
)

type SpinalConfig struct {
	LogDir string `yaml:"log_dir,omitempty" json:"log_dir,omitempty"`
	DbDir  string `yaml:"db_dir,omitempty"  json:"db_dir,omitempty"`
	// DbType string `yaml:"db_type,omitempty"  json:"db_type,omitempty"`
}

func LoadConfig(fname string) SpinalConfig {
	cfg := SpinalConfig{}

	text, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println("Cannot read config file!")
		return cfg
	}
	if err := yml.Unmarshal([]byte(text), &cfg); err != nil {
		fmt.Println("Cannot parse YAML config!")
		return cfg
	}

	return cfg
}
