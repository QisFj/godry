package retry

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	requireResultEqual := func(exp, act Result) {
		if exp.Error == nil {
			require.NoError(t, act.Error)
		} else {
			require.EqualError(t, act.Error, exp.Error.Error())
		}
		exp.Error = nil
		act.Error = nil
		require.Equal(t, exp, act)
	}
	t.Run("normal", func(t *testing.T) {
		// execute 5 times
		retry := New(Option{
			RetryInterval: 10 * time.Millisecond,
			F: func(r *Retry) error {
				if r.RunCount() < 3 {
					// won't return error at the 4th and 5th time execute
					// error wrapper would panic at 4th execute
					return fmt.Errorf("rest call: %d", 4-r.RunCount())
				}
				return nil
			},
			ResultSize: 3,
			ErrorWrapper: ErrorWrapperChain(
				FmtErrorWrapper("wrap1: %w"),
				FmtErrorWrapper("wrap2: %w"),
				func(r *Retry, err error) error { // panic wrapper
					if r.RunCount() == 3 {
						require.NoError(t, err) // verify error is nil
						panic("panic wrapper")
					}
					return err
				},
				nil, // nil wrapper
			),
		})
		tt := time.Now()
		retry.Start()
		retry.Start() // start more than once
		retry.Wait()
		require.True(t, retry.HasDone())
		require.Greater(t, int64(time.Since(tt)), int64(40*time.Millisecond))
		requireResultEqual(Result{}, retry.Result(1))
		requireResultEqual(Result{}, retry.Result(2))
		requireResultEqual(Result{Valid: true, Error: errors.New("wrap2: wrap1: rest call: 2")}, retry.Result(3))
		requireResultEqual(Result{Valid: true, Error: errors.New("error wrapper panic; err: nil error, panic: panic wrapper")}, retry.Result(4))
		requireResultEqual(Result{Valid: true, Error: nil}, retry.Result(5))
		requireResultEqual(Result{Valid: true, Error: nil}, retry.LatestResult())
		requireResultEqual(Result{}, retry.Result(6))
	})
	t.Run("zero interval", func(t *testing.T) {
		// execute 5 times
		restCallCount := 5
		retry := New(Option{
			F: func(r *Retry) error {
				restCallCount--
				require.GreaterOrEqual(t, restCallCount, 0)
				if restCallCount == 0 {
					return nil
				}
				return fmt.Errorf("rest call: %d", restCallCount)
			},
			ResultSize: 3,
		}).Start()
		tt := time.Now()
		retry.Wait()
		require.True(t, retry.HasDone())
		require.Less(t, int64(time.Since(tt)), int64(time.Millisecond))
		requireResultEqual(Result{}, retry.Result(1))
		requireResultEqual(Result{}, retry.Result(2))
		requireResultEqual(Result{Valid: true, Error: errors.New("rest call: 2")}, retry.Result(3))
		requireResultEqual(Result{Valid: true, Error: errors.New("rest call: 1")}, retry.Result(4))
		requireResultEqual(Result{Valid: true, Error: nil}, retry.Result(5))
		requireResultEqual(Result{Valid: true, Error: nil}, retry.LatestResult())
		requireResultEqual(Result{}, retry.Result(6))
	})
	t.Run("stop", func(t *testing.T) {
		retry := New(Option{
			RetryInterval: 30 * time.Millisecond,
			F:             func(r *Retry) error { return errors.New("always error") },
		}).Start()
		require.False(t, retry.HasDone())
		retry.Stop()
		retry.Stop() // stop more than once
		require.True(t, retry.HasDone())
	})
	t.Run("panic", func(t *testing.T) {
		// execute 5 times
		restCallCount := 5
		retry := New(Option{
			RetryInterval: 10 * time.Millisecond,
			F: func(r *Retry) error {
				restCallCount--
				require.GreaterOrEqual(t, restCallCount, 0)
				if restCallCount == 0 {
					return nil
				}
				panic(fmt.Errorf("rest call: %d", restCallCount))
			},
		}).Start()
		tt := time.Now()
		retry.Wait()
		require.Greater(t, int64(time.Since(tt)), int64(40*time.Millisecond))
	})
	t.Run("stop retry with error - by special error", func(t *testing.T) {
		retry := New(Option{
			RetryInterval: 30 * time.Millisecond,
			F: func(r *Retry) error {
				return StopRetryWithError(errors.New("always error"))
			},
		}).Start().Wait()
		requireResultEqual(Result{Valid: true, Error: errors.New("always error")}, retry.Result(1))
		requireResultEqual(Result{}, retry.Result(2))
	})
	t.Run("stop retry with error - by stop", func(t *testing.T) {
		retry := New(Option{
			RetryInterval: 30 * time.Millisecond,
			F: func(r *Retry) error {
				r.Stop()
				return errors.New("always error")
			},
		}).Start().Wait()
		requireResultEqual(Result{Valid: true, Error: errors.New("always error")}, retry.Result(1))
		requireResultEqual(Result{}, retry.Result(2))
	})
}
