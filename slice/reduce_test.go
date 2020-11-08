package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	t.Run("len", func(t *testing.T) {
		length := Reduce([]int{1, 2, 3, 4}, 0, func(before interface{}, i int, v interface{}) interface{} {
			return before.(int) + 1
		}).(int)
		require.Equal(t, 4, length)
	})
	t.Run("sum", func(t *testing.T) {
		sum := Reduce([]int{1, 2, 3, 4}, 0, func(before interface{}, i int, v interface{}) interface{} {
			return before.(int) + v.(int)
		}).(int)
		require.Equal(t, 10, sum)
	})
}
