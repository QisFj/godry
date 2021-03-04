package config

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMultiAccessor(t *testing.T) {
	ma := &MultiAccessor{}

	newCounterGetter := func(counter int) Getter {
		return func() (string, error) {
			counter += 2
			return strconv.Itoa(counter), nil
		}
	}

	explainer := func(raw string) (interface{}, error) {
		return strconv.Atoi(raw)
	}

	ma.AddAccessor("1", NewAccessor(newCounterGetter(-1), explainer, nil))
	ma.AddAccessor("2", NewAccessor(newCounterGetter(0), explainer, nil))

	require.NoError(t, ma.Reload()) // reload all Accessor
	require.Equal(t, 1, ma.Config("1").Value.(int))
	require.Equal(t, 2, ma.Config("2").Value.(int))

	var exist bool
	// use GetAccessor to get Accessor, and check if Accessor exist
	_, exist = ma.GetAccessor("1")
	require.True(t, exist)
	_, exist = ma.GetAccessor("2")
	require.True(t, exist)
	_, exist = ma.GetAccessor("3")
	require.False(t, exist)
	// use GetConfig to get Config, and check if Config exist
	_, exist = ma.GetConfig("1")
	require.True(t, exist)
	_, exist = ma.GetConfig("2")
	require.True(t, exist)
	_, exist = ma.GetConfig("3")
	require.False(t, exist)
	// use Accessor to get Accessor, get nil if not exist
	require.NotNil(t, ma.Accessor("1"))
	require.NotNil(t, ma.Accessor("2"))
	require.Nil(t, ma.Accessor("3"))
	// use Config to get Config, get empty(zero value) Config if not exist
	require.NotZero(t, ma.Config("1"))
	require.NotZero(t, ma.Config("2"))
	require.Zero(t, ma.Config("3"))

	ctx, cancel := context.WithCancel(context.Background())
	go ma.ReloadEach(100*time.Millisecond, ctx.Done()) // reload each 100ms
	// wait reload several times
	time.Sleep(550 * time.Millisecond)
	cancel() // stop reload

	// reload at 100ms, 200ms, 300ms, 400ms, 500ms, total 5 times
	// value for config 1 should be 11
	// value for config 2 should be 12
	configs := ma.Configs() // use Configs get all Config
	require.Equal(t, 11, configs["1"].Value.(int))
	require.Equal(t, 12, configs["2"].Value.(int))
}
