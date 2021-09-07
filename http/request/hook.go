package request

import (
	"context"
	"net/http"
)

type RequestHookFunc func(ctx context.Context, req *http.Request) error

type RequestHookFuncs []RequestHookFunc

func (fs RequestHookFuncs) Hook(ctx context.Context, req *http.Request) error {
	for _, f := range fs {
		if err := f(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

type ResponseHookFunc func(ctx context.Context, resp *http.Response) error

type ResponseHookFuncs []ResponseHookFunc

func (fs ResponseHookFuncs) Hook(ctx context.Context, resp *http.Response) error {
	for _, f := range fs {
		if err := f(ctx, resp); err != nil {
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
