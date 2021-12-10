package set

import (
	"encoding/json"
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

func TestIntJSON(t *testing.T) {
	MarshalOrder = 1 // asc
	t.Run("marshal", func(t *testing.T) {
		requireJSON := func(t *testing.T, value interface{}, expJSON string) {
			actJSON, err := json.Marshal(value)
			require.NoError(t, err)
			require.JSONEq(t, expJSON, string(actJSON))
		}
		t.Run("nil", func(t *testing.T) {
			requireJSON(t, Int(nil), `null`)
		})
		t.Run("empty", func(t *testing.T) {
			requireJSON(t, Int{}, `[]`)
		})
		t.Run("non-empty", func(t *testing.T) {
			requireJSON(t, FromInts(1), `[1]`)
		})
		t.Run("duplicates", func(t *testing.T) {
			requireJSON(t, FromInts(1, 1), `[1]`)
		})
		t.Run("order", func(t *testing.T) {
			requireJSON(t, FromInts(4, 3, 2, 1), `[1, 2, 3, 4]`)
		})
	})
	t.Run("unmarshal", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			var set Int
			require.NoError(t, json.Unmarshal([]byte(`null`), &set))
			require.Equal(t, Int(nil), set)
		})
		t.Run("empty", func(t *testing.T) {
			var set Int
			require.NoError(t, json.Unmarshal([]byte(`[]`), &set))
			require.Equal(t, Int{}, set)
		})
		t.Run("non-empty", func(t *testing.T) {
			var set Int
			require.NoError(t, json.Unmarshal([]byte(`[1]`), &set))
			require.Equal(t, FromInts(1), set)
		})
		t.Run("duplicates", func(t *testing.T) {
			var set Int
			require.NoError(t, json.Unmarshal([]byte(`[1, 1]`), &set))
			require.Equal(t, FromInts(1), set)
		})
		t.Run("order", func(t *testing.T) {
			var set Int
			require.NoError(t, json.Unmarshal([]byte(`[4, 3, 2, 1]`), &set))
			require.Equal(t, FromInts(1, 2, 3, 4), set)
		})
		t.Run("embed", func(t *testing.T) {
			var p struct{ Set Int }
			require.NoError(t, json.Unmarshal([]byte(`{ "Set": [1, 2, 3] }`), &p))
			require.Equal(t, FromInts(1, 2, 3), p.Set)
		})
		t.Run("embed not exist", func(t *testing.T) {
			var p struct{ Set Int }
			require.NoError(t, json.Unmarshal([]byte(`{}`), &p))
			require.Equal(t, Int(nil), p.Set)
		})
	})
}
