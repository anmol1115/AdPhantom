package config

import (
	"gopkg.in/ini.v1"
)

type Config struct {
	Logging LoggingConfig `ini:"logging"`
	DNS     Dns           `ini:"dns"`
	Filter  Filter        `ini:"filter"`
}

type LoggingConfig struct {
	Level string `ini:"level"`
	Size  int    `ini:"size"`
	Count int    `ini:"count"`
}

type Dns struct {
	Upstream         string `ini:"upstream"`
	UpstreamFailover string `ini:"upstream_failover"`
}

type Filter struct {
	Lists string `ini:"lists"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	err := ini.MapTo(cfg, configPath)

	if cfg.DNS.Upstream == "" && cfg.DNS.UpstreamFailover == "" {
		var upstream string
		upstream, err = LoadDefaultUpstream()
		cfg.DNS.Upstream = upstream
	}
	return cfg, err
}
