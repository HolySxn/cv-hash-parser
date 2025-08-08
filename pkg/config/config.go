package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type ServerConfig struct {
	Port   string `env:"PORT,default=8000"`
	Host   string `env:"HOST,default=0.0.0.0"`
	LogLvl string `env:"LOG_LEVEL,default=debug"`
}

type SMTPConfig struct {
	From string `env:"FROM,default=test@gmail.com"`
	Host string `env:"SMTP_HOST,default=smtp.gmail.com"`
	Port string `env:"SMTP_PORT,default=587"`
	Pass string `env:"SMTP_PASS,default=yourpassword"`
}

type Config struct {
	Server ServerConfig `envPrefix:"SERVER_"`
	SMTP   SMTPConfig   `envPrefix:"SMTP_"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
