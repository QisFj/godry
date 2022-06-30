package concurrent

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGo(t *testing.T) {
	sum := int32(0)
	options := []GoOption{
		WithWrapper(func(f FuncWithResultMayError) FuncWithResultMayError {
			return func() (result interface{}, err error) {
				traceID := rand.Int()
				t.Logf("START [%d]", traceID)
				defer func() {
					t.Logf("STOP [%d]", traceID)
				}()
				return f()
			}
		}, true),
		WithSkipNilFunc(true),
		WithLimiter(NewLimiter(2)),
		WithWaitGroup(&sync.WaitGroup{}),
		WithRecoverPanic(func(p interface{}) {
			t.Logf("GOT PANIC: %s", p)
		}),
		WithErrorHandler(func(err error) {
			t.Logf("GOT ERROR: %s", err)
		}),
		WithResultHandler(func(result interface{}) {
			atomic.AddInt32(&sum, int32(result.(int)))
		}),
	}

	for i := 0; i < 10; i++ {
		ii := i
		Go(func() (result interface{}, err error) {
			time.Sleep(20 * time.Millisecond)
			return ii, nil
		}, options...)
	}

	Go(nil, options...)
	Go(func() (result interface{}, err error) {
		return nil, fmt.Errorf("some error")
	}, options...)
	Go(func() (result interface{}, err error) {
		panic("some panic")
	}, options...)

	require.Equal(t, int32(45), sum)
}
