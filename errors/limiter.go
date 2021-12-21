package errors

import (
	"fmt"
	"time"
)

type RateLimitExceeded struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

func ErrorLimitExceeded(limit int, reset time.Time) error {
	return RateLimitExceeded{
		Limit: limit,
		Reset: reset,
	}
}

func (e RateLimitExceeded) Error() string {
	return fmt.Sprintf(
		"rate limit of %d has been exceeded and resets at %v",
		e.Limit,
		e.Reset,
	)
}

type GrcaRateLimtExceeded struct {
	Quantity   int64
	Remaining  int64
	RetryAfter int64
	ResetAfter int64
}

func GrcaRateLimitExceeded(quantity int64, remaining int64, retryAfter int64, resetAfter int64) error {
	return GrcaRateLimtExceeded{
		Quantity:   quantity,
		Remaining:  remaining,
		RetryAfter: retryAfter,
		ResetAfter: resetAfter,
	}
}

func (e GrcaRateLimtExceeded) Error() string {
	return fmt.Sprintf(
		"rate limit of %d has been exceeded with %d remaining and and resets at %v",
		e.Quantity,
		e.Remaining,
		e.ResetAfter,
	)
}
