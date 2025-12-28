package config

import "gopkg.in/ini.v1"

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	err := ini.MapTo(cfg, configPath)
	return cfg, err
}
