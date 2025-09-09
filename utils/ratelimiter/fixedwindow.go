package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func IsAllowFixedWindow(ctx context.Context, client *redis.Client, userId string, limiter int64, window time.Duration) (bool, error) {
	nano := time.Now().UnixNano()
	//fixedwindow:test-user1:1746063952
	key := fmt.Sprintf("fixedwindow:%s:%d", userId, nano/int64(window))
	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return false, nil
	}
	if count == 1 {
		client.Expire(ctx, key, window)
	}
	//fixedwindow:test-user1:1746063952
	return count <= limiter, nil
}
