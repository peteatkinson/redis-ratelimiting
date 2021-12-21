package limiters

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rwxpeter/sliding-rate-limit/errors"
)

type Redis struct {
	client      *redis.Client
	limitPeriod time.Duration
	rate        int
}

type TokenBucketLimiter interface {
	Update(context.Context, string) error
}

func New(client *redis.Client, limitPeriod time.Duration, rate int) *Redis {
	return &Redis{client, limitPeriod, rate}
}

func (r *Redis) Update(ctx context.Context, key string) error {
	now := time.Now()
	timeStamp := now.Unix()

	values, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		// Redis Connection issue
		return err
	}

	if len(values) == 0 {
		// No hash exists at that key, we can create a new one (-1 on rate as first hit is one use of a token)
		take(ctx, r.client, key, r.rate, timeStamp)
		return nil
	}

	// Get timestamp from Redis HASH
	ts, err := strconv.ParseInt(values["ts"], 10, 32)
	if err != nil {
		return err
	}

	deadline := time.Unix(ts, 0).Add(r.limitPeriod).Unix()
	// IF the current time distance between now and the last access timestamp is higher than the deadline timestamp
	if timeStamp >= deadline {
		// refill the bucket
		update(ctx, r.client, key, r.rate, timeStamp)
	} else {
		// Otherwise check if we have any tokens left in the bucket
		if tokens, _ := strconv.Atoi(values["tokens"]); tokens > 0 {
			// If we do then we can update the Hash with the new token count and timestamp
			take(ctx, r.client, key, r.rate, timeStamp)
		} else {
			// Otherwise, we've exceeded the limit before the deadline, return error
			return errors.ErrorLimitExceeded(r.rate, time.Unix(deadline, 0))
		}
	}

	return nil
}

func update(ctx context.Context, client *redis.Client, key string, rate int, ts int64) {
	client.HSet(ctx, key, "tokens", fmt.Sprint(rate), "ts", fmt.Sprint(ts))
}

func take(ctx context.Context, client *redis.Client, key string, rate int, ts int64) {
	update(ctx, client, key, rate-1, ts)
}

func NewClient(limitPeriod time.Duration, rate int) (Redis, func()) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := New(client, limitPeriod, rate)
	return *r, func() { client.Close() }
}
