package remote

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	ServicePort uint
	ServiceName string
	Path        string
	Body        interface{}
	Query       interface{}
	Headers     map[string]string
	ctx         *Context
	Timeout     *time.Duration
	err         error
	Secure      bool
}

// Req returns a new Request instance.
func Req() *Request {
	ctx := &Context{Context: context.Background()}
	req := &Request{
		Timeout: globalTimeout, // timeout can be override later by calling WithTimeout() in Request
		Headers: globalHeaders,
	}
	ctx.Request = req
	req.ctx = ctx
	ctx.handlers = append(ctx.handlers, globalHandlers...)
	return req
}

// Use adds middlewares to current Request.
func (r *Request) Use(handlerFuncs ...HandlerFunc) *Request {
	r.ctx.handlers = append(r.ctx.handlers, handlerFuncs...)
	return r
}

// Host sets remote service call name to current Request.
func (r *Request) Host(serviceName string) *Request {
	if serviceName == "" {
		r.err = fmt.Errorf("request service name not set")
		return r
	}
	r.ServiceName = serviceName
	return r
}

func (r *Request) Port(port uint) *Request {
	if port > 65535 {
		r.err = fmt.Errorf("invalid service port %d", port)
		return r
	}
	r.ServicePort = port
	return r
}

// HostAndPort sets remote service call name and port to current Request.
func (r *Request) HostAndPort(serviceName string, port uint) *Request {
	return r.Host(serviceName).Port(port)
}

// WithPath sets request path of current Request.
func (r *Request) WithPath(path string) *Request {
	r.Path = strings.TrimPrefix(path, "/")
	return r
}

// WithHeaders adds request headers of current Request.
// If the given headers contains same key with current existed headers, this header will be override.
func (r *Request) WithHeaders(headers map[string]string) *Request {
	if len(headers) > 0 {
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
	return r
}

func (r *Request) ContentType(contentType string) *Request {
	return r.contentHandle(contentType, ContentTypeHeader)
}

func (r *Request) Accept(contentType string) *Request {
	return r.contentHandle(contentType, AcceptTypeHeader)
}

func (r *Request) contentHandle(contentType string, headerKey string) *Request {
	switch contentType {
	case ContentTypeJson:
		fallthrough
	case ContentTypeFrom:
		if r.Headers == nil {
			r.Headers = make(map[string]string)
		}
		r.Headers[headerKey] = contentType
	default:
		r.err = fmt.Errorf("unsupported %s %s", headerKey, contentType)
	}
	return r
}

// WithQueries sets request query params of current Request.
func (r *Request) WithQueries(query interface{}) *Request {
	r.Query = query
	return r
}

// WithBody sets request body of current Request.
func (r *Request) WithBody(body interface{}) *Request {
	r.Body = body
	return r
}

// WithTimeout sets timeout for the whole processing chain of current Request. Note that this will override global timeout.
// Once timeout be reached, the HTTP request and all middlewares will be canceled from executing.
func (r *Request) WithTimeout(d time.Duration) *Request {
	r.Timeout = &d
	return r
}

func (r *Request) WithSecure(secure bool) *Request {
	r.Secure = secure
	return r
}

func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx.Context = ctx
	return r
}

// Get executes current Request using HTTP Get.
func (r *Request) Get(rsp Response) {
	r.ctx.Method = http.MethodGet
	r.do(rsp)
}

// Post executes current Request using HTTP Post.
func (r *Request) Post(rsp Response) {
	r.ctx.Method = http.MethodPost
	if _, p := r.ctx.Request.Headers[ContentTypeHeader]; !p {
		r.ContentType(ContentTypeJson)
	}
	r.do(rsp)
}

// Patch executes current Request using HTTP Post.
func (r *Request) Patch(rsp Response) {
	r.ctx.Method = http.MethodPatch
	r.do(rsp)
}

// Delete executes current Request using HTTP Delete.
func (r *Request) Delete(rsp Response) {
	r.ctx.Method = http.MethodDelete
	if _, p := r.ctx.Request.Headers[ContentTypeHeader]; !p {
		r.ContentType(ContentTypeJson)
	}
	r.do(rsp)
}

// Put executes current Request using HTTP Put.
func (r *Request) Put(rsp Response) {
	r.ctx.Method = http.MethodPut
	if _, p := r.ctx.Request.Headers[ContentTypeHeader]; !p {
		r.ContentType(ContentTypeJson)
	}
	r.do(rsp)
}

func (r *Request) do(rsp Response) {
	if r.err != nil {
		rsp.ErrorSave(r.err)
		return
	}
	r.ctx.Response = rsp
	r.ctx.handlers = append(r.ctx.handlers, doHttpReq)
	if r.Timeout != nil {
		var cancelFn context.CancelFunc
		r.ctx.Context, cancelFn = context.WithTimeout(context.TODO(), *r.Timeout)
		timer := time.NewTimer(*r.Timeout)
		done := make(chan struct{})
		go func() {
			r.ctx.handlers[0](r.ctx)
			close(done)
		}()
		select {
		case <-timer.C:
			cancelFn()
		case <-done:
			cancelFn()
			timer.Stop()
		}
	} else {
		r.ctx.handlers[0](r.ctx)
	}
}

func (r *Request) String() string {
	return fmt.Sprintf("[%s] %s://%s:%d/%s", r.ctx.Method, r.Scheme(), r.ServiceName, r.ServicePort, r.Path)
}

func (r *Request) Scheme() string {
	scheme := schemeHttp
	if r.Secure {
		scheme = schemeHttps
	}
	return scheme
}

func (r *Request) Error() error {
	return r.err
}
