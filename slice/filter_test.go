package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	type E struct{ V1, V2 int }
	es := []E{
		{V1: 10, V2: 1},
		{V1: 10, V2: 2},
		{V1: 10, V2: 3},
		{V1: 10, V2: 5},
		{V1: 10, V2: 10},
		{V1: 10, V2: 11},
		{V1: 15, V2: 1},
		{V1: 15, V2: 2},
		{V1: 15, V2: 3},
		{V1: 15, V2: 4},
		{V1: 15, V2: 5},
	}
	clone := Clone(es)
	filter := func(_ int, e E) bool {
		return e.V1 >= e.V2 && e.V1%e.V2 == 0
	}
	t.Run("normal", func(t *testing.T) {
		filtered := Filter(es, filter)
		require.Equal(t, clone, es) // no change
		require.Equal(t, []E{
			{V1: 10, V2: 1},
			{V1: 10, V2: 2},
			{V1: 10, V2: 5},
			{V1: 10, V2: 10},
			{V1: 15, V2: 1},
			{V1: 15, V2: 3},
			{V1: 15, V2: 5},
		}, filtered)
	})
	t.Run("nil", func(t *testing.T) {
		filtered := Filter(nil, filter)
		require.Nil(t, filtered)
	})
	t.Run("zero", func(t *testing.T) {
		filtered := Filter([]E{}, filter)
		require.Equal(t, []E{}, filtered)
	})
}

func TestFilterOn(t *testing.T) {
	type E struct{ V1, V2 int }
	es := []E{
		{V1: 10, V2: 1},
		{V1: 10, V2: 2},
		{V1: 10, V2: 3},
		{V1: 10, V2: 5},
		{V1: 10, V2: 10},
		{V1: 10, V2: 11},
		{V1: 15, V2: 1},
		{V1: 15, V2: 2},
		{V1: 15, V2: 3},
		{V1: 15, V2: 4},
		{V1: 15, V2: 5},
	}
	filter := func(_ int, e E) bool {
		return e.V1 >= e.V2 && e.V1%e.V2 == 0
	}
	t.Run("normal", func(t *testing.T) {
		FilterOn(&es, filter)
		require.Equal(t, []E{
			{V1: 10, V2: 1},
			{V1: 10, V2: 2},
			{V1: 10, V2: 5},
			{V1: 10, V2: 10},
			{V1: 15, V2: 1},
			{V1: 15, V2: 3},
			{V1: 15, V2: 5},
		}, es)
	})
	t.Run("nil pointer", func(t *testing.T) {
		FilterOn((*[]E)(nil), filter)
	})
	t.Run("nil slice", func(t *testing.T) {
		slice := []E(nil)
		FilterOn(&slice, filter)
	})
}
