package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	client                      *http.Client
	timeout                     time.Duration
	header                      http.Header
	method                      string
	log                         Log
	statusCodeCheck             func(int) error
	responseCheck               func([]byte) error
	responseCheckAfterUnmarshal func(interface{}) error
}

// Do request with default context
func (r Request) Do(u string, v url.Values, req, resp interface{}) error {
	return r.DoCtx(context.Background(), u, v, req, resp)
}

// Do request
func (r Request) DoCtx(ctx context.Context, u string, v url.Values, req, resp interface{}) error {
	var (
		reqBytes  []byte
		reqReader io.Reader
		httpReq   *http.Request

		err error
	)
	if req != nil {
		reqBytes, err = json.Marshal(req)
		if err != nil {
			return fmt.Errorf("http|reqeust marshal error: %w", err)
		}
		reqReader = bytes.NewBuffer(reqBytes)
	}
	httpReq, err = http.NewRequestWithContext(ctx, r.method, u, reqReader)
	if err != nil {
		return fmt.Errorf("http|new request error: %w", err)
	}
	if v != nil {
		for key, values := range httpReq.URL.Query() {
			for _, value := range values {
				v.Add(key, value)
			}
		}
		httpReq.URL.RawQuery = v.Encode()
	}
	for key, values := range r.header {
		for _, value := range values {
			httpReq.Header.Add(key, value)
		}
	}
	if r.client == nil {
		r.client = http.DefaultClient
	}
	r.client.Timeout = r.timeout
	if r.log.Logger != nil && r.log.URL {
		r.log.Logger("http|%s %s", r.method, httpReq.URL)
	}

	if r.log.Logger != nil && r.log.RequestBody && len(reqBytes) != 0 {
		if len(reqBytes) < r.log.RequestBodyLimit {
			r.log.Logger("http|request body: %s", string(reqBytes))
		} else {
			r.log.Logger("http|request body: %s ...", string(reqBytes[:r.log.RequestBodyLimit]))
		}
	}
	var httpResp *http.Response
	httpResp, err = r.client.Do(httpReq)
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
		if len(respBytes) < r.log.ResponseBodyLimit {
			r.log.Logger("http|response body: %s", string(respBytes))
		} else {
			r.log.Logger("http|response body: %s ...", string(respBytes[:r.log.ResponseBodyLimit]))
		}
	}
	if r.responseCheck != nil {
		if err = r.responseCheck(respBytes); err != nil {
			return fmt.Errorf("http|response response check error: %w", err)
		}
	}
	if resp != nil { // ignore response, if resp == nil
		if err = json.Unmarshal(respBytes, resp); err != nil {
			return fmt.Errorf("http|response unmarshal error: %w", err)
		}
		if r.responseCheckAfterUnmarshal != nil {
			if err = r.responseCheckAfterUnmarshal(resp); err != nil {
				return fmt.Errorf("http|response check after unmarshal error: %w", err)
			}
		}
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
