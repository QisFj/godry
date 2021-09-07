package registry

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type testInterface interface {
	method() string
}

type impl1 struct{}
type impl2 struct{}
type impl3 struct{}

func (impl1) method() string { return "impl1" }
func (impl2) method() string { return "impl2" }
func (impl3) method() string { return "impl3" }

type notImpl struct{}

func Test(t *testing.T) {
	t.Run("interface", func(t *testing.T) {
		r := New(reflect.TypeOf((*testInterface)(nil)).Elem(), func(name string) interface{} {
			return impl3{}
		})
		r.Register("1", impl1{})
		r.Register("2", impl2{})
		require.Panics(t, func() {
			// expect panic
			r.Register("wow", notImpl{})
		})

		require.Equal(t, "impl1", r.Get("1").(testInterface).method())
		require.Equal(t, "impl2", r.Get("2").(testInterface).method())
		require.Equal(t, "impl3", r.Get("other").(testInterface).method())
	})

	t.Run("struct", func(t *testing.T) {
		r := New(reflect.TypeOf(impl1{}), nil)
		r.Register("impl1", impl1{})
		require.Panics(t, func() {
			// expect panic
			r.Register("impl1p", &impl1{})
		})
		r.Register("impl2", impl2{})
		r.Register("notImpl", notImpl{})

		_ = r.Get("impl1").(impl1) // should not panic
		_ = r.Get("impl2").(impl1) // should not panic
		require.Panics(t, func() {
			// expect panic
			_ = r.Get("impl2").(impl2)
		})
		require.Nil(t, r.Get("impl3"))
	})
	t.Run("func", func(t *testing.T) {
		type fn func(v int) int
		type fn1 fn
		type fn2 fn

		r := New(reflect.TypeOf(fn1(nil)), nil)

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
		r.Register("add(2)", fn1(add(2)))
		r.Register("mul(3)", fn2(mul(3)))
		r.Register("nilf", fn1(nil))
		require.Panics(t, func() {
			r.Register("nil", nil) // can't register nil
		})

		require.Equal(t, 2, r.Get("add(1)").(fn1)(1))
		require.Equal(t, 4, r.Get("add(2)").(fn1)(2))
		require.Equal(t, 9, r.Get("mul(3)").(fn1)(3))
		require.Panics(t, func() {
			// expect panic
			_ = r.Get("add(1)").(fn)
		})

		require.True(t, r.Get("nilf") != nil) // return fn1(nil), but func can't be compared
		require.True(t, r.Get("nilf").(fn1) == nil)
	})
}
