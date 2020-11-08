package concurrent

import (
	uslice "github.com/QisFj/godry/slice"
)

func Foreach(slice interface{}, f func(i int, v interface{}) error) error {
	var fs []FuncMayError
	uslice.Foreach(slice, func(i int, v interface{}) {
		fs = append(fs, func() error {
			return f(i, v)
		})
	})
	return Do(fs...)
}
