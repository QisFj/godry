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

// about forgetAfter:
// - timer start before fn be called
// - only make sense if this call not hit cache, so you can't reset forget time
// - non-positive value means never forget
func (g *Group) Do(key string, fn func() (interface{}, error), forgetAfter time.Duration) (interface{}, error) {
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
	return c.Do()
}

func (g *Group) Forget(key string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.m, key)
}
