package limiters

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Limiter interface {
	Run(ctx context.Context, key string)
}

type GRCARateLimiter struct {
	client                  *redis.Client
	delayVariationTolerance time.Duration
	emissionInterval        time.Duration
	limit                   int64
}

func NewGRCARateLimiter(client *redis.Client, quota Quota) (*GRCARateLimiter, error) {
	if quota.maxburst < 0 {
		return nil, fmt.Errorf("invalid Quota  %#v; MaxBurst must be greater than zero", quota)
	}
	if quota.maxRate.period <= 0 {
		return nil, fmt.Errorf("invalid RateQuota %#v; MaxRate must be greater than zero", quota)
	}
	return &GRCARateLimiter{
		delayVariationTolerance: quota.maxRate.period * (time.Duration(quota.maxburst) + 1),
		emissionInterval:        quota.maxRate.period,
		limit:                   int64(quota.maxburst) + 1,
		client:                  client,
	}, nil
}

func NewQuota(maxRate Rate, maxBurst int) Quota {
	return Quota{
		maxRate:  maxRate,
		maxburst: maxBurst,
	}
}

type Quota struct {
	maxRate  Rate
	maxburst int
}

type Rate struct {
	period time.Duration
	count  int
}

func PerSecond(n int) Rate { return Rate{time.Second / time.Duration(n), n} }

func PerMinute(n int) Rate { return Rate{time.Minute / time.Duration(n), n} }

func PerHour(n int) Rate { return Rate{time.Hour / time.Duration(n), n} }

func PerDay(n int) Rate { return Rate{24 * time.Hour / time.Duration(n), n} }

func PerDuration(n int, d time.Duration) Rate { return Rate{d / time.Duration(n), n} }

func calculateExpectedTimeOfArrival(tm string, timeStamp int64) int64 {
	var tat int64
	if tm == "" {
		tat = timeStamp
	} else {
		tat, _ = strconv.ParseInt(tm, 10, 32)
	}
	return tat
}

func calculateNewTimeOfArrival(tat int64, increment time.Duration, now time.Time) time.Time {
	var newTat time.Time
	if now.After(time.Unix(tat, 0)) {
		newTat = now.Add(increment)
	} else {
		newTat = time.Unix(tat, 0).Add(increment)
	}
	return newTat
}

func (c *GRCARateLimiter) Execute(ctx context.Context, key string, quantity int64) error {
	now, _ := c.client.Time(ctx).Result()

	timeStamp := now.Unix()

	tm, _ := c.client.Get(ctx, key).Result()

	tat := calculateExpectedTimeOfArrival(tm, timeStamp)

	fmt.Println(tat)

	limit := c.limit
	increment := time.Duration(quantity) * c.emissionInterval

	newTat := calculateNewTimeOfArrival(tat, increment, now)

	allowAt := newTat.Add(-(c.delayVariationTolerance))
	allowAtStr := allowAt.String()

	var ttl time.Duration
	limited := false
	if diff := now.Sub(allowAt); diff < 0 {
		if increment <= c.delayVariationTolerance {
			ttl = time.Unix(tat, 0).Sub(now)
		}
		limited = true
	}

	ttl = newTat.Sub(now)

	c.client.Set(ctx, key, fmt.Sprint(newTat.Unix()), ttl)

	fmt.Println(newTat, limit, allowAt, allowAtStr, ttl, limited)
	return nil
}
