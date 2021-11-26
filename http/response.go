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
	err    error
	Status int         `json:"status"`
	ResMsg string      `json:"resMsg"`
	Data   interface{} `json:"data"`
	raw    *http.Response
}

type ListResponse struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

func (d *DefaultResponse) Error() error {
	return d.err
}

func (d *DefaultResponse) ErrorSave(e error) {
	d.err = e
}

func (d *DefaultResponse) SetRaw(raw *http.Response) {
	d.raw = raw
}

func (d *DefaultResponse) HttpResponse() *http.Response {
	return d.raw
}

func (d *DefaultResponse) String() string {
	if d.err != nil {
		return fmt.Sprintf("err: %v", d.err)
	}
	return fmt.Sprintf("status: %d, errMsg: %s, data: %v", d.Status, d.ResMsg, d.Data)
}

type HarborResponse struct {
	err  error
	Data interface{} `json:"data"`
	raw  *http.Response
}

func (d *HarborResponse) Error() error {
	return d.err
}

func (d *HarborResponse) ErrorSave(e error) {
	d.err = e
}

func (d *HarborResponse) SetRaw(raw *http.Response) {
	d.raw = raw
}

func (d *HarborResponse) HttpResponse() *http.Response {
	return d.raw
}

func (d *HarborResponse) String() string {
	if d.err != nil {
		return fmt.Sprintf("err: %v", d.err)
	}
	return fmt.Sprintf("status: %d, errMsg: %s, data: %v", d.HttpResponse().StatusCode, d.err, d.Data)
}

func (d *HarborResponse) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, d.Data)
	return err
}
