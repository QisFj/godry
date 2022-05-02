package run

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_LazyRunner(t *testing.T) {
	var result []int
	before := time.Now()
	runner := LazyRunner{
		Run: func(stopCh <-chan struct{}) {
			t.Log("Start")
			EachUntilImmediately(func() {
				ms := time.Since(before) / time.Millisecond * time.Millisecond
				t.Logf("run at %s", ms)
				result = append(result, int(ms/50/time.Millisecond))
			}, 50*time.Millisecond, stopCh)
			t.Log("Stop")
		},
		Locker: new(sync.RWMutex),
	}
	after := func(dur time.Duration) <-chan struct{} {
		ch := make(chan struct{})
		time.AfterFunc(dur, func() {
			close(ch)
		})
		return ch
	}

	runner.AddSupervisor(after(100 * time.Millisecond))
	runner.AddSupervisor(after(330 * time.Millisecond))
	time.Sleep(500 * time.Millisecond)

	time.Sleep(500 * time.Millisecond)
	runner.AddSupervisor(after(100 * time.Millisecond))
	runner.AddSupervisor(after(330 * time.Millisecond))
	time.Sleep(500 * time.Millisecond)

	require.Len(t, result, 14)
	// 1st run: 0, 50, 100, 150, 200, 250, 300
	// 2nd run: 0, 50, 100, 150, 200, 250, 300
}
