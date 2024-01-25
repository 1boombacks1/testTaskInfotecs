package config

import (
	"fmt"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Postgres PostgresDatabase
		Server   Server
		Log      Log
	}

	PostgresDatabase struct {
		DSN string `env:"DSN"`
	}

	Server struct {
		Host string `env:"SRV_HOST"`
		Port string `env:"SRV_PORT" env-default:"3003"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}
)

func NewConfig(configPath ...string) (*Config, error) {
	cfg := &Config{}
	var err error

	if len(configPath) == 0 {
		err = cleanenv.ReadEnv(cfg)
	} else {
		err = cleanenv.ReadConfig(path.Join("./", configPath[0]), cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
