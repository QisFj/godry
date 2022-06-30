package concurrent

import "sync"

func Go(f FuncWithResultMayError, options ...GoOption) {
	var goOptions GoOptions
	for _, option := range options {
		if option != nil {
			option(&goOptions)
		}
	}
	for _, w := range goOptions.wrappers {
		f = w(f)
	}
	if goOptions.skipNilFunc && f == nil {
		return
	}
	if goOptions.wg != nil {
		goOptions.wg.Add(1)
	}
	if goOptions.l != nil {
		goOptions.l.Acquire()
	}
	go func() {
		if goOptions.wg != nil {
			defer goOptions.wg.Done()
		}
		if goOptions.l != nil {
			defer goOptions.l.Release()
		}
		if goOptions.panicHandler != nil {
			defer func() {
				if p := recover(); p != nil {
					goOptions.panicHandler(p)
				}

			}()
		}
		result, err := f()
		if err != nil {
			if goOptions.errHandler != nil {
				goOptions.errHandler(err)
			}
			return
		}
		if goOptions.resultHandler != nil {
			goOptions.resultHandler(result)
		}
	}()
}

type GoOption func(*GoOptions)

type GoOptions struct {
	wrappers []func(f FuncWithResultMayError) FuncWithResultMayError

	skipNilFunc bool

	wg *sync.WaitGroup
	l  *Limiter

	panicHandler  func(p interface{})
	errHandler    func(err error)
	resultHandler func(result interface{})
}

// set checkFuncNil to return nil when f is nil
func WithWrapper(w func(f FuncWithResultMayError) FuncWithResultMayError, checkFuncNil bool) GoOption {
	if w == nil {
		return nil
	}
	if checkFuncNil {
		originW := w
		w = func(f FuncWithResultMayError) FuncWithResultMayError {
			if f == nil {
				return nil
			}
			return originW(f)
		}
	}
	return func(options *GoOptions) {
		options.wrappers = append(options.wrappers, w)
	}
}

func WithSkipNilFunc(skip bool) GoOption {
	return func(options *GoOptions) {
		options.skipNilFunc = skip
	}
}

func WithWaitGroup(wg *sync.WaitGroup) GoOption {
	return func(options *GoOptions) {
		options.wg = wg
	}
}

func WithLimiter(l *Limiter) GoOption {
	return func(options *GoOptions) {
		options.l = l
	}
}

// nil to not recover
func WithRecoverPanic(f func(p interface{})) GoOption {
	return func(options *GoOptions) {
		options.panicHandler = f
	}
}

func WithErrorHandler(f func(err error)) GoOption {
	return func(options *GoOptions) {
		options.errHandler = f
	}
}

func WithResultHandler(f func(result interface{})) GoOption {
	return func(options *GoOptions) {
		options.resultHandler = f
	}
}
