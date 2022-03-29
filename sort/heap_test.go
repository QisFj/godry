package sort

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeap(t *testing.T) {
	randArray := func() []string {
		n := 1000 // 000 -> 999
		array := make([]string, n)
		for i := 0; i < n; i++ {
			array[i] = fmt.Sprintf("%03d", i)
		}
		for i := 0; i < n; i++ {
			j := i + rand.Intn(n-i) // nolint: gosec
			array[i], array[j] = array[j], array[i]
		}
		return array
	}
	t.Run("normal", func(t *testing.T) {
		heap := Heap[string]{
			Less: func(o1, o2 string) bool {
				return o1 < o2
			},
		}
		for _, v := range randArray() {
			heap.Append(v)
		}
		for i, v := range heap.Dump() {
			require.Equal(t, fmt.Sprintf("%03d", i), v)
		}
	})
	t.Run("sized", func(t *testing.T) {
		heap := Heap[string]{
			Less: func(o1, o2 string) bool {
				return o1 < o2
			},
			Size: 10,
		}
		for _, v := range randArray() {
			heap.Append(v)
		}
		require.Equal(t, []interface{}{
			"000", "001", "002", "003", "004", "005", "006", "007", "008", "009",
		}, heap.Dump())
	})
}
