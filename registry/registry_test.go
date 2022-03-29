package registry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	type fn func(v int) int
	r := New("", AlwaysReturn(func(v int) int { return v }))
	add := func(n int) fn {
		return func(v int) int {
			return v + n
		}
	}
	mul := func(n int) fn {
		return func(v int) int {
			return v * n
		}
	}
	r.Register("add(1)", add(1))
	r.Register("add(2)", add(2))
	r.Register("mul(3)", mul(3))

	require.Panics(t, func() {
		// expect panic
		r.Register("add(1)", add(1))
	})

	require.Equal(t, 2, r.Get("add(1)")(1))
	require.Equal(t, 3, r.Get("add(2)")(1))
	require.Equal(t, 3, r.Get("mul(3)")(1))
}
