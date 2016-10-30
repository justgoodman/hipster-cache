package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	MetricsPort   int    `json:"metrics_port"`
	ServerPort    int    `json:"server_port"`
	Address       string `json:"address"`
	ConsulAddress string `json:"consul_address"`
	MaxBytesSize  int64  `json:"maximum_bytes_size"`
	MaxLenghtKey  int64  `json:"maximum_lenght_key"`
	InitCapacity  int64  `json:"init_capacity"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadFile(configPath string) error {
	configString, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configString, c)
	if err != nil {
		return err
	}

	if serverIP := os.Getenv("SERVER_IP"); serverIP != "" {
		c.Address = serverIP
	}
	if consulURL := os.Getenv("CONSUL_URL"); consulURL != "" {
		c.ConsulAddress = consulURL
	}
	fmt.Printf("%#v", c)
	return nil
}
