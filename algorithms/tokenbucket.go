package algorithms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client      *redis.Client
	limitPeriod time.Duration
	rate        int
}

type RateLimiter interface {
	Update(context.Context, int) error
}

type RateLimitExceeded struct {
	error
	Limit int
	Reset time.Time
}

func ErrorLimitExceeded(limit int, reset time.Time) error {
	return RateLimitExceeded{
		Limit: limit,
		Reset: reset,
	}
}

func (e RateLimitExceeded) ErrorString() string {
	return fmt.Sprintf(
		"rate limit of %d has been exceeded and resets at %v",
		e.Limit,
		e.Reset,
	)
}

func New(client *redis.Client, limitPeriod time.Duration, rate int) *Redis {
	return &Redis{client, limitPeriod, rate}
}

func (r *Redis) Update(ctx context.Context, key string) error {
	now := time.Now()
	unixNow := now.Unix()

	values, err := r.client.HGetAll(ctx, key).Result()

	if err != nil {
		// Redis Connection issue
		return nil
	}

	if len(values) == 0 {
		// No hash exists at that key, we can create a new one (-1 on rate as first hit is one use of a token)
		r.client.HSet(ctx, key, "tokens", fmt.Sprint(r.rate-1), "ts", fmt.Sprint(unixNow))
		return nil
	}

	// Get timestamp from Redis HASH
	ts, err := strconv.ParseInt(values["ts"], 10, 32)

	if err != nil {
		return err
	}

	// Get tokens from Redis HASH
	tokens, err := strconv.Atoi(values["tokens"])

	if err != nil {
		return err
	}

	// Do a bunch of convertions on the STRING values
	tsUnix := time.Unix(ts, 0)
	delta := unixNow

	deadline := tsUnix.Add(r.limitPeriod).Unix()

	// IF the current time distance between now and the last access timestamp is higher than the deadline timestamp
	if delta >= deadline {
		// refill the bucket
		r.client.HSet(ctx, key, "tokens", fmt.Sprint(r.rate), "ts", fmt.Sprint(unixNow))
	} else {
		// Otherwise check if we have any tokens left in the bucket
		if tokens > 0 {
			// If we do then we can update the Hash with the new token count and timestamp
			remainder := tokens - 1
			r.client.HSet(ctx, key, "tokens", fmt.Sprint(remainder), "ts", ts)
		} else {
			// Otherwise, we've exceeded the limit before the deadline, return error
			return ErrorLimitExceeded(r.rate, time.Unix(deadline, 0))
		}
	}

	return nil
}

func NewClient(limitPeriod time.Duration, rate int) (Redis, func()) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := New(client, limitPeriod, rate)
	return *r, func() { client.Close() }
}
