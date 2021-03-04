package concurrent

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	var count int32
	var wg sync.WaitGroup
	var limit = 4
	var limiter = NewLimiter(limit - 1)
	// it's ok to use Release without Acquire to increase limit or Acquire without Release to decrease
	// eg.
	limiter.Release()
	base := time.Now()
	function := func() {
		c := atomic.AddInt32(&count, 1)
		if c > int32(limit) {
			panic("too much running")
		}
		time.Sleep(30 * time.Millisecond)
		atomic.AddInt32(&count, -1)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// a simple ticker print log to tell how many functions are running
		// 10 job, echo cost 30 ms to run, but 4 job can run at the same time
		// the 1st  start at 5ms,  stop at 35ms
		// the 5th  start at 35ms, stop at 65ms
		// the 9th  start at 65ms, stop at 95ms
		// the 10th start at 70ms, stop at 100ms
		// ticker run 23 times, each time sleep 5ms
		for i := 0; i < 21; i++ {
			t.Logf("% 4dms run count: %d", time.Since(base).Milliseconds(), atomic.LoadInt32(&count))
			time.Sleep(5 * time.Millisecond)
		}
	}()
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Millisecond)
		wg.Add(1)
		go func() {
			defer wg.Done()
			limiter.Acquire()
			defer limiter.Release()
			function()
		}()
	}
	wg.Wait()
}
