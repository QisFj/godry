package channels

import (
	"reflect"
)

type ReadableCh[T any] interface {
	~chan T | ~<-chan T
}

func Read[T any, Ch ReadableCh[T]](chs []Ch) (chosenIndex int, value T, ok bool) {
	cases := make([]reflect.SelectCase, len(chs))
	for i, ch := range chs {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	var rValue reflect.Value
	chosenIndex, rValue, ok = reflect.Select(cases)
	if ok {
		// not closed, should read value
		value = rValue.Interface().(T)
	}
	return
}
