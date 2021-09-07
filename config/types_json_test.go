package config

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	dur := Duration(time.Hour)

	marshaled, err := json.Marshal(dur)
	require.NoError(t, err)
	t.Logf("marshaled: %s", string(marshaled))

	var result Duration
	err = json.Unmarshal(marshaled, &result)
	require.NoError(t, err)

	require.Equal(t, dur, result)
}

func TestRegexp(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		r := Regexp{Regexp: regexp.MustCompile("\",\\\\,{,\\[")}

		marshaled, err := json.Marshal(r)
		require.NoError(t, err)
		t.Logf("marshaled: %s", string(marshaled))

		var result Regexp
		err = json.Unmarshal(marshaled, &result)
		require.NoError(t, err)

		require.Equal(t, r, result)
	})
	t.Run("empty", func(t *testing.T) {
		r := Regexp{}

		marshaled, err := json.Marshal(r)
		require.NoError(t, err)
		t.Logf("marshaled: %s", string(marshaled))

		var result Regexp
		err = json.Unmarshal(marshaled, &result)
		require.NoError(t, err)

		require.Equal(t, r, result)
	})
}

func TestRegexp_Unmarshal(t *testing.T) {
	t.Run("embed, not exist", func(t *testing.T) {
		var v struct{ Regexp Regexp }
		require.NoError(t, json.Unmarshal([]byte(`{}`), &v))
		require.Nil(t, v.Regexp.Regexp)
	})
	t.Run("null", func(t *testing.T) {
		var v Regexp
		require.NoError(t, json.Unmarshal([]byte(`null`), &v))
		require.Nil(t, v.Regexp)
	})
	t.Run("{}", func(t *testing.T) {
		var v Regexp
		require.Error(t, json.Unmarshal([]byte(`{}`), &v))
	})
	t.Run(`""`, func(t *testing.T) {
		var v Regexp
		require.NoError(t, json.Unmarshal([]byte(`""`), &v))
		require.Equal(t, "", v.Regexp.String()) // not nil
	})
}
