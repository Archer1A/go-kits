package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response interface {
	Error() error
	ErrorSave(e error)
	SetRaw(raw *http.Response)
	HttpResponse() *http.Response
}

type DefaultResponse struct {
	err  error
	Data interface{} `json:"data"`
	Raw  *http.Response
}

func (d *DefaultResponse) Error() error {
	return d.err
}

func (d *DefaultResponse) ErrorSave(e error) {
	d.err = e
}

func (d *DefaultResponse) SetRaw(raw *http.Response) {
	d.Raw = raw
}

func (d *DefaultResponse) HttpResponse() *http.Response {
	return d.Raw
}

func (d *DefaultResponse) String() string {
	if d.err != nil {
		return fmt.Sprintf("err: %v", d.err)
	}
	return fmt.Sprintf("status: %d, errMsg: %s, data: %v", d.HttpResponse().StatusCode, d.err, d.Data)
}

func (d *DefaultResponse) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, d.Data)
	return err
}
