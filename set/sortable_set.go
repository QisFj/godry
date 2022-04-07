package set

import "golang.org/x/exp/constraints"

type SortableSet[V constraints.Ordered] map[V]struct{}
