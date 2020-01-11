package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	timeout         time.Duration
	header          http.Header
	method          string
	log             Log
	statusCodeCheck func(int) error
	responseCheck   func([]byte) error
}

func (r Request) Do(u string, v url.Values, req, resp interface{}) error {
	var reqReader io.Reader
	var reqBytes []byte
	var err error
	if req != nil {
		reqBytes, err = json.Marshal(req)
		if err != nil {
			return fmt.Errorf("http|reqeust marshal error: %w", err)
		}
		reqReader = bytes.NewBuffer(reqBytes)
	}
	var httpReq *http.Request
	httpReq, err = http.NewRequest(r.method, u, reqReader)
	if err != nil {
		return fmt.Errorf("http|new request error: %w", err)
	}
	for key, values := range httpReq.URL.Query() {
		for _, value := range values {
			v.Add(key, value)
		}
	}
	httpReq.URL.RawQuery = v.Encode()
	for k, vs := range r.header {
		for _, v := range vs {
			httpReq.Header.Add(k, v)
		}
	}
	httpClient := http.Client{
		Timeout: r.timeout,
	}
	if r.log.Logger != nil && r.log.URL {
		r.log.Logger("http|%s %s", r.method, httpReq.URL)
	}

	if r.log.Logger != nil && r.log.RequestBody && len(reqBytes) != 0 {
		if len(reqBytes) >= r.log.RequestBodyLimit {
			r.log.Logger("http|request body: %s", string(reqBytes))
		} else {
			r.log.Logger("http|request body: %s ...", string(reqBytes))
		}
	}
	var httpResp *http.Response
	httpResp, err = httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http|do request error: %w", err)
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()
	if r.statusCodeCheck != nil {
		err = r.statusCodeCheck(httpResp.StatusCode)
		if err != nil {
			return fmt.Errorf("http|response status code check error: %w", err)
		}
	}
	respBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("http|read response error: %w", err)
	}
	if r.log.Logger != nil && r.log.ResponseBody {
		if len(respBytes) >= r.log.ResponseBodyLimit {
			r.log.Logger("http|response body: %s", string(respBytes))
		} else {
			r.log.Logger("http|response body: %s ...", string(respBytes))
		}
	}
	if r.responseCheck != nil {
		err = r.responseCheck(respBytes)
		if err != nil {
			return fmt.Errorf("http|response response check error: %w", err)
		}
	}
	err = json.Unmarshal(respBytes, resp)
	if err != nil {
		return fmt.Errorf("http|response unmarshal error: %w", err)
	}
	return nil
}

func (r *Request) With(opts ...RequestOpt) *Request {
	for _, opt := range opts {
		if opt != nil {
			opt.On(r)
		}
	}
	return r
}
