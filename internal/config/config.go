package config

import "gopkg.in/ini.v1"

type Config struct {
	Logging LoggingConfig `ini:"logging"`
}

type LoggingConfig struct {
	Level string `ini:"level"`
	Size  int    `ini:"size"`
	Count int    `ini:"count"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	err := ini.MapTo(cfg, configPath)
	return cfg, err
}
