package http

import (
	"errors"
	"golang.org/x/time/rate"
	"sync"
)

type buckets struct {
	buks map[interface{}]*rate.Limiter
	mu   sync.Mutex
}

func (b *buckets) getOrCreate(key interface{}, limit rate.Limit, burst int) *rate.Limiter {
	limiter, ok := b.buks[key]
	if !ok {
		b.mu.Lock()
		if _, ok = b.buks[key]; !ok {
			limiter = rate.NewLimiter(limit, burst)
			b.buks[key] = limiter
		}
		b.mu.Unlock()
	}
	return limiter
}

// BucketGetter defines how to get specified bucket when check rate limit
type BucketGetter func(ctx *Context) interface{}

var RateLimitExceedError = errors.New("rate limit exceed")

// RateLimitAllow returns a rate limiter which will abort request when here is no token could be obtained in bucket in now time.
func RateLimitAllow(limit rate.Limit, burst int, getter BucketGetter) HandlerFunc {
	return RateLimitWithHandle(limit, burst, getter, func(ctx *Context, limiter *rate.Limiter) {
		if limiter.Allow() {
			ctx.Next()
		} else {
			ctx.Abort()
			ctx.Response.ErrorSave(RateLimitExceedError)
		}
	})
}

// RateLimitWait returns a rate limiter which will wait util here is one token could be obtained in bucket, or abort after deadline exceed.
//
// Please make sure that context deadline is set
func RateLimitWait(limit rate.Limit, burst int, getter BucketGetter) HandlerFunc {
	return RateLimitWithHandle(limit, burst, getter, func(ctx *Context, limiter *rate.Limiter) {
		if err := limiter.Wait(ctx.Context); err != nil {
			ctx.Abort()
			ctx.Response.ErrorSave(err)
		}
	})
}

// RateLimitWithHandle return a rate limiter and pass it into the handle function
func RateLimitWithHandle(limit rate.Limit, burst int, getter BucketGetter, handle func(ctx *Context, limiter *rate.Limiter)) HandlerFunc {
	buks := buckets{
		buks: make(map[interface{}]*rate.Limiter),
		mu:   sync.Mutex{},
	}

	return func(ctx *Context) {
		limiter := buks.getOrCreate(getter(ctx), limit, burst)
		handle(ctx, limiter)
	}
}
