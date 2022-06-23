package graph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/QisFj/godry/slice"
)

func TestNewIter(t *testing.T) {
	repeatLayer := func(layer LayerI, n int) []LayerI {
		layers := make([]LayerI, n)
		for i := 0; i < n; i++ {
			layers[i] = layer
		}
		return layers
	}
	splitStringToLayer := func(s string) LayerI {
		return Layer(slice.Map([]rune(s), func(index int, value rune) NodeI {
			return value
		}))
	}
	t.Run("decimal", func(t *testing.T) {
		n09 := NumberRange{From: 0, Length: 10}
		g := Graph(repeatLayer(n09, 4)) // 4 bits
		it := NewIter(g, false)
		t.Run("iterate", func(t *testing.T) {
			exp := 0
			for it.Next() {
				nodes := it.Get()
				act := 0
				for _, node := range nodes {
					act = 10*act + node.(int)
				}
				require.Equal(t, exp, act)
				exp++
			}
			require.Equal(t, 10000, exp) // after iterate all decimal 4 bits number
		})
		t.Run("reuse", func(t *testing.T) {
			exp := 0
			for it.Next() {
				nodes := it.Get()
				act := 0
				for _, node := range nodes {
					act = 10*act + node.(int)
				}
				require.Equal(t, exp, act)
				exp++
			}
			require.Equal(t, 10000, exp) // after iterate all decimal 4 bits number
		})
	})
	t.Run("binary", func(t *testing.T) {
		// for all uint8
		n01 := NumberRange{From: 0, Length: 2}
		g := Graph(repeatLayer(n01, 8)) // 8 bits
		it := NewIter(g, false)
		results := [256]bool{}
		for it.Next() {
			nodes := it.Get()
			value := uint8(0)
			for _, node := range nodes {
				value = (value << 1) | uint8(node.(int))
			}
			results[value] = true
		}
		for index, result := range results {
			require.True(t, result, "%d", index)
		}
	})
	t.Run("reverse", func(t *testing.T) {
		// ABC abc 123
		g := Graph{
			splitStringToLayer("ABC"),
			splitStringToLayer("abc"),
			splitStringToLayer("123"),
		}
		t.Run("without reverse", func(t *testing.T) {
			it := NewIter(g, false) // iterate start from "123"
			results := make([]string, 0, 27)
			for it.Next() {
				nodes := it.Get()
				result := new(strings.Builder)
				for _, node := range nodes {
					result.WriteRune(node.(rune))
				}
				results = append(results, result.String())
			}
			require.Equal(t, []string{
				"Aa1", "Aa2", "Aa3",
				"Ab1", "Ab2", "Ab3",
				"Ac1", "Ac2", "Ac3",

				"Ba1", "Ba2", "Ba3",
				"Bb1", "Bb2", "Bb3",
				"Bc1", "Bc2", "Bc3",

				"Ca1", "Ca2", "Ca3",
				"Cb1", "Cb2", "Cb3",
				"Cc1", "Cc2", "Cc3",
			}, results)
		})
		t.Run("within reverse", func(t *testing.T) {
			it := NewIter(g, true) // iterate start from "ABC"
			results := make([]string, 0, 27)
			for it.Next() {
				nodes := it.Get()
				result := new(strings.Builder)
				for _, node := range nodes {
					result.WriteRune(node.(rune))
				}
				results = append(results, result.String())
			}
			require.Equal(t, []string{
				"Aa1", "Ba1", "Ca1",
				"Ab1", "Bb1", "Cb1",
				"Ac1", "Bc1", "Cc1",

				"Aa2", "Ba2", "Ca2",
				"Ab2", "Bb2", "Cb2",
				"Ac2", "Bc2", "Cc2",

				"Aa3", "Ba3", "Ca3",
				"Ab3", "Bb3", "Cb3",
				"Ac3", "Bc3", "Cc3",
			}, results)
		})
	})
	//it := NewIter(Graph{
	//	Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	//	Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	//	Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	//	Layout{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	//}, false)
	//t.Run("first loop", func(t *testing.T) {
	//	e := 0
	//	for it.Next() {
	//		nodes := it.Get()
	//		a := 0
	//		for _, node := range nodes {
	//			a = 10*a + node.(int)
	//		}
	//		require.Equal(t, e, a)
	//		e++
	//	}
	//	require.Equal(t, 10000, e)
	//})
	//t.Run("second loop", func(t *testing.T) {
	//	e := 0
	//	for it.Next() {
	//		nodes := it.Get()
	//		a := 0
	//		for _, node := range nodes {
	//			a = 10*a + node.(int)
	//		}
	//		require.Equal(t, e, a)
	//		e++
	//	}
	//	require.Equal(t, 10000, e)
	//})
	//t.Run("reverse", func(t *testing.T) {
	//	it.reverse = true
	//	e := 0
	//	for it.Next() {
	//		nodes := it.Get()
	//		a := 0
	//		for i := len(nodes) - 1; i >= 0; i-- {
	//			a = 10*a + nodes[i].(int)
	//		}
	//		require.Equal(t, e, a)
	//		e++
	//	}
	//	require.Equal(t, 10000, e)
	//})
}
