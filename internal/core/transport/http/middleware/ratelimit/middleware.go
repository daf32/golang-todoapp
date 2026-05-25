package core_ratelimit

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/daf32/golang-todoapp/internal/core/errors"
	core_logger "github.com/daf32/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/daf32/golang-todoapp/internal/core/transport/http/middleware"
	core_http_response "github.com/daf32/golang-todoapp/internal/core/transport/http/response"
)

type KeyFunc func(r *http.Request) string

func ByIP() KeyFunc {
	return func(r *http.Request) string {
		return ClientIP(r)
	}
}

func Middleware(limiter Limiter, keyFn KeyFunc) core_http_middleware.Middleware {
	if keyFn == nil {
		keyFn = ByIP()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			allowed, retryAfter := limiter.Allow(keyFn(r))
			if !allowed {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

				if retryAfter > 0 {
					secs := int(retryAfter.Seconds())
					if secs < 1 {
						secs = 1
					}
					rw.Header().Set("Retry-After", strconv.Itoa(secs))
				}

				responseHandler.ErrorResponse(
					fmt.Errorf("rate limit exceeded: %w", core_errors.ErrTooManyRequests),
					"rate limit exceeded",
				)

				return 
			}

			next.ServeHTTP(rw, r)
		})
	}
}
