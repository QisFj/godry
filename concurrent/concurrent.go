package concurrent

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/QisFj/godry/multierr"
)

type FuncMayError func() error

func Do(functions ...FuncMayError) error {
	fs := make([]FuncMayError, 0, len(functions))
	for _, f := range functions {
		if f == nil {
			continue
		}
		fs = append(fs, panicRecoverWrap(f))
	}
	if len(fs) == 0 {
		return nil
	}
	errs := multierr.New(nil, nil)
	if len(fs) == 1 {
		errs.AppendOnlyNotNil(fs[0]())
		return errs.Error()
	} else {
		wg := sync.WaitGroup{}
		for _, function := range fs {
			wg.Add(1)
			go func(f FuncMayError) {
				defer wg.Done()
				errs.AppendOnlyNotNil(f())
			}(function)
		}
		wg.Wait()
	}
	return errs.Error()
}

func panicRecoverWrap(f FuncMayError) FuncMayError {
	return func() (err error) {
		defer func() {
			if p := recover(); p != nil {
				switch pp := p.(type) {
				case error:
					err = fmt.Errorf("panic: %w, stack:\n%s", pp, string(debug.Stack()))
				case string:
					err = fmt.Errorf("panic: %s, stack:\n%s", pp, string(debug.Stack()))
				default:
					err = fmt.Errorf("panic: %#v, stack:\n%s", pp, string(debug.Stack()))
				}
			}
		}()
		err = f()
		return
	}
}
