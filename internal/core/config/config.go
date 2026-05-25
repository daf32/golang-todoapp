package core_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TimeZone     *time.Location `envconfig:"TIME_ZONE"     default:"UTC"`
	AppBaseURL   string         `envconfig:"APP_BASE_URL"  required:"true"`
	CookieSecure bool           `envconfig:"COOKIE_SECURE" default:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("get core config: %w", err)
		panic(err)
	}

	return config
}
