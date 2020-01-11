package http

import (
	"fmt"
	"net/http"
	"time"
)

type RequestOpt interface {
	internal()
	On(r *Request)
}

type optInternal struct{}

func (optInternal) internal() {}

type Timeout time.Duration

func (Timeout) internal() {}
func (opt Timeout) On(r *Request) {
	r.timeout = time.Duration(opt)
}

type Header struct {
	optInternal
	Key, Value string
}

func (opt Header) On(r *Request) {
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.header.Add(opt.Key, opt.Value)
}

type POST struct {
	optInternal
}

func (opt POST) On(r *Request) {
	r.method = http.MethodPost
}

type StatusCodeCheck struct {
	optInternal
	// Equal 和 CheckFunc 只有一个生效，Equal优先
	// 实际上可能会设置多个StatusCodeCheck的Option，按顺序执行，同时通过才通过
	Equal     int
	CheckFunc func(statusCode int) error
}

func (opt StatusCodeCheck) On(r *Request) {
	oldFunc := r.statusCodeCheck
	r.statusCodeCheck = func(statusCode int) error {
		var err error
		if oldFunc != nil {
			err = oldFunc(statusCode)
			if err != nil {
				return err
			}
		}
		if opt.Equal != 0 {
			if opt.Equal != statusCode {
				return fmt.Errorf("expected code[%d] != actually code[%d]", opt.Equal, statusCode)
			}
			return nil
		}
		if opt.CheckFunc != nil {
			err = opt.CheckFunc(statusCode)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

type ResponseCheck func(response []byte) error

func (ResponseCheck) internal() {}
func (opt ResponseCheck) On(r *Request) {
	oldFunc := r.responseCheck
	r.responseCheck = func(bytes []byte) error {
		var err error
		if oldFunc != nil {
			err = oldFunc(bytes)
			if err != nil {
				return err
			}
		}
		if opt != nil {
			err = opt(bytes)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

type Log struct {
	optInternal
	Logger            func(format string, v ...interface{})
	URL               bool
	RequestBody       bool
	RequestBodyLimit  int
	ResponseBody      bool
	ResponseBodyLimit int
}

func (opt Log) On(r *Request) {
	r.log = opt
}
