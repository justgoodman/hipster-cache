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
	MaxBytesSize  int64  `json:"maximum_bytes_size"`
	MaxLenghtKey  int64  `json:maximum_lenght_key`
	InitCapacity  int64  `json:init_capacity`
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
