package request

import "net/http"

// AddHeader Add header to Request
func (options) AddHeader(k, v string) Option {
	return func(r *Request) {
		if r.requestHeader == nil {
			r.requestHeader = http.Header{}
		}
		r.requestHeader.Add(k, v)
	}
}

// SetHeader Set header to Request
func (options) SetHeader(k, v string) Option {
	return func(r *Request) {
		if r.requestHeader == nil {
			r.requestHeader = http.Header{}
		}
		r.requestHeader.Set(k, v)
	}
}

// ReplaceHeader can replace whole request header by f's returned value
// can used to Del header
func (options) ReplaceHeader(f func(oldHeader http.Header) http.Header) Option {
	return func(r *Request) {
		r.requestHeader = f(r.requestHeader)
	}
}

// GetResponseHeader get response's header, add it into given header
// happen before CheckResponseBeforeUnmarshalFunc be called
func (options) GetResponseHeader(header http.Header) Option {
	if header == nil {
		panic("header must not nil")
	}
	return func(r *Request) {
		r.responseHeaders = append(r.responseHeaders, header)
	}
}
