package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type ServerConfig struct {
	Port string `env:"PORT,default=8000"`
	Host string `env:"HOST,default=0.0.0.0"`
}

type SMTPConfig struct {
	Login     string `env:"LOGIN,default=test@gmail.com"`
	Host      string `env:"HOST,default=smtp.gmail.com"`
	Port      string `env:"PORT,default=587"`
	Password  string `env:"PASSWORD,default=yourpassword"`
	Recipient string `env:"RECIPIENT,default=test@gmail.com"`
}

type Config struct {
	Server ServerConfig `env:", prefix=SERVER_"`
	SMTP   SMTPConfig   `env:", prefix=SMTP_"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
