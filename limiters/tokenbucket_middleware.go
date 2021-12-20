package limiters

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/rwxpeter/sliding-rate-limit/errors"
)

const (
	RateLimitLimit     = "X-RateLimit-Limit"
	RateLimitReset     = "X-RateLimit-Limit"
	RateLimitRemaining = "X-RateLimit-Remaining"
	RetryAfter         = "Retry-After"
)

type KeyFunc func(r *http.Request) (string, error)

func RealIP(headers ...string) KeyFunc {
	return (func(r *http.Request) (string, error) {
		for _, h := range headers {
			if v := r.Header.Get(h); v != "" {
				return v, nil
			}
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return "", nil
		}

		return ip, nil
	})
}

func HttpRateLimiter(l TokenBucketLimiter, handle KeyFunc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			key, err := handle(r)

			if err != nil {
				return
			}

			exceeded, val := l.Update(ctx, key).(errors.RateLimitExceeded)

			if val {
				limit := exceeded.Limit
				remaining := exceeded.Remaining
				reset := exceeded.Reset.UTC().Format(time.RFC1123)

				// Set HTTP headers for X-RateLimit-Limit, X-RateLimit-Limit, X-RateLimit-Remaining and Retry-After
				w.Header().Set(RateLimitLimit, strconv.FormatUint(uint64(limit), 10))
				w.Header().Set(RateLimitRemaining, strconv.FormatUint(uint64(remaining), 10))
				w.Header().Set(RateLimitReset, reset)

				w.Header().Set(RetryAfter, reset)
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
