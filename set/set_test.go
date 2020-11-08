package set

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func testEQAfterSorted(t *testing.T, exp, act []int) {
	sort.Ints(act)
	require.Equal(t, exp, act)
}

func TestT(t *testing.T) {
	set1, set2 := NewT(reflect.TypeOf(0)), NewT(reflect.TypeOf(0))
	set1.Add(1, 2, 3, 4, 5)
	set2.Add(3, 4, 5, 6, 7)
	list := TMerge(set1, set2).List().([]int)
	testEQAfterSorted(t, []int{1, 2, 3, 4, 5, 6, 7}, list)
	both, only1, only2 := TDiff(set1, set2)
	testEQAfterSorted(t, []int{3, 4, 5}, both.List().([]int))
	testEQAfterSorted(t, []int{1, 2}, only1.List().([]int))
	testEQAfterSorted(t, []int{6, 7}, only2.List().([]int))
}

func TestInt(t *testing.T) {
	set1, set2 := Int{}, Int{}
	set1.Add(1, 2, 3, 4, 5)
	set2.Add(3, 4, 5, 6, 7)
	list := IntMerge(set1, set2).List()
	testEQAfterSorted(t, []int{1, 2, 3, 4, 5, 6, 7}, list)
	both, only1, only2 := IntDiff(set1, set2)
	testEQAfterSorted(t, []int{3, 4, 5}, both.List())
	testEQAfterSorted(t, []int{1, 2}, only1.List())
	testEQAfterSorted(t, []int{6, 7}, only2.List())
}
