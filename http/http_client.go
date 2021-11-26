package remote

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	schemeHttp  = "http"
	schemeHttps = "https"
)

func doHttpReq(ctx *Context) {
	req := ctx.Request

	rsp := ctx.Response
	errHandle := func(err error) {
		rsp.ErrorSave(err)
	}

	contentType := req.Headers[ContentTypeHeader]
	if contentType == "" {
		contentType = ContentTypeJson
	}

	acceptType := req.Headers[AcceptTypeHeader]
	if acceptType == "" {
		acceptType = ContentTypeJson
	}

	contentTypeResolver, ok := contentTypeRegistry[contentType]
	if !ok {
		errHandle(fmt.Errorf("unrecoginzed content type %s", contentType))
		return
	}

	acceptTypeResolver, ok := contentTypeRegistry[acceptType]
	if !ok {
		errHandle(fmt.Errorf("unrecognized accept type %s", acceptType))
		return
	}

	reqUrl, err := GetUrl(req)

	if err != nil {
		errHandle(err)
		return
	}

	var body io.Reader
	if req.Body != nil {
		if r, ok := req.Body.(io.Reader); !ok {
			var bodyBytes []byte
			var marshalErr error

			if b, ok := req.Body.([]byte); ok {
				bodyBytes = b
			} else {
				bodyBytes, marshalErr = contentTypeResolver.Marshal(req.Body)
			}
			if marshalErr != nil {
				errHandle(marshalErr)
				return
			}
			body = bytes.NewReader(bodyBytes)
		} else {
			body = r
		}
	}
	httpRequest, err := http.NewRequestWithContext(ctx.Context, ctx.Method, reqUrl, body)

	if err != nil {
		errHandle(err)
		return
	}

	if req.Headers != nil {
		for key, value := range req.Headers {
			httpRequest.Header.Add(key, value)
		}
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		errHandle(err)
		return
	}
	rsp.SetRaw(httpResponse)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResponse.Body)
	read, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		errHandle(err)
		return
	}
	if len(read) > 0 {
		if err = acceptTypeResolver.Unmarshal(read, rsp); err != nil {
			errHandle(err)
		}
	}
}

func GetUrl(req *Request) (string, error) {
	if req.ServicePort == 0 {
		if req.Secure {
			req.ServicePort = 443
		} else {
			req.ServicePort = 80
		}
	}
	reqUrl := url.URL{
		Host:   fmt.Sprintf("%s:%d", req.ServiceName, req.ServicePort),
		Path:   req.Path,
		Scheme: schemeHttp,
	}
	if req.ServicePort == 80 {
		reqUrl = url.URL{
			Host:   fmt.Sprintf("%s", req.ServiceName),
			Path:   req.Path,
			Scheme: schemeHttp,
		}
	}
	if req.Secure {
		reqUrl.Scheme = schemeHttps
	}
	if req.Query != nil {
		queriesMap, err := formToMap(req.Query, queryTagName)
		if err != nil {
			return "", err
		}
		reqUrl.RawQuery = encodeForm(queriesMap)
	}
	return reqUrl.String(), nil
}
