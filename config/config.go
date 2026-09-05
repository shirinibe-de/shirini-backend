package config

import (
	"fmt"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/providers/env/v2"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Google   GoogleConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func Load() (*Config, error) {
	k := koanf.New(".")

	err := k.Load(env.Provider("", env.Opt{}), nil)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: k.String("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			URL: k.String("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret: k.String("JWT_SECRET"),
		},
		Google: GoogleConfig{
			ClientID:     k.String("GOOGLE_CLIENT_ID"),
			ClientSecret: k.String("GOOGLE_CLIENT_SECRET"),
			RedirectURI:  k.String("GOOGLE_REDIRECT_URI"),
		},
	}

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	return cfg, nil
}