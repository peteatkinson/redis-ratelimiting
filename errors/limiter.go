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
