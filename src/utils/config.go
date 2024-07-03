// Copyright 2024 Bitcrush Testing

package utils

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var defaultAgentNames = []string{
	"Alois",
	"Anton",
	"August",
	"Ferdinand",
	"Franz",
	"Frederick",
	"Gustav",
	"Heinrich",
	"Hermann",
	"Joseph",
	"Karl",
	"Leopold",
	"Matthias",
	"Maximilian",
	"Otto",
	"Rudolf",
}

var viperInstance *viper.Viper

// Config represents the configuration structure
type Config struct {
	Agent struct {
		Name string `yaml:"name"`
		Id   string `yaml:"id"`
	} `yaml:"agent"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Key  string `yaml:"key"`
	} `yaml:"server"`
}

type SerialPortConfig struct {
	Port     string `yaml:"port"`
	BaudRate int    `yaml:"baud_rate"`
	DataBits int    `yaml:"data_bits"`
	Parity   string `yaml:"parity"`
	StopBits int    `yaml:"stop_bits"`
}

func NewConfig() *Config {

	id := uuid.New()
	viperInstance = viper.New()
	viperInstance.AddConfigPath(".")
	viperInstance.AddConfigPath("/etc/bon-voyage-agent")
	viperInstance.AddConfigPath("$HOME/.bon-voyage-agent")
	viperInstance.SetConfigName("config.yaml")
	viperInstance.SetConfigType("yaml")
	viperInstance.ReadInConfig()

	c := new(Config)
	c.Agent.Name = defaultAgentNames[rand.Intn(len(defaultAgentNames))]
	c.Agent.Id = id.String()
	c.Server.Host = "127.0.0.1"
	c.Server.Port = "6666"
	c.Server.Key = ""
	return c
}

func (c *Config) LoadConfig() error {

	err := viperInstance.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	err = viperInstance.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	return nil
}

func (c *Config) SaveConfig() error {

	return viperInstance.SafeWriteConfig()
}

func (c *Config) PluginConfig(name string, params *any) error {

	if !viperInstance.IsSet(name) {
		return fmt.Errorf("plugin %s configuration not found", name)
	}

	if err := viperInstance.Sub(name).Unmarshal(params); err != nil {
		return fmt.Errorf("error unmarshalling plugin config: %w", err)
	}
	return nil
}
