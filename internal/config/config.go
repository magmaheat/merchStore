package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
	"path"
	"time"
)

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type App struct {
	Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
	Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
}

type Log struct {
	Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL" `
}

type JWT struct {
	SigningKey     string        `env-required:"true" env:"JWT_SIGNING_KEY" `
	TokenAccessTTL time.Duration `env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
}

type PG struct {
	URL string `env:"PG_URL" env-required:"true"`
}

type Config struct {
	App  `yaml:"app"`
	HTTP `yaml:"http"`
	Log  `yaml:"log"`
	JWT  `yaml:"jwt"`
	PG   `yaml:"postgres"`
}

func New(configPath string) *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return cfg
}
