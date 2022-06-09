package set

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type SortableSet[V constraints.Ordered] map[V]struct{}

// OrderedList return an ordered list
// can use AscLess, DescLess as less function
func (s SortableSet[V]) OrderedList(less func(v1, v2 V) bool) []V {
	list := Set[V](s).List()
	sort.Slice(list, func(i, j int) bool {
		return less(list[i], list[j])
	})
	return list
}

func AscLess[V constraints.Ordered](v1, v2 V) bool  { return v1 < v2 }
func DescLess[V constraints.Ordered](v1, v2 V) bool { return v1 > v2 }
