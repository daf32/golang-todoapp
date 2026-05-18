package core_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TimeZone           *time.Location `envconfig:"TIME_ZONE"                  default:"UTC"`
	JWTSecret          string         `envconfig:"AUTH_JWT_SECRET"            required:"true"`
	AccessTokenExpiry  time.Duration  `envconfig:"AUTH_ACCESS_TOKEN_EXPIRY"   default:"15m"`
	RefreshTokenExpiry time.Duration  `envconfig:"AUTH_REFRESH_TOKEN_EXPIRY"  default:"7d"`
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
