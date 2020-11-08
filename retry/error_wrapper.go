package retry

import (
	"errors"
	"fmt"
)

type ErrorWrapper func(r *Retry, err error) error

func ErrorWrapperChain(wrappers ...ErrorWrapper) ErrorWrapper {
	return func(r *Retry, err error) error {
		for _, wrapper := range wrappers {
			if wrapper == nil {
				continue
			}
			err = wrapper(r, err)
		}
		return err
	}
}

func FmtErrorWrapper(format string) ErrorWrapper {
	return SkipNoErrorWrap(func(r *Retry, err error) error {
		return fmt.Errorf(format, err)
	})
}

func SkipNoErrorWrap(wrapper ErrorWrapper) ErrorWrapper {
	return func(r *Retry, err error) error {
		if err == nil {
			return nil
		}
		return wrapper(r, err)
	}
}

func errorWrapperRecoverWrap(wrapper ErrorWrapper) ErrorWrapper {
	if wrapper == nil {
		return nil
	}
	return func(r *Retry, er error) (err error) {
		defer func() {
			if p := recover(); p != nil {
				if er == nil {
					er = errors.New("nil error")
				}
				err = fmt.Errorf("error wrapper panic; err: %w, panic: %s", er, p)
			}
		}()
		err = wrapper(r, er)
		return
	}
}
