package algorithms

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func FlushRedisDB(ctx context.Context, client *redis.Client) {
	val := client.FlushAll(ctx)
	if val.Err() != nil {
		panic("Unable to flush Redis with command FlushAll")
	}

	val = client.FlushDB(ctx)
	if val.Err() != nil {
		panic("Unable to flush Redis with command FlushDB")
	}
}

func NewRedis(ctx context.Context, limitPeriod time.Duration, rate int) (*Redis, func()) {
	r, close := NewClient(limitPeriod, rate)

	FlushRedisDB(ctx, r.client)

	return &r, close
}

func NewUser() string {
	return uuid.NewString()
}

func TestReturnErrorWhenBucketCapacityExhausted(t *testing.T) {
	ctx := context.Background()
	user := NewUser()
	r, close := NewRedis(ctx, time.Hour*1, 10)
	defer close()

	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.Error(t, r.Update(ctx, user))
}
func TestReturnErrorWhenMaxRateExceeded(t *testing.T) {
	ctx := context.Background()
	user := NewUser()
	r, close := NewRedis(ctx, time.Hour*1, 1)
	defer close()

	assert.NoError(t, r.Update(ctx, user))
	assert.Error(t, r.Update(ctx, user))
}

func TestRefilBucketCapacity(t *testing.T) {
	ctx := context.Background()
	user := NewUser()
	r, close := NewRedis(ctx, time.Millisecond*1, 1)
	defer close()

	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
	assert.NoError(t, r.Update(ctx, user))
}
