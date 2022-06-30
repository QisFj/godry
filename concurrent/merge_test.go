package concurrent

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.NoError(t, Merge(nil, nil))
	})
	t.Run("with nil", func(t *testing.T) {
		require.NoError(t, Merge([]FuncWithResultMayError{nil}, nil))
	})
	t.Run("normal", func(t *testing.T) {
		slice := []int{0, 0, 0}
		sum := 0
		require.NoError(t, Merge([]FuncWithResultMayError{
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				return 0, nil
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				return 1, nil
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return 2, nil
			},
		}, func(_ int, v interface{}) error {
			sum += v.(int)
			return nil
		}))
		require.Equal(t, []int{0, 1, 2}, slice)
		require.Equal(t, sum, 3)
	})
	t.Run("error", func(t *testing.T) {
		slice := []int{0, 0, 0}
		sum := 0
		require.Error(t, Merge([]FuncWithResultMayError{
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				return 0, errors.New("error")
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				return 1, errors.New("error")
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return 2, nil
			},
		}, func(_ int, v interface{}) error {
			sum += v.(int)
			return nil
		}), `2 errors occurred:
	* error
	* error
`)
		require.Equal(t, []int{0, 1, 2}, slice)
		require.Equal(t, sum, 2) // () + () + (2)
	})
	t.Run("error on merge", func(t *testing.T) {
		slice := []int{0, 0, 0}
		sum := 0
		require.Error(t, Merge([]FuncWithResultMayError{
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				return 0, nil
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				return 1, nil
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return 2, nil
			},
		}, func(_ int, v interface{}) error {
			sum += v.(int)
			return errors.New("error")
		}), `3 errors occurred:
	* error
	* error
	* error
`)
		require.Equal(t, []int{0, 1, 2}, slice)
		require.Equal(t, sum, 3)
	})
	t.Run("panic", func(t *testing.T) {
		slice := []int{0, 0, 0}
		sum := 0
		err := Merge([]FuncWithResultMayError{
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				panic("error")
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				panic("error")
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return 2, nil
			},
		}, func(_ int, v interface{}) error {
			sum += v.(int)
			return nil
		})
		require.Error(t, err)
		t.Logf("got error: %s", err)
		require.Equal(t, []int{0, 1, 2}, slice)
		require.Equal(t, sum, 2)
	})
	t.Run("panic on merge", func(t *testing.T) {
		slice := []int{0, 0, 0}
		sum := 0
		err := Merge([]FuncWithResultMayError{
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				return 0, nil
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				return 1, nil
			},
			func() (interface{}, error) {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return 2, nil
			},
		}, func(_ int, v interface{}) error {
			sum += v.(int)
			panic("error")
		})
		require.Error(t, err)
		t.Logf("got error: %s", err)
		require.Equal(t, []int{0, 1, 2}, slice)
		require.Equal(t, sum, 3)
	})
}

func TestMergeForeach(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	var value int
	require.NoError(t, MergeForeach(ints, func(i int, v interface{}) (interface{}, error) {
		val := v.(int)
		return val * val, nil
	}, func(_ int, v interface{}) error {
		t.Logf("value: %d, v: %v", value, v)
		value += v.(int) // no concurrency here, it's ok add direct
		return nil
	}))
	// value = 1 + 4 + 9 + 16 + 25 = 55
	require.Equal(t, 55, value)
}
