package slice

import "reflect"

func Foreach(slice interface{}, f func(i int, v interface{})) {
	foreach(valueOf(slice), f)
}

func foreach(sv reflect.Value, f func(i int, v interface{})) {
	for i := 0; i < sv.Len(); i++ {
		f(i, sv.Index(i).Interface())
	}
}

// reflect.ValueOf for slice, compared to reflect.ValueOf, it requires Kind to be reflect.Slice
func valueOf(slice interface{}) reflect.Value {
	sv := reflect.ValueOf(slice)
	if sv.Kind() != reflect.Slice {
		panic("slice must be a slice")
	}
	return sv
}
