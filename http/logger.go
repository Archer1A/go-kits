package remote

import (
	"fmt"
	"io"
	"os"
	"time"
)

var DefaultWriter io.Writer = os.Stdout

type LogFormatterParams struct {
	Request    *Request
	Response   Response
	Method     string
	Path       string
	StatusCode int
	ErrMessage string
	Timestamp  time.Time
	Latency    time.Duration
}

type LoggerFormatter func(param LogFormatterParams) string

var defaultLogFormatter = func(param LogFormatterParams) string {
	if param.ErrMessage != "" {
		return fmt.Sprintf("[%s] | %10v | %s | %s://%s:%d/%s\n", param.Method, param.Latency, param.ErrMessage, param.Request.Scheme(), param.Request.ServiceName, param.Request.ServicePort, param.Path)
	}
	return fmt.Sprintf("[%s] | %10v | %d | %s://%s:%d/%s\n", param.Method, param.Latency, param.StatusCode, param.Request.Scheme(), param.Request.ServiceName, param.Request.ServicePort, param.Path)
}

type LoggerConfig struct {
	Formatter LoggerFormatter
	Output    io.Writer
}

// Logger returns a middleware that log into remote.DefaultWriter.
// By default remote.DefaultWriter = os.Stdout
func Logger() HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithFormatter returns a middleware with the specified log format.
func LoggerWithFormatter(formatter LoggerFormatter) HandlerFunc {
	return LoggerWithConfig(LoggerConfig{Formatter: formatter})
}

func LoggerWithConfig(config LoggerConfig) HandlerFunc {
	formatter := config.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := config.Output
	if out == nil {
		out = DefaultWriter
	}

	return func(ctx *Context) {
		start := time.Now()

		ctx.Next()

		param := LogFormatterParams{
			Request:  ctx.Request,
			Response: ctx.Response,
			Method:   ctx.Method,
			Path:     ctx.Request.Path,
		}

		param.Timestamp = time.Now()
		param.Latency = param.Timestamp.Sub(start)

		if ctx.Response.HttpResponse() != nil {
			param.StatusCode = ctx.Response.HttpResponse().StatusCode
		}
		if ctx.Request.err != nil {
			param.ErrMessage = ctx.Request.err.Error()
		} else if ctx.Response.Error() != nil {
			param.ErrMessage = ctx.Response.Error().Error()
		}

		_, _ = fmt.Fprint(out, formatter(param))
	}
}
