package request

import "net/http"

// Set Header to Request
func (options) Header(k, v string) Option {
	return func(r *Request) {
		if r.requestHeader == nil {
			r.requestHeader = http.Header{}
		}
		r.requestHeader.Add(k, v)
	}
}

// get response's header, add it into given header
// happen before CheckResponseBeforeUnmarshalFunc be called
func (options) GetResponseHeader(header http.Header) Option {
	if header == nil {
		panic("header must not nil")
	}
	return func(r *Request) {
		r.responseHeaders = append(r.responseHeaders, header)
	}
}
