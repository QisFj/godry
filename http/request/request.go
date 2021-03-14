package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// zero value is not ready for use, call New
type Request struct {
	client                            *http.Client
	header                            http.Header
	method                            string
	log                               Log
	codec                             Codec
	checkResponseBeforeUnmarshalFuncs CheckResponseBeforeUnmarshalFuncs
	checkResponseAfterUnmarshalFuncs  CheckResponseAfterUnmarshalFuncs
}

func New() *Request {
	return new(Request).With(defaultOptions...)
}

var defaultOptions = []Option{
	Options.Client(http.DefaultClient),
	Options.Codec(CodecFuncs{
		MarshalFunc:   json.Marshal,
		UnmarshalFunc: json.Unmarshal,
	}),
}

// Do request with default context
func (r Request) Do(u string, v url.Values, req, resp interface{}) error {
	return r.DoCtx(context.Background(), u, v, req, resp)
}

// Do request
func (r Request) DoCtx(ctx context.Context, u string, v url.Values, reqObj, respObj interface{}) (err error) {
	var (
		reqBody   []byte
		reqReader io.Reader
	)
	if reqObj != nil {
		if reqBody, err = r.codec.Marshal(reqObj); err != nil {
			return fmt.Errorf("request|marshal request object error: %w", err)
		}
		reqReader = bytes.NewBuffer(reqBody)
	}
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, r.method, u, reqReader); err != nil {
		return fmt.Errorf("request|new http request error: %w", err)
	}
	if v != nil {
		for key, values := range req.URL.Query() {
			for _, value := range values {
				v.Add(key, value)
			}
		}
		req.URL.RawQuery = v.Encode()
	}
	for key, values := range r.header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	if r.client == nil {
		r.client = http.DefaultClient
	}
	r.log.LogURL(req.Method, req.URL)
	r.log.LogRequestBody(reqBody)

	var resp *http.Response
	if resp, err = r.client.Do(req); err != nil {
		return fmt.Errorf("requset|do http request error: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var respBody []byte
	if respBody, err = ioutil.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("requset|read response body error: %w", err)
	}
	r.log.LogResponseBody(reqBody)
	if err = r.checkResponseBeforeUnmarshalFuncs.Check(resp.StatusCode, respBody); err != nil {
		return fmt.Errorf("request|check response before unmarshal failed: %w", err)
	}

	if respObj != nil { // ignore response, if respObj == nil
		if err = r.codec.Unmarshal(respBody, respObj); err != nil {
			return fmt.Errorf("request|response body unmarshal error: %w", err)
		}
		if err = r.checkResponseAfterUnmarshalFuncs.Check(resp.StatusCode, respObj); err != nil {
			return fmt.Errorf("request|check response after unmarshal failed: %w", err)
		}
	}
	return nil
}

// important: With changed the original Request
func (r *Request) With(opts ...Option) *Request {
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

func (r Request) Clone() *Request {
	r.header = r.header.Clone()
	return &r
}