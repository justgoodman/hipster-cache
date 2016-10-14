package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	MetricsPort   int    `json:"metrics_port"`
	ServerPort    int    `json:"server_port"`
	Address       string `json:"address"`
	ConsulAddress string `json:"consul_address"`
}

func NewConfig() *Config {
	return &Config{}
}

func (this *Config) LoadFile(configPath string) error {
	configString, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configString, this)
	if err != nil {
		return err
	}
	fmt.Printf("%#v", this)
	return nil
}
