package binding

import (
	"net/http"
)

type Binding interface {
	Bind(*http.Request, interface{}) error
}

type Func func(*http.Request, interface{}) error

func (f Func) Bind(r *http.Request, v interface{}) error {
	return f(r, v)
}
