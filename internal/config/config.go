package config

type Config struct {
	Logging LoggingConfig `ini:"logging"`
}

type LoggingConfig struct {
	Level string `ini:"level"`
}
