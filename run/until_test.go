package run

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEachUntil(t *testing.T) {
	stopCh := make(chan struct{})
	before := time.Now()
	count := 0
	go func() {
		EachUntil(func() {
			t.Logf("run at %s", time.Since(before)/time.Millisecond*time.Millisecond)
			count++
		}, 50*time.Millisecond, stopCh)
	}()
	time.Sleep(220 * time.Millisecond)
	close(stopCh)
	require.Equal(t, 4, count) // 50, 100, 150, 200
}

func TestEachUntilImmediately(t *testing.T) {
	stopCh := make(chan struct{})
	before := time.Now()
	count := 0
	go func() {
		EachUntilImmediately(func() {
			t.Logf("run at %s", time.Since(before)/time.Millisecond*time.Millisecond)
			count++
		}, 50*time.Millisecond, stopCh)
	}()
	time.Sleep(220 * time.Millisecond)
	close(stopCh)
	require.Equal(t, 5, count) // 0, 50, 100, 150, 200
}
