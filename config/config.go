package config

import (
	"fmt"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/providers/env/v2"
)

type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	JWT      JWTConfig      `koanf:"jwt"`
	Google   GoogleConfig   `koanf:"google"`
}

type ServerConfig struct {
	Port string `koanf:"port"`
}

type DatabaseConfig struct {
	URL string `koanf:"url"`
}

type JWTConfig struct {
	Secret string `koanf:"secret"`
}

type GoogleConfig struct {
	ClientID            string `koanf:"client_id"`
	ClientSecret        string `koanf:"client_secret"`
	RedirectURI         string `koanf:"redirect_uri"`
	FrontendCallbackURL string `koanf:"frontend_callback_url"`
}

func Load() (*Config, error) {
	k := koanf.New(".")

	err := k.Load(
		env.Provider("", env.Opt{
			Prefix: "",
		}),
		nil,
	)

	if err != nil {
		return nil, err
	}

	// تبدیل env های flat به ساختار config
	err = k.Set("server.port", k.String("SERVER_PORT"))
	if err != nil {
		return nil, err
	}

	err = k.Set("database.url", k.String("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	err = k.Set("jwt.secret", k.String("JWT_SECRET"))
	if err != nil {
		return nil, err
	}

	err = k.Set("google.client_id", k.String("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}

	err = k.Set("google.client_secret", k.String("GOOGLE_CLIENT_SECRET"))
	if err != nil {
		return nil, err
	}

	err = k.Set("google.redirect_uri", k.String("GOOGLE_REDIRECT_URI"))
	if err != nil {
		return nil, err
	}

	err = k.Set("google.frontend_callback_url", k.String("FRONTEND_CALLBACK_URL"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	err = k.Unmarshal("", cfg)

	if err != nil {
		return nil, err
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