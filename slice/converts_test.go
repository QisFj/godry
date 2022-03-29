package slice

import (
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToKVs(t *testing.T) {
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	kvs := KVsOfMap(m)
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].Key < kvs[j].Key })
	require.Equal(t, []KV[int, string]{
		{Key: 1, Value: "1"},
		{Key: 2, Value: "2"},
		{Key: 3, Value: "3"},
	}, kvs)
}

func TestToMap(t *testing.T) {
	require.Equal(t, map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}, ToMap([]int{1, 2, 3},
		strconv.Itoa,
		func(i int) int { return i },
	))
}

func TestToMapItself(t *testing.T) {
	require.Equal(t, map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}, ToMapItself([]int{1, 2, 3},
		strconv.Itoa,
	))
}
