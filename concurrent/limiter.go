package concurrent

import "sync"

// wrap sync.Cond on a uint, use to limit concurrent
//
// zero value is not read to use, use NewLimiter to create
// a nil *Limiter is a noop Limiter
type Limiter struct {
	rest int
	cond sync.Cond
}

func NewLimiter(n int) *Limiter {
	return &Limiter{
		rest: n,
		cond: sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

func (l *Limiter) Release() {
	if l == nil {
		return
	}
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	l.rest++
	l.cond.Signal()
}

func (l *Limiter) Acquire() {
	if l == nil {
		return
	}
	l.cond.L.Lock()
	defer l.cond.L.Unlock()
	if l.rest > 0 {
		l.rest--
		return
	}
	l.cond.Wait()
	l.rest--
}
