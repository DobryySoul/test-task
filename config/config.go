package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP `yaml:"http"`
	Log  `yaml:"logger"`
	PG   `yaml:"postgres"`
}

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"PORT"`
}

type Log struct {
	Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
}

type PG struct {
	URL string `env-required:"true" yaml:"pg_url" env:"PG_URL"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read yml config: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read env: %w", err)
	}

	return cfg, nil
}
