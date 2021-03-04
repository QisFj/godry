package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCommonJSONExplainer(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v, err := explain(NewCommonJSONExplainer(0), `1`)
		require.NoError(t, err)
		expected := 1
		require.Equal(t, &expected, v)
	})
	t.Run("string", func(t *testing.T) {
		v, err := explain(NewCommonJSONExplainer("0"), `"1"`)
		require.NoError(t, err)
		expected := "1"
		require.Equal(t, &expected, v)
	})
	t.Run("struct", func(t *testing.T) {
		type Struct struct {
			V int
		}
		v, err := explain(NewCommonJSONExplainer(Struct{}), `{"V":1}`)
		require.NoError(t, err)
		expected := Struct{V: 1}
		require.Equal(t, &expected, v)
	})
	t.Run("ptr & interface{}", func(t *testing.T) {
		type Struct struct {
			V int
		}
		v, err := explain(NewCommonJSONExplainer(interface{}(&Struct{})), `{"V":1}`)
		require.NoError(t, err)
		expected := Struct{V: 1}
		require.Equal(t, &expected, v)
	})
	t.Run("array", func(t *testing.T) {
		v, err := explain(NewCommonJSONExplainer([]int{}), `[1,2,3]`)
		require.NoError(t, err)
		expected := []int{1, 2, 3}
		require.Equal(t, &expected, v)
	})
	t.Run("map", func(t *testing.T) {
		v, err := explain(NewCommonJSONExplainer(map[string]int{}), `{"1":1,"2":2,"3":3}`)
		require.NoError(t, err)
		expected := map[string]int{"1": 1, "2": 2, "3": 3}
		require.Equal(t, &expected, v)
	})
}
