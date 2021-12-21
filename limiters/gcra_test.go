package limiters

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rwxpeter/sliding-rate-limit/errors"
)

func Test1(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_ = rdb.FlushDB(ctx).Err()

	quota := NewQuota(PerHour(5), 0)

	limiter, _ := NewGRCARateLimiter(rdb, quota)
	user := uuid.NewString()
	val, _ := limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)
	val, _ = limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)
	val, _ = limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)
	val, _ = limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)
	val, _ = limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)
	val, _ = limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)
	val, _ = limiter.Execute(ctx, user, 1).(errors.GrcaRateLimtExceeded)


	fmt.Println(val)
}
