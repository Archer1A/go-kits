package logger

import (
	"context"
	"errors"
	"testing"
	"time"
)

func Test_Info(t *testing.T) {
	type req struct {
		Url  string
		Time time.Time
	}
	r := req{
		Url:  "www.baidu.com",
		Time: time.Now(),
	}
	//zapLogger.Info("hello", zap.String("name", "vic"))
	//sugarLogger.Info("hello", "a", "b")
	//sugarLogger.Infow("hello", &r)
	Infow(context.TODO(), "hello ", "name", "vic")
	Infow(context.TODO(), "request ", "req", r)
	Infow(context.TODO(), "error ", errors.New("404"))

}
