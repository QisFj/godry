package retry

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Retry struct {
	option    Option
	startOnce sync.Once
	stopOnce  sync.Once
	runCount  uint64
	results   results
	run       chan struct{}
	stop      chan struct{}
	done      chan struct{}
}

type Func func(r *Retry) error

func recoverWrap(f Func) Func {
	return func(r *Retry) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("panic: %s", p)
			}
		}()
		return f(r)
	}
}

func New(option Option) *Retry {
	if option.MaxRunTime == 0 {
		option.MaxRunTime = math.MaxUint64
	}
	retry := Retry{
		option: option,
		run:    make(chan struct{}, 1),
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	retry.results.changeSize(option.ResultSize)
	return &retry
}

func (r *Retry) Start() *Retry {
	go r.startOnce.Do(func() {
		defer func() {
			close(r.done)
		}()
		f := recoverWrap(r.option.F)
		errorWrapper := errorWrapperRecoverWrap(r.option.ErrorWrapper)
		r.run <- struct{}{}
		for ; r.RunCount() < r.option.MaxRunTime; atomic.AddUint64(&r.runCount, 1) {
			select {
			case <-r.stop:
				return
			case <-r.run:
				err := f(r)
				var sre stopRetryError
				if errors.As(err, &sre) {
					err = sre.originError
				}
				if r.option.ErrorWrapper != nil {
					err = errorWrapper(r, err)
				}
				r.results.append(Result{
					Valid: true,
					Error: err,
				})
				if err == nil || sre.originError != nil {
					return
				}
				if r.option.RetryInterval == 0 {
					r.run <- struct{}{}
				} else {
					time.AfterFunc(r.option.RetryInterval, func() {
						r.run <- struct{}{}
					})
				}
			}
		}
	})
	return r
}

func (r *Retry) RunCount() uint64 { // return how many times have been executed
	return atomic.LoadUint64(&r.runCount)
}

func (r *Retry) Stop() {
	r.stopOnce.Do(func() {
		close(r.stop)
	})
}

func (r *Retry) HasDone() bool {
	select {
	case <-r.done:
		return true
	case <-r.stop:
		return true
	default:
		return false
	}
}

func (r *Retry) Wait() *Retry {
	<-r.done
	return r
}

func (r *Retry) Result(runTime int) Result {
	return r.results.get(runTime - 1)
}

func (r *Retry) LatestResult() Result {
	return r.results.getLatest()
}
