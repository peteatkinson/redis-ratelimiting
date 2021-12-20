# Redis Ratelimiting in Go


This repository provides a few exporatory examples of rate limiting algorithms written in Go. 
It includes (or atleast will include) a few examples of using Redis to create Ratelimiters by using a few of the common algorithms for rate limiting.
These include, the token-bucket, leaky-bucket and sliding-log alogrithms.

The approaches are taken from [this blog post](https://www.figma.com/blog/an-alternative-approach-to-rate-limiting) from the Figma Engineering team

This is honestly just an experimental bit of coding for my own personal development, so take all this code with caution.
## Token Bucket

The first and proberbly the most simpliest of algorithms to implmenet rate-limiting with is known as the Token bucket. 

Simply, how it works - is we keep track of the requests coming in with a Redis Hash.

See below for an example for each unique request that comes in.
```

127.0.0.1 (user-1): {"ts": "1639986575", "tokens": 5}
```

The token bucket keeps track of the timestamp and the total remaining tokens left. If all tokens are exhuasted within a given time window then we drop an incoming request for that particular user. After the time window is up, we refill the tokens within the hash and the cycle continues.


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
