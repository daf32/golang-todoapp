package core_ratelimit

import "time"

type Limiter interface {
	Allow(key string) (allowed bool, retryAfter time.Duration)
}
