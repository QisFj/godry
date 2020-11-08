package multierr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		errs := New(nil, nil)
		errs.AppendOnlyNotNil(nil)
		errs.AppendOnlyNotNil(nil)
		require.NoError(t, errs.Error())
	})
	t.Run("with errors", func(t *testing.T) {
		errs := New(func(err error) error {
			return fmt.Errorf("-> %w", err)
		}, nil)
		for i := 0; i < 6; i++ {
			var err error
			if i%2 == 0 {
				err = fmt.Errorf("error-%d", i+1)
			}
			errs.AppendOnlyNotNil(err)
		}
		require.EqualError(t, errs.Error(), `3 errors occurred:
	* -> error-1
	* -> error-3
	* -> error-5
`)
	})
}
