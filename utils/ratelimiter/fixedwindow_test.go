package ratelimiter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func Test_isAllowFixedWindow(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		DB:       0,
	})
	type args struct {
		ctx     context.Context
		client  *redis.Client
		userId  string
		limiter int64
		window  time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:     context.Background(),
				client:  rdb,
				userId:  "test-user1",
				limiter: 5,
				window:  time.Second,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsAllowFixedWindow(tt.args.ctx, tt.args.client, tt.args.userId, tt.args.limiter, tt.args.window)
			if (err != nil) != tt.wantErr {
				t.Errorf("isAllowFixedWindow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isAllowFixedWindow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAllowFixedWindow(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		DB:       0,
	})
	for i := 0; i < 10; i++ {
		allowed, err := IsAllowFixedWindow(context.Background(), rdb, "user1", 5, time.Second)
		if err != nil {
			return
		}
		if !allowed {
			fmt.Printf("user1 not allowed pass index %d \n", i)
		}

	}
}
