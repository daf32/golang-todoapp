package core_auth

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	JWTSecret                    string        `envconfig:"JWT_SECRET"                       required:"true"`
	AccessTokenExpiry            time.Duration `envconfig:"ACCESS_TOKEN_EXPIRY"              default:"15m"`
	RefreshTokenExpiry           time.Duration `envconfig:"REFRESH_TOKEN_EXPIRY"             default:"168h"`
	EmailConfirmationTokenExpiry time.Duration `envconfig:"EMAIL_CONFIRMATION_TOKEN_EXPIRY"  default:"24h"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("AUTH", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err := fmt.Errorf("get Auth config: %w", err)
		panic(err)
	}

	return config
}
