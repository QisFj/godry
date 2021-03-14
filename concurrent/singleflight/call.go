package singleflight

import "sync"

// singleflight Call, only call fn once
// it's a wrapper of sync.Once to store return values
type Call struct {
	once sync.Once
	fn   func() (interface{}, error)

	value interface{}
	err   error
}

func NewCall(fn func() (interface{}, error)) *Call {
	return &Call{fn: fn}
}

func (c *Call) Do() (interface{}, error) {
	c.once.Do(func() {
		c.value, c.err = c.fn()
	})
	return c.value, c.err
}
