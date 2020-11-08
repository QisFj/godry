// Code generated by go generate; DO NOT EDIT.
package slice

type StringMapFunc func(i int, v interface{}) string
type IntMapFunc func(i int, v interface{}) int
type UintMapFunc func(i int, v interface{}) uint

func MapString(slice interface{}, f StringMapFunc) []string {
	sv := valueOf(slice)
	if sv.IsNil() {
		return nil
	}
	list := make([]string, 0, sv.Len())
	foreach(sv, func(i int, v interface{}) {
		list = append(list, f(i, v))
	})
	return list
}
func MapInt(slice interface{}, f IntMapFunc) []int {
	sv := valueOf(slice)
	if sv.IsNil() {
		return nil
	}
	list := make([]int, 0, sv.Len())
	foreach(sv, func(i int, v interface{}) {
		list = append(list, f(i, v))
	})
	return list
}
func MapUint(slice interface{}, f UintMapFunc) []uint {
	sv := valueOf(slice)
	if sv.IsNil() {
		return nil
	}
	list := make([]uint, 0, sv.Len())
	foreach(sv, func(i int, v interface{}) {
		list = append(list, f(i, v))
	})
	return list
}
