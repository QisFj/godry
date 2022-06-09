package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortableSet_OrderedList(t *testing.T) {
	s := SortableSet[int](Of(1, 1, 4, 2, 3, 1, 4))

	ascList := s.OrderedList(AscLess[int])
	require.Equal(t, []int{1, 2, 3, 4}, ascList)

	descList := s.OrderedList(DescLess[int])
	require.Equal(t, []int{4, 3, 2, 1}, descList)
}
