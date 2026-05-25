package core_oauth

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type GoogleConfig struct {
	ClientID     string `envconfig:"CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"CLIENT_SECRET" required:"true"`
}

type Config struct {
	Google GoogleConfig
}

func NewConfig() (Config, error) {
	var google GoogleConfig
	if err := envconfig.Process("GOOGLE_OAUTH", &google); err != nil {
		return Config{}, fmt.Errorf("process Google OAuth config: %w", err)
	}

	return Config{Google: google}, nil
}

func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get OAuth config: %w", err))
	}

	return cfg
}
