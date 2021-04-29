package request

import "net/http"

type RequestHookFunc func(req *http.Request) error

type RequestHookFuncs []RequestHookFunc

func (fs RequestHookFuncs) Hook(req *http.Request) error {
	for _, f := range fs {
		if err := f(req); err != nil {
			return err
		}
	}
	return nil
}

type ResponseHookFunc func(resp *http.Response) error

type ResponseHookFuncs []ResponseHookFunc

func (fs ResponseHookFuncs) Hook(resp *http.Response) error {
	for _, f := range fs {
		if err := f(resp); err != nil {
			return err
		}
	}
	return nil
}

// happen before log, after other request option
func (options) HookRequest(f RequestHookFunc) Option {
	return func(r *Request) {
		r.requestHookFuncs = append(r.requestHookFuncs, f)
	}
}

// happen before other response option
func (options) HookResponse(f ResponseHookFunc) Option {
	return func(r *Request) {
		r.responseHookFuncs = append(r.responseHookFuncs, f)
	}
}
