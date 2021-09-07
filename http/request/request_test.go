package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		name       string      // default "[unnamed]"
		method     string      // default "GET"
		reqHeader  http.Header // default nil
		respHeader http.Header // default nil
		params     url.Values  // default url.Values{}
		req        interface{} // default nil
		resp       interface{} // default struct{}{}
	}{
		{
			name: "GET",
		},
		{
			name:   "GET with Params",
			params: url.Values{"k1": nil, "k2": []string{}, "k3": []string{"v"}, "k4": []string{"v", "v"}},
		},
		{
			// Not Good, But Allowed
			name: "GET with Request Body",
			req: struct {
				String string `json:"string"`
				Int    int    `json:"int"`
			}{
				String: "string",
				Int:    123,
			},
		},
		{
			name: "GET with Response",
			resp: struct {
				String string `json:"string"`
				Int    int    `json:"int"`
			}{
				String: "string",
				Int:    123,
			},
		},
		{
			name:   "POST",
			method: http.MethodPost,
		},
		{
			name:   "POST with Params",
			method: http.MethodPost,
			params: url.Values{"k1": nil, "k2": []string{}, "k3": []string{"v"}, "k4": []string{"v", "v"}},
		},
		{
			name:   "POST with Request Body",
			method: http.MethodPost,
			req: struct {
				String string `json:"string"`
				Int    int    `json:"int"`
			}{
				String: "string",
				Int:    123,
			},
		},
		{
			name:   "POST with Response",
			method: http.MethodPost,
			resp: struct {
				String string `json:"string"`
				Int    int    `json:"int"`
			}{
				String: "string",
				Int:    123,
			},
		},
		{
			name:      "Header",
			reqHeader: http.Header{"Header-Key": {"Header-Value"}},
		},
		{
			name:       "GetHeader",
			respHeader: http.Header{"Header-Key": {"Header-Value"}},
		},
	}
	pof := func(i interface{}) interface{} {
		return reflect.New(reflect.TypeOf(i)).Interface()
	}
	eof := func(i interface{}) interface{} {
		return reflect.ValueOf(i).Elem().Interface()
	}
	for idx, tt := range tests {
		test := tt
		if test.name == "" {
			test.name = "[unnamed]"
		}
		if test.method == "" {
			test.method = "GET"
		}
		if test.reqHeader == nil {
			test.reqHeader = make(http.Header)
		}
		// basic reqHeader
		basicHeader := http.Header{
			"Accept-Encoding": {"gzip"},
			"User-Agent":      {"Go-http-client/1.1"},
		}
		if test.params == nil {
			test.params = url.Values{}
		}
		if test.method == "POST" || test.req != nil {
			if test.req != nil {
				reqBytes, err := json.Marshal(test.req)
				require.NoError(t, err)
				basicHeader.Add("Content-Length", strconv.Itoa(len(reqBytes)))
			} else {
				basicHeader.Add("Content-Length", "0")
			}
		}
		if test.resp == nil {
			test.resp = struct{}{}
		}
		t.Run(fmt.Sprintf("%d|%s", idx, test.name), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// check if the Method meets expectations
				require.Equal(t, test.method, req.Method)

				// check if the Header meets expectations
				require.Equal(t, test.reqHeader, req.Header)

				// check if the Params meets expectations
				require.Equal(t, test.params.Encode(), req.URL.Query().Encode())

				// check if the Request Body meets expectations
				reqBytes, err := ioutil.ReadAll(req.Body)
				require.NoError(t, err)
				if test.req == nil {
					require.Len(t, reqBytes, 0)
				} else {
					_req := pof(test.req)
					require.NoError(t, json.Unmarshal(reqBytes, _req))
					require.Equal(t, test.req, eof(_req))
				}

				// send response
				for k, vs := range test.respHeader {
					for _, v := range vs {
						rw.Header().Add(k, v)
					}
				}
				respBytes, err := json.Marshal(test.resp)
				require.NoError(t, err)
				_, err = rw.Write(respBytes)
				require.NoError(t, err)
			}))
			defer server.Close()
			// verify that the Client received the expected response
			_resp := pof(test.resp)
			r := New().With(Options.Log(Log{
				Logger:            t.Logf,
				URL:               true,
				RequestBody:       true,
				RequestBodyLimit:  0,
				ResponseBody:      true,
				ResponseBodyLimit: 0,
			}))
			if test.method == http.MethodPost {
				r.With(Options.Method(http.MethodPost))
			}
			for k, vs := range test.reqHeader {
				for _, v := range vs {
					r.With(Options.AddHeader(k, v))
				}
			}
			for k, vs := range basicHeader {
				for _, v := range vs {
					test.reqHeader.Add(k, v)
				}
			}
			respHeader := http.Header{}
			r.With(Options.GetResponseHeader(respHeader))
			r.With(Options.CheckResponseBeforeUnmarshal(func(statusCode int, body []byte) error {
				require.Equal(t, http.StatusOK, statusCode)
				return nil
			}))
			r.With(Options.CheckResponseAfterUnmarshal(func(statusCode int, v interface{}) error {
				require.Equal(t, http.StatusOK, statusCode)
				require.Equal(t, test.resp, eof(v))
				return nil
			}))
			r.With(Options.HookRequest(func(ctx context.Context, req *http.Request) error {
				dumpBytes, err := httputil.DumpRequest(req, true)
				if err != nil {
					return err
				}
				t.Logf("request dump: %s", string(dumpBytes))
				return nil
			}))
			r.With(Options.HookResponse(func(ctx context.Context, resp *http.Response) error {
				dumpBytes, err := httputil.DumpResponse(resp, true)
				if err != nil {
					return err
				}
				t.Logf("response dump: %s", string(dumpBytes))
				return nil
			}))
			require.NoError(t, r.Do(server.URL, test.params, test.req, _resp))
			require.Equal(t, test.resp, eof(_resp))
			// check if the Header meets expectations
			// only check if all expected header exsit and meets expectations
			for k, vs := range test.respHeader {
				require.Equal(t, vs, respHeader.Values(k), "response header %s", k)
			}
		})
	}
}
