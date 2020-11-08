package slice

import (
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeysOfMap(t *testing.T) {
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	list := KeysOfMap(m).([]int)
	sort.Ints(list)
	require.Equal(t, []int{1, 2, 3}, list)
}

func TestValuesOfMap(t *testing.T) {
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	list := ValuesOfMap(m).([]string)
	sort.Strings(list)
	require.Equal(t, []string{"1", "2", "3"}, list)
}

func TestToKVs(t *testing.T) {
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	kvs := KVsOfMap(m)
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Key.(int) < kvs[j].Key.(int)
	})
	require.Equal(t, []KV{
		{Key: 1, Value: "1"},
		{Key: 2, Value: "2"},
		{Key: 3, Value: "3"},
	}, kvs)
}

func TestToMap(t *testing.T) {
	require.Equal(t, map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}, ToMap([]int{1, 2, 3},
		reflect.TypeOf(0), func(i int, v interface{}) interface{} {
			return v
		},
		reflect.TypeOf(""), func(i int, v interface{}) interface{} {
			return strconv.Itoa(v.(int))
		},
	))
}
