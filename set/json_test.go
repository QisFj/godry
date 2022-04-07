package set

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSON(t *testing.T) {
	MarshalOrder = 1 // asc
	t.Run("normal", func(t *testing.T) {
		t.Run("marshal", func(t *testing.T) {
			requireJSON := func(t *testing.T, value interface{}, expJSON string) {
				actJSON, err := json.Marshal(value)
				require.NoError(t, err)
				require.JSONEq(t, expJSON, string(actJSON))
			}
			t.Run("nil", func(t *testing.T) {
				requireJSON(t, Set[int](nil), `null`)
			})
			t.Run("empty", func(t *testing.T) {
				requireJSON(t, Set[int]{}, `[]`)
			})
			t.Run("non-empty", func(t *testing.T) {
				requireJSON(t, Of(1), `[1]`)
			})
			t.Run("duplicates", func(t *testing.T) {
				requireJSON(t, Of(1, 1), `[1]`)
			})
		})
		t.Run("unmarshal", func(t *testing.T) {
			t.Run("nil", func(t *testing.T) {
				var set Set[int]
				require.NoError(t, json.Unmarshal([]byte(`null`), &set))
				require.Equal(t, Set[int](nil), set)
			})
			t.Run("empty", func(t *testing.T) {
				var set Set[int]
				require.NoError(t, json.Unmarshal([]byte(`[]`), &set))
				require.Equal(t, Set[int]{}, set)
			})
			t.Run("non-empty", func(t *testing.T) {
				var set Set[int]
				require.NoError(t, json.Unmarshal([]byte(`[1]`), &set))
				require.Equal(t, Of(1), set)
			})
			t.Run("duplicates", func(t *testing.T) {
				var set Set[int]
				require.NoError(t, json.Unmarshal([]byte(`[1, 1]`), &set))
				require.Equal(t, Of(1), set)
			})
			t.Run("order", func(t *testing.T) {
				var set Set[int]
				require.NoError(t, json.Unmarshal([]byte(`[4, 3, 2, 1]`), &set))
				require.Equal(t, Of(1, 2, 3, 4), set)
			})
			t.Run("embed", func(t *testing.T) {
				var p struct{ Set Set[int] }
				require.NoError(t, json.Unmarshal([]byte(`{ "Set": [1, 2, 3] }`), &p))
				require.Equal(t, Of(1, 2, 3), p.Set)
			})
			t.Run("embed not exist", func(t *testing.T) {
				var p struct{ Set Set[int] }
				require.NoError(t, json.Unmarshal([]byte(`{}`), &p))
				require.Equal(t, Set[int](nil), p.Set)
			})
		})

	})
	t.Run("sortable", func(t *testing.T) {
		t.Run("marshal", func(t *testing.T) {
			requireJSON := func(t *testing.T, value interface{}, expJSON string) {
				actJSON, err := json.Marshal(value)
				require.NoError(t, err)
				require.JSONEq(t, expJSON, string(actJSON))
			}
			t.Run("nil", func(t *testing.T) {
				requireJSON(t, SortableSet[int](nil), `null`)
			})
			t.Run("empty", func(t *testing.T) {
				requireJSON(t, SortableSet[int]{}, `[]`)
			})
			t.Run("non-empty", func(t *testing.T) {
				requireJSON(t, SortableSet[int](Of(1)), `[1]`)
			})
			t.Run("duplicates", func(t *testing.T) {
				requireJSON(t, SortableSet[int](Of(1, 1)), `[1]`)
			})
			t.Run("order", func(t *testing.T) {
				requireJSON(t, SortableSet[int](Of(4, 3, 2, 1)), `[1, 2, 3, 4]`)
			})
		})
		t.Run("unmarshal", func(t *testing.T) {
			t.Run("nil", func(t *testing.T) {
				var set SortableSet[int]
				require.NoError(t, json.Unmarshal([]byte(`null`), &set))
				require.Equal(t, SortableSet[int](nil), set)
			})
			t.Run("empty", func(t *testing.T) {
				var set SortableSet[int]
				require.NoError(t, json.Unmarshal([]byte(`[]`), &set))
				require.Equal(t, SortableSet[int]{}, set)
			})
			t.Run("non-empty", func(t *testing.T) {
				var set SortableSet[int]
				require.NoError(t, json.Unmarshal([]byte(`[1]`), &set))
				require.Equal(t, SortableSet[int](Of(1)), set)
			})
			t.Run("duplicates", func(t *testing.T) {
				var set SortableSet[int]
				require.NoError(t, json.Unmarshal([]byte(`[1, 1]`), &set))
				require.Equal(t, SortableSet[int](Of(1)), set)
			})
			t.Run("order", func(t *testing.T) {
				var set SortableSet[int]
				require.NoError(t, json.Unmarshal([]byte(`[4, 3, 2, 1]`), &set))
				require.Equal(t, SortableSet[int](Of(1, 2, 3, 4)), set)
			})
			t.Run("embed", func(t *testing.T) {
				var p struct{ Set SortableSet[int] }
				require.NoError(t, json.Unmarshal([]byte(`{ "Set": [1, 2, 3] }`), &p))
				require.Equal(t, SortableSet[int](Of(1, 2, 3)), p.Set)
			})
			t.Run("embed not exist", func(t *testing.T) {
				var p struct{ Set SortableSet[int] }
				require.NoError(t, json.Unmarshal([]byte(`{}`), &p))
				require.Equal(t, SortableSet[int](nil), p.Set)
			})
		})

	})
}
