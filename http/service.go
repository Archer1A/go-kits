package remote

import "time"

// Service represents a remote service
type Service struct {
	Host        string
	Port        uint
	Middlewares []HandlerFunc
	Headers     map[string]string
	Timeout     string
	Secure      bool
}

// Serve create Request from Service
func (s *Service) Serve() *Request {
	request := Req().Host(s.Host).Port(s.Port).Use(s.Middlewares...).WithHeaders(s.Headers).WithSecure(s.Secure)
	if s.Timeout != "" {
		duration, err := time.ParseDuration(s.Timeout)
		request.err = err
		request.WithTimeout(duration)
	}
	return request
}
