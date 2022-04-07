package set

import (
	"reflect"
	"testing"
)

func TestT(t *testing.T) {
	set1, set2 := NewT(reflect.TypeOf(0)), NewT(reflect.TypeOf(0))

	set1.Add(1, 2, 3, 4, 5)
	set2.Add(3, 4, 5, 6, 7)

	t.Run("union", func(t *testing.T) {
		list := TUnion(set1, set2).List().([]int)
		requireEqualAfterSort(t, []int{1, 2, 3, 4, 5, 6, 7}, list)
	})

	t.Run("intersect", func(t *testing.T) {
		list := TIntersect(set1, set2).List().([]int)
		requireEqualAfterSort(t, []int{3, 4, 5}, list)
	})

	t.Run("subtract", func(t *testing.T) {
		list := TSubtract(set1, set2).List().([]int)
		requireEqualAfterSort(t, []int{1, 2}, list)
	})

	t.Run("diff", func(t *testing.T) {
		both, only1, only2 := TDiff(set1, set2)
		requireEqualAfterSort(t, []int{3, 4, 5}, both.List().([]int))
		requireEqualAfterSort(t, []int{1, 2}, only1.List().([]int))
		requireEqualAfterSort(t, []int{6, 7}, only2.List().([]int))
	})

}
