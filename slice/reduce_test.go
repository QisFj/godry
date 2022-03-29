package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	t.Run("len", func(t *testing.T) {
		length := Reduce([]int{1, 2, 3, 4}, 0, func(before, i, v int) int {
			return before + 1
		})
		require.Equal(t, 4, length)
	})
	t.Run("sum", func(t *testing.T) {
		sum := Reduce([]int{1, 2, 3, 4}, 0, func(before, i, v int) int {
			return before + v
		})
		require.Equal(t, 10, sum)
	})
}
