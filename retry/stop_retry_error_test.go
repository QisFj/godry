package retry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_stopRetryError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.NoError(t, StopRetryWithError(nil))
	})
	t.Run("error", func(t *testing.T) {
		err := StopRetryWithError(errors.New("inner error"))
		var sre stopRetryError
		require.True(t, errors.As(err, &sre))
		require.NotNil(t, sre.originError)
		require.Equal(t, "stop retry with error: inner error", sre.Error())
	})
}
