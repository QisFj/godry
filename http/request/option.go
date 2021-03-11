package request

import (
	"net/http"
)

type Option func(r *Request)

//  option factory
type options struct{}

var Options options

func (options) Header(k, v string) Option {
	return func(r *Request) {
		if r.header == nil {
			r.header = http.Header{}
		}
		r.header.Add(k, v)
	}
}

func (options) Method(method string) Option {
	return func(r *Request) {
		r.method = method
	}
}

func (options) CheckResponseBeforeUnmarshal(f CheckResponseBeforeUnmarshalFunc) Option {
	return func(r *Request) {
		r.checkResponseBeforeUnmarshalFuncs = append(r.checkResponseBeforeUnmarshalFuncs, f)
	}
}

func (options) CheckResponseAfterUnmarshal(f CheckResponseAfterUnmarshalFunc) Option {
	return func(r *Request) {
		r.checkResponseAfterUnmarshalFuncs = append(r.checkResponseAfterUnmarshalFuncs, f)
	}
}

func (options) Log(log Log) Option {
	return func(r *Request) {
		r.log = log
	}
}

func (options) Client(client *http.Client) Option {
	return func(r *Request) {
		r.client = client
	}
}

func (options) Codec(codec Codec) Option {
	return func(r *Request) {
		r.codec = codec
	}
}

// patch marshall only
func (options) Marshaller(m Marshaller) Option {
	return func(r *Request) {
		r.codec = CodecFuncs{
			MarshalFunc:   MarshalFuncOf(m),
			UnmarshalFunc: UnmarshalFuncOf(r.codec),
		}
	}
}

// path unmarshal only
func (options) Unmarshaller(u Unmarshaller) Option {
	return func(r *Request) {
		r.codec = CodecFuncs{
			MarshalFunc:   MarshalFuncOf(r.codec),
			UnmarshalFunc: UnmarshalFuncOf(u),
		}
	}
}
