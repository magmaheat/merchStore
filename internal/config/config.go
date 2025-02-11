package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
	"time"
)

type HTTP struct {
	Host string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"HTTP_PORT" env-required:"true"`
}

type App struct {
	Name    string `env:"APP_NAME" env-required:"true"`
	Version string `env:"APP_VERSION" env-required:"true"`
}

type Log struct {
	Level string `env:"LOG_LEVEL" env-required:"true"`
}

type JWT struct {
	SigningKey     string        `env:"JWT_SIGNING_KEY" env-required:"true"`
	TokenAccessTTL time.Duration `env:"JWT_TOKEN_ACCESS_TTL" env-required:"true"`
}

type PG struct {
	URL string `env:"PG_URL" env-required:"true"`
}

type Config struct {
	App
	HTTP
	Log
	JWT
	PG
}

func New() *Config {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	return cfg
}
