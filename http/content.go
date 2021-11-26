package http

import (
	"encoding/json"
	"errors"
	"fmt"
)

var contentTypeRegistry map[string]ContentTypeResolver

func init() {
	contentTypeRegistry = make(map[string]ContentTypeResolver)
	contentTypeRegistry[ContentTypeJson] = &contentTypeJson{}
	contentTypeRegistry[ContentTypeFrom] = &contentTypeForm{}
}

type ContentTypeResolver interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// RegisterContentResolver register specified ContentTypeResolver for the given contentType
func RegisterContentResolver(contentType string, resolver ContentTypeResolver) {
	if _, ok := contentTypeRegistry[contentType]; ok {
		_, _ = fmt.Fprintf(DefaultWriter, "trying to override ContentTypeResolver for ContentType %s to %v", contentType, resolver)
	}
	contentTypeRegistry[contentType] = resolver
}

const (
	ContentTypeHeader = "Content-Type"
	AcceptTypeHeader  = "Accept"
	ContentTypeJson   = "application/json"
	ContentTypeFrom   = "application/x-www-form-urlencoded"
)

type contentTypeJson struct{}

func (c *contentTypeJson) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *contentTypeJson) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type contentTypeForm struct{}

func (c *contentTypeForm) Marshal(v interface{}) ([]byte, error) {
	m, err := formToMap(v, formTagName)
	if err != nil {
		return nil, err
	}
	queries := encodeForm(m)
	return []byte(queries), nil
}

func (c *contentTypeForm) Unmarshal(bytes []byte, v interface{}) error {
	return errors.New("implement me")
}
