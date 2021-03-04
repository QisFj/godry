package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewIter(t *testing.T) {
	it := NewIter(Graph{
		Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}, false)
	t.Run("first loop", func(t *testing.T) {
		e := 0
		for it.Next() {
			nodes := it.Get()
			a := 0
			for _, node := range nodes {
				a = 10*a + node.(int)
			}
			require.Equal(t, e, a)
			e++
		}
		require.Equal(t, 10000, e)
	})
	t.Run("second loop", func(t *testing.T) {
		e := 0
		for it.Next() {
			nodes := it.Get()
			a := 0
			for _, node := range nodes {
				a = 10*a + node.(int)
			}
			require.Equal(t, e, a)
			e++
		}
		require.Equal(t, 10000, e)
	})
	t.Run("reverse", func(t *testing.T) {
		it.reverse = true
		e := 0
		for it.Next() {
			nodes := it.Get()
			a := 0
			for i := len(nodes) - 1; i >= 0; i-- {
				a = 10*a + nodes[i].(int)
			}
			require.Equal(t, e, a)
			e++
		}
		require.Equal(t, 10000, e)
	})
}
