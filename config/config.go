package config

import (
	"fmt"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	DatabaseURL string `koanf:"DATABASE_URL"`
}

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	err := k.Load(env.Provider("SHIRINI", "_", nil), nil)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	var cfg Config
	err = k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{
		Tag:       "koanf",
		FlatPaths: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &cfg, nil
}
