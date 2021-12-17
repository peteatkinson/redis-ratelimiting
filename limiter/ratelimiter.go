package limiter

import (
	"context"
	"fmt"
	"time"
)

type RateLimiter interface {
	Incr(context.Context, string, int) error
}

type RateLimitExceeded struct {
	error
	Remaining int
	Limit     int
	Period    time.Duration
	Reset     time.Time
}

func ErrorLimitExceeded(remaining int, limit int, period time.Duration, reset time.Time) error {
	return RateLimitExceeded{
		Remaining: remaining,
		Limit:     limit,
		Period:    period,
		Reset:     reset,
	}
}

func (e RateLimitExceeded) ErrorString() string {
	return fmt.Sprintf(
		"rate limit of %d per %v has been exceeded and resets at %v",
		e.Limit,
		e.Period,
		e.Reset,
	)
}
