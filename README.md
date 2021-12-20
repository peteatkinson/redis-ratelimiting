# Redis Ratelimiting in Go


This repository provides a few exporatory examples of rate limiting algorithms written in Go. 
It includes (or atleast will include) a few examples of using Redis to create Ratelimiters by using a few of the common algorithms for rate limiting.
These include, the token-bucket, leaky-bucket and sliding-log alogrithms.

The approaches are taken from [this blog post](https://www.figma.com/blog/an-alternative-approach-to-rate-limiting) from the Figma Engineering team

## Token Bucket
Needs writting....

## Leaky Bucket
Needs writting....

## Sliding Log (or window Sliding Log)
Needs writting....

## Usage
### Token Bucket

The token bucket satisfies this interface:

```go
type TokenBucketLimiter interface {
	Update(context.Context, string) error
}
```

The implementation for the token bucket is backed by Redis and uses the go-redis library.

```go
client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

rateLimiter := limiters.New(client, time.Minute*10, 5)

ratelimiter.Update(ctx, "user_id")
```

## Middleware
