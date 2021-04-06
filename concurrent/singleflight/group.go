package singleflight

import (
	"sync"
	"time"
)

// a group of Call, use key to identify
type Group struct {
	mu sync.Mutex
	m  map[string]*Call
}

func (g *Group) getCall(key string, fn func() (interface{}, error), forgetAfter time.Duration) *Call {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*Call)
	}
	c, ok := g.m[key]
	if !ok {
		c = NewCall(fn)
		g.m[key] = c
		if forgetAfter > 0 {
			time.AfterFunc(forgetAfter, func() { g.Forget(key) })
		}
	}
	g.mu.Unlock()
	return c
}

// about forgetAfter:
// - timer start before fn be called
// - only make sense if this call not hit cache, so you can't reset forget time
// - non-positive value means never forget
func (g *Group) Do(key string, fn func() (interface{}, error), forgetAfter time.Duration) (interface{}, error) {
	return g.getCall(key, fn, forgetAfter).Do()
}

// unblock Do, return a result channel
func (g *Group) DoChan(key string, fn func() (interface{}, error), forgetAfter time.Duration) <-chan Result {
	return g.getCall(key, fn, forgetAfter).DoChan()
}

func (g *Group) Forget(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.m, key)
}
