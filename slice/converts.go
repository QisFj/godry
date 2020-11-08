package slice

import (
	"reflect"
)

// The functions in this file(converts.go) are related to the type conversion between map and slice

// convert a map to a key slice
func KeysOfMap(v interface{}) interface{} {
	// v must be a map
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		panic("v must be a map")
	}
	list := reflect.MakeSlice(reflect.SliceOf(rv.Type().Key()), 0, rv.Len())
	list = reflect.Append(list, rv.MapKeys()...)
	return list.Interface()
}

// convert a map to a value slice
func ValuesOfMap(v interface{}) interface{} {
	// v must be a map
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		panic("v must be a map")
	}
	list := reflect.MakeSlice(reflect.SliceOf(rv.Type().Elem()), 0, rv.Len())
	it := rv.MapRange()
	for it.Next() {
		list = reflect.Append(list, it.Value())
	}
	return list.Interface()
}

type KV struct {
	Key   interface{}
	Value interface{}
}

// convert a map to a kv slice
// so that other function in this package can be used
func KVsOfMap(v interface{}) []KV {
	// v must be a map, result's Key is map's key, Value is map's value
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		panic("v must be a map")
	}
	kvs := make([]KV, 0, rv.Len())
	it := rv.MapRange()
	for it.Next() {
		kvs = append(kvs, KV{
			Key:   it.Key().Interface(),
			Value: it.Value().Interface(),
		})
	}
	return kvs
}

var typeOfEmptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()

// convert a slice to a map
// kt, getK must not nil
// vt, getV can be nil
// value type would be interface{} when vt == nil
// value would be elem of slice when getV == nil
func ToMap(slice interface{},
	kt reflect.Type, getK func(i int, v interface{}) interface{},
	vt reflect.Type, getV func(i int, v interface{}) interface{},
) interface{} {
	if vt == nil {
		vt = typeOfEmptyInterface
	}
	if getV == nil {
		getV = func(i int, v interface{}) interface{} {
			return v
		}
	}
	m := reflect.MakeMap(reflect.MapOf(kt, vt))
	Foreach(slice, func(i int, v interface{}) {
		m.SetMapIndex(reflect.ValueOf(getK(i, v)), reflect.ValueOf(getV(i, v)))
	})
	return m.Interface()
}
