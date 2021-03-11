package request

import (
	"net/url"
)

type Log struct {
	Logger            func(format string, v ...interface{})
	URL               bool
	RequestBody       bool
	RequestBodyLimit  int // 0 means no limit
	ResponseBody      bool
	ResponseBodyLimit int // 0 means no limit
}

func (log Log) LogURL(m string, u *url.URL) {
	if log.Logger == nil && !log.URL {
		return
	}
	log.Logger("request|%s %s", m, u)
}

func (log Log) LogRequestBody(body []byte) {
	if log.Logger == nil && !log.RequestBody {
		return
	}
	log.Logger("request|request body: %s", logLimit(body, log.RequestBodyLimit))
}

func (log Log) LogResponseBody(body []byte) {
	if log.Logger == nil && !log.ResponseBody {
		return
	}
	log.Logger("request|response body: %s", logLimit(body, log.ResponseBodyLimit))
}

func logLimit(data []byte, limit int) string {
	if limit == 0 || len(data) < limit {
		return string(data)
	}
	return string(data[:limit]) + " ..."
}
