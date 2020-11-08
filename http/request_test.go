package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	tests := []struct {
		name   string      // default "[unnamed]"
		method string      // default "GET"
		header http.Header // default nil
		params url.Values  // default url.Values{}
		req    interface{} // default nil
		resp   interface{} // default struct{}{}
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
			name:   "Header",
			header: http.Header{"Header-Key": {"Header-Value"}},
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
		if test.header == nil {
			test.header = make(http.Header)
		}
		// basic header
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
				require.Equal(t, test.header, req.Header)

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
				respBytes, err := json.Marshal(test.resp)
				require.NoError(t, err)
				_, err = rw.Write(respBytes)
				require.NoError(t, err)
			}))
			defer server.Close()
			// verify that the Client received the expected response
			_resp := pof(test.resp)
			r := &Request{}
			if test.method == http.MethodPost {
				r.With(POST{})
			}
			for k, vs := range test.header {
				for _, v := range vs {
					r.With(Header{Key: k, Value: v})
				}
			}
			for k, vs := range basicHeader {
				for _, v := range vs {
					test.header.Add(k, v)
				}
			}
			require.NoError(t, r.Do(server.URL, test.params, test.req, _resp))
			require.Equal(t, test.resp, eof(_resp))
		})
	}
}
