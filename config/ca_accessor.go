package config

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type Getter func() (string, error)

type Validator func(raw string, v interface{}) error

func validate(validator Validator, raw string, v interface{}) error {
	if validator == nil {
		return nil
	}
	return validator(raw, v)
}

type Config struct {
	Raw       string      `json:"raw,omitempty"`
	Value     interface{} `json:"value"` // value is Ptr to real config, but must treat it as read only
	UpdatedAt time.Time   `json:"updated_at"`
}

// zero value is not ready to use, call NewAccessor to create
type Accessor struct {
	getter    Getter
	explainer Explainer
	validator Validator

	rw sync.RWMutex
	sf singleflight.Group

	config Config

	// callback
	onError  func(err error)
	onChange func(oldConfig, newConfig Config)
}

// for accessor itself, only getter is required
// for most of use case, explainer is expected
func NewAccessor(getter Getter, explainer Explainer, validator Validator) *Accessor {
	return &Accessor{
		getter:    getter,
		explainer: explainer,
		validator: validator,
	}
}

func (a *Accessor) OnError(onError func(err error)) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.onError = onError
}

func (a *Accessor) callOnError(err error) {
	a.rw.RLock()
	defer a.rw.RUnlock()
	if a.onError == nil {
		return
	}
	go a.onError(err)
}

func (a *Accessor) OnChange(onChange func(oldConfig, newConfig Config)) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.onChange = onChange
}

func (a *Accessor) callOnChange(oldConfig, newConfig Config) {
	a.rw.RLock()
	defer a.rw.RUnlock()
	if a.onChange == nil {
		return
	}
	go a.onChange(oldConfig, newConfig)
}

//nolint:gocritic
func (a *Accessor) Config() Config {
	a.rw.RLock()
	defer a.rw.RUnlock()

	return a.config
}

func (a *Accessor) Reload() error {
	iRaw, err, _ := a.sf.Do("", func() (interface{}, error) {
		return a.getter()
	})
	raw := iRaw.(string) // nolint: errcheck
	if err != nil {
		a.callOnError(err)
		return err
	}
	oldConfig := a.Config()
	if raw == oldConfig.Raw {
		return nil
	}
	// config changed
	v, err := explain(a.explainer, raw)
	if err != nil {
		a.callOnError(err)
		return err
	}
	if err = validate(a.validator, raw, v); err != nil {
		a.callOnError(err)
		return err
	}
	a.rw.Lock()
	a.config = Config{
		Raw:       raw,
		Value:     v,
		UpdatedAt: time.Now(),
	}
	a.rw.Unlock()
	a.callOnChange(oldConfig, a.Config())
	return nil
}
