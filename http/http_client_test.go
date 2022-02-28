package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_queryMap(t *testing.T) {
	type args struct {
		query interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "test basic",
			args: args{query: struct {
				Page     int `query:"page"`
				PageSize int `query:"pageSize"`
			}{
				Page:     1,
				PageSize: 20,
			}},
			want: map[string]interface{}{
				"page":     1,
				"pageSize": 20,
			},
		},
		{
			name: "test required",
			args: args{query: struct {
				Name string `query:",required"`
			}{}},
			wantErr: true,
		},
		{
			name: "test default",
			args: args{query: struct {
				Name string `query:",default=abc"`
			}{}},
			want: map[string]interface{}{
				"name": "abc",
			},
			wantErr: false,
		},
		{
			name: "test omit",
			args: args{query: struct {
				Name string `query:"name,omitempty"`
			}{
				Name: "",
			}},
			want: map[string]interface{}{},
		},
		{
			name: "test private member",
			args: args{query: struct {
				name string `query:"name"`
			}{
				name: "abc",
			}},
			want: map[string]interface{}{},
		}, {
			name: "test struct ptr",
			args: args{query: &struct {
				Name string `query:""`
			}{
				Name: "abc",
			}},
			want: map[string]interface{}{
				"name": "abc",
			},
		},
		{
			name: "test map",
			args: args{query: map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			}},
			want: map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			},
		},
		{
			name: "test string map",
			args: args{query: map[string]string{
				"page":     "1",
				"pageSize": "10",
			}},
			want: map[string]interface{}{
				"page":     "1",
				"pageSize": "10",
			},
		},
		{
			name: "test map ptr",
			args: args{query: &map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			}},
			want: map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			},
		},
		{
			name:    "test invalid type",
			args:    args{query: "abc"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formToMap(tt.args.query, queryTagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("formToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("formToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getQueries(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test basic",
			args: args{m: map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			}},
			want: "page=1&pageSize=10",
		},
		{
			name: "test escape",
			args: args{m: map[string]interface{}{
				"visitor": "Bob John",
			}},
			want: "visitor=Bob+John",
		},
		{
			name: "test empty",
			args: args{m: nil},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeForm(tt.args.m); got != tt.want {
				t.Errorf("encodeForm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUrl(t *testing.T) {
	type args struct {
		req *Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test basic",
			args: args{req: Req().HostAndPort("localhost", 8080).WithPath("abc")},
			want: "http://localhost:8080/abc",
		},
		{
			name: "test query",
			args: args{req: Req().Host("localhost").WithPath("abc").WithQueries(map[string]interface{}{
				"page":     1,
				"pageSize": 10,
			})},
			want: "http://localhost:80/abc?page=1&pageSize=10",
		},
		{
			name:    "test error",
			args:    args{req: Req().Host("localhost").WithPath("abc").WithQueries("123")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUrl(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrl()  error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doHttpReq(t *testing.T) {
	type args struct {
		method string
		req    *Request
		rsp    Response
	}

	type reqBody struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name    string
		args    args
		want    Response
		wantErr bool
	}{
		{
			name: "test get",
			args: args{
				req:    Req().HostAndPort("localhost", 8080).WithPath("echo"),
				method: http.MethodGet,
				rsp:    &DefaultResponse{},
			},
			want: &DefaultResponse{},
		},
		{
			name: "test post",
			args: args{
				req:    Req().HostAndPort("localhost", 8080).WithPath("echo").WithBody(reqBody{Name: "abc"}),
				method: http.MethodPost,
				rsp:    &DefaultResponse{Data: &reqBody{}},
			},
			want: &DefaultResponse{Data: &reqBody{
				Name: "abc",
			}},
		},
		{
			name: "test []byte body",
			args: args{
				method: http.MethodPost,
				req:    Req().HostAndPort("localhost", 8080).WithPath("echo").WithBody([]byte(`{"name": "abc"}`)),
				rsp:    &DefaultResponse{Data: &reqBody{}},
			},
			want: &DefaultResponse{Data: &reqBody{
				Name: "abc",
			}},
		},
		{
			name: "test malformed json",
			args: args{
				method: http.MethodPost,
				req:    Req().HostAndPort("localhost", 8080).WithPath("echo").WithBody("name=abc"),
				rsp:    &DefaultResponse{Data: &reqBody{}},
			},
			wantErr: true,
		},
		{
			name: "test timeout",
			args: args{
				method: http.MethodGet,
				req:    Req().HostAndPort("localhost", 8080).WithPath("echo").WithTimeout(time.Second * 2),
				rsp:    &DefaultResponse{},
			},
		},
		{
			name: "test header",
			args: args{
				method: http.MethodGet,
				req: Req().HostAndPort("localhost", 8080).WithPath("echo").WithHeaders(map[string]string{
					"Content-Type": "application/json",
				}),
				rsp: &DefaultResponse{},
			},
		},
	}
	go func() {
		http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(r.Body)
				readAll, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte(err.Error()))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(readAll)
			case http.MethodGet:
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status": 200}`))
			}
		})

		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("echo server start failed ", err)
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.req.ctx.Method = tt.args.method
			tt.args.req.ctx.Response = tt.args.rsp
			doHttpReq(tt.args.req.ctx)
			if (tt.args.req.ctx.Response.Error() != nil) != tt.wantErr {
				t.Errorf("doHttpReq() error = %v, wantErr %v", tt.args.req.ctx.Response.Error(), tt.wantErr)
				return
			}

		})
	}
}
