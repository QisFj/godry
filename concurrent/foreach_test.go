package concurrent

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestForeach(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	require.NoError(t, Foreach(slice, func(i int, v interface{}) error {
		t.Logf("slice[%d]=%v start", i, v)
		time.Sleep(20 * time.Millisecond)
		t.Logf("slice[%d]=%v stop", i, v)
		return nil
	}))
}
