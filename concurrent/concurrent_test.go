package concurrent

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.NoError(t, Do())
	})
	t.Run("with nil", func(t *testing.T) {
		require.NoError(t, Do(nil))
	})
	t.Run("one", func(t *testing.T) {
		slice := []int{0}
		require.NoError(t, Do(func() error {
			slice[0] = 1
			return nil
		}))
		require.Equal(t, []int{1}, slice)
	})
	t.Run("one with nil", func(t *testing.T) {
		slice := []int{0}
		require.NoError(t, Do(nil, func() error {
			slice[0] = 1
			return nil
		}, nil))
		require.Equal(t, []int{1}, slice)
	})
	t.Run("one with error", func(t *testing.T) {
		require.EqualError(t, Do(func() error {
			return errors.New("EXPECTED ERROR")
		}), "1 error occurred: EXPECTED ERROR")
	})
	t.Run("normal", func(t *testing.T) {
		slice := []int{0, 0, 0}
		require.NoError(t, Do(
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				return nil
			},
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				return nil
			},
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return nil
			},
		))
		require.Equal(t, []int{0, 1, 2}, slice)
	})
	t.Run("error", func(t *testing.T) {
		slice := []int{0, 0, 0}
		require.EqualError(t, Do(
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				return errors.New("error")
			},
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				return errors.New("error")
			},
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return nil
			},
		), `2 errors occurred:
	* error
	* error
`)
		require.Equal(t, []int{0, 1, 2}, slice)
	})
	t.Run("panic", func(t *testing.T) {
		slice := []int{0, 0, 0}
		err := Do(
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[0] = 0
				panic("error")
			},
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[1] = 1
				panic("error")
			},
			func() error {
				time.Sleep(20 * time.Millisecond)
				slice[2] = 2
				return nil
			},
		)
		require.Error(t, err)
		t.Logf("got error: %s", err)
		require.Equal(t, []int{0, 1, 2}, slice)
	})
}
