package slice

import "reflect"

//go:generate go run map.gen.go

type InterfaceMapFunc func(i int, v interface{}) interface{}

var typeAssertMapForInterface InterfaceMapFunc = func(_ int, v interface{}) interface{} { return v }

func MapInterface(slice interface{}, f InterfaceMapFunc) []interface{} {
	sv := valueOf(slice)
	if sv.IsNil() {
		return nil
	}
	if f == nil {
		f = typeAssertMapForInterface
	}
	list := make([]interface{}, 0, sv.Len())
	foreach(sv, func(i int, v interface{}) {
		list = append(list, f(i, v))
	})
	return list
}

func MapT(slice interface{}, t reflect.Type, f InterfaceMapFunc) interface{} {
	sv := valueOf(slice)
	if sv.IsNil() {
		return reflect.New(reflect.SliceOf(t)).Elem().Interface() // return a nil slice
	}
	if f == nil {
		f = typeAssertMapForInterface
	}
	list := reflect.MakeSlice(reflect.SliceOf(t), 0, reflect.ValueOf(slice).Len())
	foreach(sv, func(i int, v interface{}) {
		list = reflect.Append(list, reflect.ValueOf(f(i, v)))
	})
	return list.Interface()
}
