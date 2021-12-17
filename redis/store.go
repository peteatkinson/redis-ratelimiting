package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rwxpeter/sliding-rate-limit/limiter"
)

type Store struct {
	client        *redis.Client
	limit         int
	limitPeriod   time.Duration
	counterWindow time.Duration
}

func New(client *redis.Client, limit int, period, expiry time.Duration) *Store {
	return &Store{client, limit, period, expiry}
}

func (r *Store) Incr(ctx context.Context, key string, increment int) error {
	now := time.Now()
	timeStamp := fmt.Sprint(now.Truncate(r.counterWindow).Unix())

	// Increment the hash by 1
	val, err := r.client.HIncrBy(ctx, key, timeStamp, int64(increment)).Result()

	if err != nil {
		return err
	}

	if val == 1 {
		// Hash just created, so set the expiry time
		r.client.Expire(ctx, key, r.limitPeriod)
	} else if val >= int64(r.limit) {
		// Otherwise, check if just this fixed window counter period is over
		return limiter.ErrorLimitExceeded(0, r.limit, r.limitPeriod, now.Add(r.limitPeriod))
	}

	values, err := r.client.HGetAll(ctx, key).Result()

	if err != nil {
		return err
	}

	// The time to start summing from, any buckets before this are ignored.
	threshold := fmt.Sprint(now.Add(-r.limitPeriod).Unix())

	// NOTE: this sums ALL the values in the hash, for more information see the
	// "Practical Considerations" section of the associated Figma blog post.
	total := 0
	for k, v := range values {
		if k > threshold {
			i, _ := strconv.Atoi(v)
			total += i
		} else {
			// Clear out the old hash keys
			r.client.HDel(ctx, key, k)
		}
	}

	if total >= int(r.limit) {
		return limiter.ErrorLimitExceeded(0, r.limit, r.limitPeriod, now.Add(r.limitPeriod))
	}

	return nil
}
