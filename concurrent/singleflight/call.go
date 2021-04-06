package singleflight

import "sync"

// singleflight Call, only call fn once
// it's a wrapper of sync.Once to store return values
type Call struct {
	once sync.Once
	fn   func() (interface{}, error)

	result Result
}

func NewCall(fn func() (interface{}, error)) *Call {
	return &Call{fn: fn}
}

func (c *Call) Do() (interface{}, error) {
	c.once.Do(func() {
		c.result.Val, c.result.Err = c.fn()
	})
	return c.result.Val, c.result.Err
}

func (c *Call) DoChan() <-chan Result {
	ch := make(chan Result, 1)
	go func() {
		val, err := c.Do()
		ch <- Result{
			Val: val,
			Err: err,
		}
	}()
	return ch
}
