package set

import (
	"reflect"
)

//go:generate go run set.gen.go

type T struct {
	t reflect.Type // set key type
	m reflect.Value
}

func NewT(t reflect.Type) T {
	set := T{}
	set.New(t)
	return set
}

var setValue = struct{}{}
var setValueRV = reflect.ValueOf(setValue)
var setValueRT = reflect.TypeOf(setValue)

func (set *T) New(t reflect.Type) {
	set.t = t
	set.m = reflect.MakeMap(reflect.MapOf(t, setValueRT))
}

func (set *T) add(e reflect.Value) {
	set.m.SetMapIndex(e, setValueRV)
}

func (set *T) Add(es ...interface{}) {
	for _, e := range es {
		set.add(reflect.ValueOf(e))
	}
}

func (set T) List() interface{} {
	l := reflect.MakeSlice(reflect.SliceOf(set.t), 0, set.m.Len())
	l = reflect.Append(l, set.m.MapKeys()...)
	return l.Interface()
}

func (set T) contains(e reflect.Value) bool {
	return set.m.MapIndex(e) != reflect.Value{}
}

func (set T) Contains(e interface{}) bool {
	return set.contains(reflect.ValueOf(e))
}

func (set T) ContainsAny(es ...interface{}) bool {
	for _, e := range es {
		if set.Contains(e) {
			return true
		}
	}
	return false
}

func (set T) ContainsAll(es ...interface{}) bool {
	for _, e := range es {
		if !set.Contains(e) {
			return false
		}
	}
	return true
}

func TDiff(set1, set2 T) (both, only1, only2 T) {
	if set1.t != set2.t {
		panic("two sets must be same type")
	}
	both, only1, only2 = NewT(set1.t), NewT(set1.t), NewT(set1.t)
	for _, v := range set1.m.MapKeys() {
		if set2.contains(v) {
			both.add(v)
		} else {
			only1.add(v)
		}
	}
	for _, v := range set2.m.MapKeys() {
		if set1.contains(v) {
			both.add(v)
		} else {
			only2.add(v)
		}
	}
	return
}

func TMerge(set1, set2 T) T {
	if set1.t != set2.t {
		panic("two sets must be same type")
	}
	set := NewT(set1.t)
	for _, v := range set1.m.MapKeys() {
		set.m.SetMapIndex(v, setValueRV)
	}
	for _, v := range set2.m.MapKeys() {
		set.m.SetMapIndex(v, setValueRV)
	}
	return set
}
