package multierr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatterList(t *testing.T) {
	newErrors := func(length int) []error {
		errs := make([]error, 0, length)
		for i := 0; i < length; i++ {
			errs = append(errs, fmt.Errorf("error-%d", i+1))
		}
		return errs
	}
	t.Run("len:0", func(t *testing.T) {
		require.Equal(t, "no error occurred", FormatterList(newErrors(0)))
	})
	t.Run("len:1", func(t *testing.T) {
		require.Equal(t, "1 error occurred: error-1", FormatterList(newErrors(1)))
	})
	t.Run("len:2", func(t *testing.T) {
		require.Equal(t, `2 errors occurred:
	* error-1
	* error-2
`, FormatterList(newErrors(2)))
	})
}
