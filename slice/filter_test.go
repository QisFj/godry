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
	filter := func(index int) bool {
		return es[index].V1 >= es[index].V2 && es[index].V1%es[index].V2 == 0
	}
	t.Run("normal", func(t *testing.T) {
		Filter(&es, filter)
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
	t.Run("interface", func(t *testing.T) {
		var i interface{} = &es
		Filter(i, filter)
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
	t.Run("panic", func(t *testing.T) {
		defer func() {
			if p := recover(); p != nil {
				t.Logf("Expected Panic: %#v", p)
				return
			}
			t.Error("Expected Panic, But Not")
		}()
		Filter(es, filter)
	})
	t.Run("nil", func(t *testing.T) {
		Filter(nil, filter)
	})
	t.Run("nil pointer", func(t *testing.T) {
		var p *[]E
		Filter(p, filter)
	})
	t.Run("nil interface", func(t *testing.T) {
		var i interface{}
		Filter(i, filter)
	})

}
