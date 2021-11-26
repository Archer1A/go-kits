package remote

import (
	"context"
	"fmt"
	"time"
)

type HandlerFunc func(*Context)

var globalHandlers []HandlerFunc
var globalTimeout *time.Duration
var globalHeaders map[string]string

// Use adds global middlewares which will take effect in every Request
func Use(middlewares ...HandlerFunc) {
	globalHandlers = append(globalHandlers, middlewares...)
}

// Timeout set global timeout for each Request.
// If a middleware intercepts the Request, and the HTTP request never be fired, this timeout will take no effect.
func Timeout(d time.Duration) {
	globalTimeout = &d
}

// Headers add global headers for each Request.
func Headers(headers map[string]string) {
	if globalHeaders == nil {
		globalHeaders = make(map[string]string)
	}
	for k, v := range headers {
		globalHeaders[k] = v
	}
}

type Context struct {
	Request  *Request
	Response Response
	Method   string
	handlers []HandlerFunc
	index    int
	params   map[string]interface{}
	context.Context
}

// Next should only be called inside middleware.
// It executes the pending handlers in the handler chain after the calling handler.
func (c *Context) Next() {
	for c.index+1 < len(c.handlers) {
		c.index++
		c.handlers[c.index](c)
	}
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
func (c *Context) Abort() {
	// When the last handler has been called, c.index = len(c.handlers).
	// So we need a larger index to indicate that Context is aborted
	c.index = len(c.handlers) + 1
}

// IsAborted returns whether current Context has been aborted
func (c *Context) IsAborted() bool {
	return c.index > len(c.handlers)
}

// Set store a new key value pair exclusively for current context.
// The under laying key-value map is lazy-initialized.
func (c *Context) Set(key string, value interface{}) {
	if c.params == nil {
		c.params = make(map[string]interface{})
	}
	c.params[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the key dose not exists it returns (nil, false).
func (c *Context) Get(key string) (val interface{}, present bool) {
	if c.params == nil {
		return nil, false
	}
	val, present = c.params[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	val, present := c.Get(key)
	if !present {
		panic(fmt.Sprintf("param %s not found in Context", key))
	}
	return val
}
