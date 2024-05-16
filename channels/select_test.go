package channels

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	N, M := 100, 100
	chs := make([]chan int, N)
	for i := range chs {
		chs[i] = make(chan int)
	}
	var closeCnt atomic.Uint32
	for _i, _ch := range chs {
		i, ch := _i, _ch // no need anymore since go1.22
		go func() {
			for j := 0; j < M; j++ {
				ch <- i*M + j
			}
			close(ch)
			closeCnt.Add(1)
		}()
	}
	vised := make([]bool, N*M)
	visCnt := 0
	for {
		_, v, ok := Read(chs)
		if !ok {
			if closeCnt.Load() == uint32(N) {
				break
			}
			continue
		}
		require.True(t, v >= 0 && v < N*M, "invalid value %d", v)
		require.False(t, vised[v], "duplicated value %d", v)
		vised[v] = true
		visCnt++
	}
	require.Equal(t, N*M, visCnt, "not all values are visited")
}
