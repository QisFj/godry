package slice

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMapT(t *testing.T) {
	mapF := func(i int, v interface{}) interface{} {
		tm, _ := time.Parse("2006-01-02 15:04:05", v.(string))
		return tm
	}
	timeType := reflect.TypeOf(time.Time{})
	t.Run("normal", func(t *testing.T) {
		timeStrings := []string{
			"1970-01-01 00:00:00",
			"1970-01-01 00:00:01",
			"1970-01-01 00:00:02",
		}
		for _, tm := range MapT(timeStrings, timeType, mapF).([]time.Time) {
			t.Logf("time: %s", tm)
		}
	})
	t.Run("nil, not slice", func(t *testing.T) {
		// would panic
		defer func() {
			if p := recover(); p != nil {
				t.Log("panic as expected")
				return
			}
			t.Error("not panic")
		}()
		_ = MapT(nil, timeType, mapF).([]time.Time)
	})
	t.Run("nil slice", func(t *testing.T) {
		times := MapT([]string(nil), timeType, mapF).([]time.Time)
		require.Nil(t, times)
		require.Len(t, times, 0)
	})
	t.Run("empty, not nil slice", func(t *testing.T) {
		times := MapT([]string{}, timeType, mapF).([]time.Time)
		require.NotNil(t, times)
		require.Len(t, times, 0)
	})

}
