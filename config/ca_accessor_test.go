package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestNewAccessor(t *testing.T) {
	type config struct {
		Value int `validate:"min=5"`
	}

	var getter Getter

	v := validator.New()
	vv := func(raw string, value interface{}) error { return v.Struct(value) }

	a := NewAccessor(func() (string, error) { return getter() }, NewCommonJSONExplainer(config{}), vv)

	// getter return a config raw string
	getter = func() (string, error) {
		return `{"Value": 10}`, nil
	}
	require.Nil(t, a.Config().Value) // before reload, config is nil
	require.NoError(t, a.Reload())
	require.Equal(t, 10, a.Config().Value.(*config).Value) // use type assert to get real type

	// when config change
	getter = func() (string, error) {
		return `{"Value": 15}`, nil
	}
	require.Equal(t, 10, a.Config().Value.(*config).Value) // before reload, config not change
	require.NoError(t, a.Reload())
	require.Equal(t, 15, a.Config().Value.(*config).Value)

	// when getter return an error
	getter = func() (string, error) {
		return `{"Value": 20}`, fmt.Errorf("some error")
	}
	require.EqualError(t, a.Reload(), "some error")        // reload report an error
	require.Equal(t, 15, a.Config().Value.(*config).Value) // config not changed

	// when returned config is invalid
	getter = func() (string, error) {
		return `{"Value": 0}`, nil
	}
	require.Error(t, a.Reload())                           // reload report an error
	require.Equal(t, 15, a.Config().Value.(*config).Value) // config not changed

	// register OnError handler
	var onErrorCalled bool
	a.OnError(func(err error) {
		time.Sleep(50 * time.Millisecond) // won't block Reload
		onErrorCalled = true
	})
	_ = a.Reload()
	require.False(t, onErrorCalled)   // OnError in sleep
	time.Sleep(50 * time.Millisecond) // wait OnError wake up
	require.True(t, onErrorCalled)    // OnError waked up, and set onErrorCalled

	// register OnChange handler
	var onChangeCalled bool
	a.OnChange(func(oldConfig, newConfig Config) {
		onChangeCalled = true
	})
	getter = func() (string, error) {
		return `{"Value": 15}`, nil
	}
	require.NoError(t, a.Reload())
	time.Sleep(50 * time.Millisecond) // make sure OnChange can be called
	require.False(t, onChangeCalled)  // cause config has no change, so OnChange not be called
	getter = func() (string, error) {
		return `{"Value": 20}`, nil
	}
	require.NoError(t, a.Reload())
	time.Sleep(50 * time.Millisecond) // make sure OnChange can be called
	require.True(t, onChangeCalled)   // cause config has been changed, so OnChange has been called
}
