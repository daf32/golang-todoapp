package users_service

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	UnverifiedCleanupInterval time.Duration `envconfig:"UNVERIFIED_CLEANUP_INTERVAL" default:"24h"`
	UnverifiedCleanupMinAge   time.Duration `envconfig:"UNVERIFIED_CLEANUP_MIN_AGE"  default:"168h"`
}

func NewConfig() (Config, error) {
	var c Config
	if err := envconfig.Process("USERS", &c); err != nil {
		return Config{}, fmt.Errorf("process users config: %w", err)
	}

	return c, nil
}

func NewConfigMust() Config {
	c, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("get users config: %w", err))
	}

	return c
}
