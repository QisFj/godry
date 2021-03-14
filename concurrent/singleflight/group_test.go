package singleflight

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	counter := int32(0)
	count := func() (interface{}, error) {
		return int(atomic.AddInt32(&counter, 1)), nil
	}
	g := Group{}
	var (
		value interface{}
		err   error
	)
	value, err = g.Do(".", count, 50*time.Millisecond)
	require.NoError(t, err)
	require.Equal(t, 1, value)
	value, err = g.Do(".", count, 10*time.Minute) // can't reset forgetAfter
	require.NoError(t, err)
	require.Equal(t, 1, value) // still return 1
	time.Sleep(100 * time.Millisecond)
	value, err = g.Do(".", count, 0)
	require.NoError(t, err)
	require.Equal(t, 2, value) // count be called
	g.Forget(".")              // forget
	g.Forget(".")              // no-op for key not exist (include already forgotten)
	value, err = g.Do(".", count, 0)
	require.NoError(t, err)
	require.Equal(t, 3, value) // count be called

	g.Forget(".")
	wg := &sync.WaitGroup{} // safe for concurrent
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Wait()
			v, e := g.Do(".", count, 50*time.Millisecond)
			require.NoError(t, e)
			require.Equal(t, 4, v)
		}()
	}
	wg.Wait()
}
