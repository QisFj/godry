package informer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_infiniteRingBuffer(t *testing.T) {
	irb := newInfiniteRingBuffer[int](1)

	irb.append(1)
	irb.append(2)
	irb.append(3)
	irb.append(4)
	irb.append(5)
	var (
		v  int
		ok bool
	)
	v, _ = irb.pop()
	require.Equal(t, 1, v)
	v, _ = irb.pop()
	require.Equal(t, 2, v)
	v, _ = irb.pop()
	require.Equal(t, 3, v)
	v, _ = irb.pop()
	require.Equal(t, 4, v)
	v, _ = irb.pop()
	require.Equal(t, 5, v)
	_, ok = irb.pop()
	require.False(t, ok)

	irb.append(6)
	irb.append(7)
	irb.append(8)

	v, _ = irb.pop()
	require.Equal(t, 6, v)
	v, _ = irb.pop()
	require.Equal(t, 7, v)
	v, _ = irb.pop()
	require.Equal(t, 8, v)
}

func Test_bufferedChannel(t *testing.T) {
	bc := NewBufferedChannel[int](0)

	input := bc.Source()
	output := bc.Sink()

	go bc.Run()

	for i := 0; i < 10; i++ {
		input <- i
	}

	for i := 0; i < 5; i++ {
		v, ok := <-output
		require.True(t, ok)
		require.Equal(t, i, v)
	}

	close(input)

	_, ok := <-output
	require.False(t, ok)
}
