package multierr

import (
	"errors"
	"fmt"
	"sync"
)

type Errs struct {
	rw     sync.RWMutex
	errors []error

	wrapper   ErrWrapper
	formatter Formatter
}

func New(wrapper ErrWrapper, formatter Formatter) *Errs {
	return &Errs{
		wrapper:   wrapper,
		formatter: formatter,
	}
}

func (errs *Errs) Append(err error) {
	errs.rw.Lock()
	defer errs.rw.Unlock()
	if errs.wrapper != nil {
		err = errs.wrapper(err)
	}
	errs.errors = append(errs.errors, err)
}

func (errs *Errs) AppendOnlyNotNil(err error) {
	if err == nil {
		return
	}
	errs.Append(err)
}

func (errs *Errs) Appendf(format string, v ...interface{}) {
	errs.Append(fmt.Errorf(format, v...))
}

func (errs *Errs) Error() error {
	if errs == nil {
		return nil
	}
	errs.rw.RLock()
	defer errs.rw.RUnlock()
	if len(errs.errors) == 0 {
		return nil
	}
	if errs.formatter == nil {
		errs.formatter = FormatterList
	}
	return errors.New(errs.formatter(errs.errors))
}

func (errs *Errs) IsNil() bool {
	if errs == nil {
		return true
	}
	errs.rw.RLock()
	defer errs.rw.RUnlock()
	return len(errs.errors) == 0
}
