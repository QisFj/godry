package http

import (
	"errors"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestOpt(t *testing.T) {
	t.Run("Timeout", func(t *testing.T) {
		timeout := time.Second * time.Duration(rand.Intn(5))
		r := new(Request).With(Timeout(timeout))
		// Verify
		require.Equal(t, timeout, r.timeout)
	})

	t.Run("Header", func(t *testing.T) {
		header := Header{Key: "Header-Key", Value: "Header-Value"}
		r := new(Request).With(header)
		r.With(header)
		// Verify
		require.Equal(t, http.Header{"Header-Key": []string{"Header-Value", "Header-Value"}}, r.header)
	})
	t.Run("POST", func(t *testing.T) {
		r := new(Request).With(POST{})
		// Verify
		require.Equal(t, "POST", r.method)
	})

	t.Run("StatusCodeCheck", func(t *testing.T) {
		t.Run("Equal", func(t *testing.T) {
			r := new(Request).With(StatusCodeCheck{
				Equal: 200,
				CheckFunc: func(statusCode int) error {
					return errors.New("no error ")
				},
			})
			// Verify
			require.NoError(t, r.statusCodeCheck(200))
			require.NotNil(t, r.statusCodeCheck(201))
		})
		t.Run("CheckFunc", func(t *testing.T) {
			r := new(Request).With(StatusCodeCheck{
				CheckFunc: func(statusCode int) error {
					if statusCode >= 200 && statusCode < 400 {
						return nil
					}
					return errors.New("error in func1 ")
				},
			}, StatusCodeCheck{
				CheckFunc: func(statusCode int) error {
					if statusCode >= 200 && statusCode < 300 {
						return nil
					}
					return errors.New("error in func2 ")
				},
			})
			// Verify
			require.NoError(t, r.statusCodeCheck(200))
			require.NoError(t, r.statusCodeCheck(201))
			require.EqualError(t, r.statusCodeCheck(400), "error in func1 ")
			require.EqualError(t, r.statusCodeCheck(300), "error in func2 ")
		})
	})
	t.Run("ResponseCheck", func(t *testing.T) {
		r := new(Request).With(
			ResponseCheck(func(response []byte) error {
				if string(response) == "equal" {
					return nil
				}
				return errors.New("error")
			}),
		)
		// Verify
		require.NoError(t, r.responseCheck([]byte("equal")))
		require.EqualError(t, r.responseCheck([]byte("not equal")), "error")
	})
}
