package remote

import (
	"context"
	"go.uber.org/atomic"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func TestRateLimitAllow(t *testing.T) {
	type args struct {
		limit  rate.Limit
		burst  int
		getter BucketGetter
	}
	tests := []struct {
		name string
		args args
		run  func(h HandlerFunc) bool
	}{
		{
			name: "10req/1s pass",
			args: args{
				limit: rate.Every(time.Second),
				burst: 10,
				getter: func(ctx *Context) interface{} {
					return 1
				},
			},
			run: func(h HandlerFunc) bool {
				for i := 0; i < 10; i++ {
					ctx := &Context{Response: &DefaultResponse{}}
					h(ctx)
					if ctx.Response.Error() != nil {
						return false
					}
				}
				return true
			},
		},
		{
			name: "10req/1s fail",
			args: args{
				limit: rate.Every(time.Second),
				burst: 10,
				getter: func(ctx *Context) interface{} {
					return 1
				},
			},
			run: func(h HandlerFunc) bool {
				for i := 0; i < 20; i++ {
					ctx := &Context{Response: &DefaultResponse{}}
					h(ctx)
					if ctx.Response.Error() != nil {
						return true
					}
				}
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RateLimitAllow(tt.args.limit, tt.args.burst, tt.args.getter)
			if !tt.run(got) {
				t.Error("RateLimitAllow() test failed")
			}
		})
	}
}

func TestRateLimitWait(t *testing.T) {
	type args struct {
		limit  rate.Limit
		burst  int
		getter BucketGetter
	}
	tests := []struct {
		name string
		args args
		run  func(h HandlerFunc) bool
	}{
		{
			name: "10req/1s",
			args: args{
				limit: rate.Every(time.Second),
				burst: 10,
				getter: func(ctx *Context) interface{} {
					return 1
				},
			},
			run: func(h HandlerFunc) bool {
				ch := make(chan bool, 20)
				for i := 0; i < 20; i++ {
					go func() {
						timeoutCtx, fn := context.WithTimeout(context.Background(), time.Millisecond*999)
						defer fn()
						ctx := &Context{Context: timeoutCtx, Response: &DefaultResponse{}}
						h(ctx)
						if ctx.Response.Error() != nil {
							ch <- false
							return
						}
						ch <- true
					}()
				}
				atomInt := atomic.Int32{}
				for i := 0; i < 20; i++ {
					b := <-ch
					if b {
						atomInt.Inc()
					}
					if i == 19 {
						close(ch)
						break
					}
				}
				return atomInt.Load() == 10
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RateLimitWait(tt.args.limit, tt.args.burst, tt.args.getter)
			if !tt.run(got) {
				t.Error("RateLimitWait() test failed")
			}
		})
	}
}
