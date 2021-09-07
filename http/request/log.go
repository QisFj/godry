package request

import (
	"context"
	"net/url"
)

type Log struct {
	Logger            func(format string, v ...interface{})
	CtxLogger         func(ctx context.Context, format string, v ...interface{}) // if this has been set, Logger make no sense
	URL               bool
	RequestBody       bool
	RequestBodyLimit  int // 0 means no limit
	ResponseBody      bool
	ResponseBodyLimit int // 0 means no limit
}

func (log Log) log(ctx context.Context, format string, v ...interface{}) {
	if log.Logger == nil && log.CtxLogger == nil {
		return
	}
	if log.CtxLogger != nil {
		log.CtxLogger(ctx, format, v...)
	} else {
		log.Logger(format, v...)
	}
}

func (log Log) LogURL(ctx context.Context, m string, u *url.URL) {
	if !log.URL {
		return
	}
	log.log(ctx, "request|%s %s", m, u)
}

func (log Log) LogRequestBody(ctx context.Context, body []byte) {
	if !log.RequestBody {
		return
	}
	log.log(ctx, "request|request body: %s", logLimit(body, log.RequestBodyLimit))
}

func (log Log) LogResponseBody(ctx context.Context, body []byte) {
	if !log.ResponseBody {
		return
	}
	log.log(ctx, "request|response body: %s", logLimit(body, log.ResponseBodyLimit))
}

func logLimit(data []byte, limit int) string {
	if limit == 0 || len(data) < limit {
		return string(data)
	}
	return string(data[:limit]) + " ..."
}
