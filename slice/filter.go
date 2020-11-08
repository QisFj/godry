package slice

import "reflect"

func Filter(slice interface{}, f func(index int) bool) {
	if slice == nil { // do nothing
		return
	}
	sv := reflect.ValueOf(slice)
	if sv.IsNil() { // do nothing
		return
	}
	if sv.Kind() != reflect.Ptr {
		panic("slice must be a pointer to a slice")
	}
	sv = sv.Elem()
	// reflect.Swapper check if sv.Interface() return a slice
	// if not, panic
	swapper := reflect.Swapper(sv.Interface())
	var newLength int
	for i := 0; i < sv.Len(); i++ {
		if f(i) {
			if i != newLength {
				swapper(i, newLength)
			}
			newLength++
		}
	}
	sv.SetLen(newLength)
}
