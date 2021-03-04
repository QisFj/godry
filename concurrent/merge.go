package concurrent

import (
	"sync"

	uslice "github.com/QisFj/godry/slice"
)

type EachFunc func() (interface{}, error)
type MergeFunc func(i int, v interface{}) error

func Merge(eachFuncs []EachFunc, merge MergeFunc) error {
	mu := sync.Mutex{}
	return Foreach(eachFuncs, func(i int, v interface{}) error {
		f := v.(EachFunc) // nolint: errcheck
		if f == nil {
			return nil
		}
		val, err := f()
		if err != nil {
			return err
		}
		if merge == nil {
			return nil
		}
		mu.Lock()
		defer mu.Unlock()
		return merge(i, val)
	})
}

func MergeForeach(slice interface{}, f func(i int, v interface{}) (interface{}, error), merge MergeFunc) error {
	var fs []EachFunc
	uslice.Foreach(slice, func(i int, v interface{}) {
		fs = append(fs, func() (interface{}, error) {
			return f(i, v)
		})
	})
	return Merge(fs, merge)
}
