package algorithms

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	r := New(client, time.Hour*1, 5)

	assert.NoError(t, r.Update(context.Background(), "user7"))
	assert.NoError(t, r.Update(context.Background(), "user7"))
	assert.NoError(t, r.Update(context.Background(), "user7"))
	assert.NoError(t, r.Update(context.Background(), "user7"))
	assert.NoError(t, r.Update(context.Background(), "user7"))
}
