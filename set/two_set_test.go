package set

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func requireEqualAfterSort[V constraints.Ordered](t *testing.T, exp, act []V) {
	// no need to sort exp
	slices.Sort(act)
	require.Equal(t, exp, act)
}

func Test(t *testing.T) {
	set1 := Of(1, 2, 3, 4, 5)
	set2 := Of(3, 4, 5, 6, 7)

	t.Run("union", func(t *testing.T) {
		list := Union(set1, set2).List()
		requireEqualAfterSort(t, []int{1, 2, 3, 4, 5, 6, 7}, list)
	})

	t.Run("intersect", func(t *testing.T) {
		list := Intersect(set1, set2).List()
		requireEqualAfterSort(t, []int{3, 4, 5}, list)
	})

	t.Run("subtract", func(t *testing.T) {
		list := Subtract(set1, set2).List()
		requireEqualAfterSort(t, []int{1, 2}, list)
	})

	t.Run("diff", func(t *testing.T) {
		both, only1, only2 := Diff(set1, set2)
		requireEqualAfterSort(t, []int{3, 4, 5}, both.List())
		requireEqualAfterSort(t, []int{1, 2}, only1.List())
		requireEqualAfterSort(t, []int{6, 7}, only2.List())
	})

}
