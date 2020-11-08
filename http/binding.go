package http

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/QisFj/godry/name"
	"github.com/QisFj/godry/set"
	"github.com/QisFj/godry/slice"
)

type Binding interface {
	Bind(*http.Request, interface{}) error
}

type BindingFunc func(*http.Request, interface{}) error

func (f BindingFunc) Bind(r *http.Request, v interface{}) error {
	return f(r, v)
}

type QueryBinding struct {
	Validate ValidateFunc
}

func (b QueryBinding) Bind(r *http.Request, v interface{}) (err error) {
	if err = b.bind(r, v); err != nil {
		return fmt.Errorf("bind error: %w", err)
	}
	if b.Validate != nil {
		if err = b.Validate(v); err != nil {
			return fmt.Errorf("validate error: %w", err)
		}
	}
	return nil
}
func (b QueryBinding) bind(r *http.Request, v interface{}) error {
	// v must be a point to struct
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("can't bind to type: %s", reflect.TypeOf(v))
	}
	rv = rv.Elem()
	s := newQueryBindingSetter(r.URL.Query())
	for i := 0; i < rv.NumField(); i++ {
		if err := s.set(rv.Field(i), nil, rv.Type().Field(i)); err != nil {
			return err
		}
	}
	return nil
}

type queryBindingSetter struct {
	values   url.Values
	usedKeys set.String
}

func newQueryBindingSetter(values url.Values) *queryBindingSetter {
	// todo[maybe]: allowed use provider key conversion function
	// todo[maybe]: allowed use make decision when key duplicated
	newValues := url.Values{}
	for k, vs := range values {
		k = name.ToSnakeCase(k)
		// original k would never duplicated, but converted k might duplicated
		// when duplicated just Add Value
		// url.Values.Add not support add multiple value, use append instead
		newValues[k] = append(newValues[k], vs...)
	}
	return &queryBindingSetter{
		values:   newValues,
		usedKeys: set.String{},
	}
}

func (s queryBindingSetter) empty() bool {
	return len(s.usedKeys) == len(s.values)
}

func (s *queryBindingSetter) set(v reflect.Value, sfs []reflect.StructField, sf reflect.StructField) (err error) {
	if s.empty() {
		return nil
	}
	if v.Kind() == reflect.Struct {
		if !v.CanSet() && !sf.Anonymous {
			return nil
		}
		for i := 0; i < v.NumField(); i++ {
			if err = s.set(v.Field(i), append(sfs, sf), v.Type().Field(i)); err != nil {
				return err
			}
		}
		return
	}
	if !v.CanSet() {
		return
	}
	key := queryBindingSetterGetKey(sfs, sf)
	values, exist := s.values[key]
	if !exist {
		return
	}
	var value interface{}
	value, err = queryBindingSetterConvertValues(values, sf.Type)
	if err != nil {
		return
	}
	s.usedKeys.Add(key)
	v.Set(reflect.ValueOf(value))
	return
}

func queryBindingSetterGetKey(sfs []reflect.StructField, sf reflect.StructField) string {
	sfs = append(sfs, sf)
	slice.Filter(&sfs, func(index int) bool {
		return !sfs[index].Anonymous
	})
	return strings.Join(slice.MapString(sfs, func(i int, v interface{}) string {
		return name.ToSnakeCase(v.(reflect.StructField).Name)
	}), ".")
}

func queryBindingSetterConvertValues(values []string, rt reflect.Type) (interface{}, error) {
	switch rt.Kind() { // Kind has a total of 27 values
	case reflect.Struct: // kind would never be Struct (1)
		panic("never reached code")
	case
		// no allowed kind (7):
		reflect.Invalid,
		reflect.Uintptr, reflect.UnsafePointer,
		reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		// unsupported kind (3):
		reflect.Array,
		reflect.Complex64, reflect.Complex128:
		return nil, fmt.Errorf("unsupported type: %v", rt)
	case reflect.Ptr: // pointer (1)
		rt = rt.Elem()
		ptr := reflect.New(rt)
		var v interface{}
		var err error
		if v, err = queryBindingSetterConvertValues(values, rt); err != nil {
			return nil, err
		}
		ptr.Elem().Set(reflect.ValueOf(v))
		return ptr.Interface(), nil
	// supported plural value kind (1):
	case reflect.Slice:
		list := reflect.MakeSlice(rt, len(values), len(values))
		var v interface{}
		var err error
		rt = rt.Elem()
		for i, value := range values {
			if v, err = queryBindingSetterConvertValue(value, rt); err != nil {
				return nil, err
			}
			list.Index(i).Set(reflect.ValueOf(v))
		}
		return list.Interface(), nil
	case // supported single value kind (14):
		reflect.Bool, reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		value := ""
		if len(values) > 0 {
			value = values[0]
		}
		return queryBindingSetterConvertValue(value, rt)
	}
	panic("never reached code")
}

func queryBindingSetterConvertValue(value string, rt reflect.Type) (interface{}, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.String:
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return queryBindingSetterConvertIntValue(value, rt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return queryBindingSetterConvertUintValue(value, rt)
	case reflect.Float32, reflect.Float64:
		return queryBindingSetterConvertFloatValue(value, rt)
	default:
		panic("never reached code")
	}
}

func queryBindingSetterConvertIntValue(value string, rt reflect.Type) (interface{}, error) {
	var bitSize int
	switch rt.Kind() {
	case reflect.Int, reflect.Int64:
		bitSize = 64
	case reflect.Int8:
		bitSize = 8
	case reflect.Int16:
		bitSize = 16
	case reflect.Int32:
		bitSize = 32
	default:
		panic("never reached code")
	}
	v, err := strconv.ParseInt(value, 10, bitSize)
	if err != nil {
		return nil, err
	}
	switch rt.Kind() {
	case reflect.Int:
		return int(v), nil
	case reflect.Int8:
		return int8(v), nil
	case reflect.Int16:
		return int16(v), nil
	case reflect.Int32:
		return int32(v), nil
	case reflect.Int64:
		return v, nil
	default:
		panic("never reached code")
	}
}

func queryBindingSetterConvertUintValue(value string, rt reflect.Type) (interface{}, error) {
	var bitSize int
	switch rt.Kind() {
	case reflect.Uint, reflect.Uint64:
		bitSize = 64
	case reflect.Uint8:
		bitSize = 8
	case reflect.Uint16:
		bitSize = 16
	case reflect.Uint32:
		bitSize = 32
	default:
		panic("never reached code")
	}
	v, err := strconv.ParseUint(value, 10, bitSize)
	if err != nil {
		return nil, err
	}
	switch rt.Kind() {
	case reflect.Uint:
		return uint(v), nil
	case reflect.Uint8:
		return uint8(v), nil
	case reflect.Uint16:
		return uint16(v), nil
	case reflect.Uint32:
		return uint32(v), nil
	case reflect.Uint64:
		return v, nil
	default:
		panic("never reached code")
	}
}

func queryBindingSetterConvertFloatValue(value string, rt reflect.Type) (interface{}, error) {
	var bitSize int
	switch rt.Kind() {
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	default:
		panic("never reached code")
	}
	v, err := strconv.ParseFloat(value, bitSize)
	if err != nil {
		return nil, err
	}
	switch rt.Kind() {
	case reflect.Float32:
		return float32(v), nil
	case reflect.Float64:
		return v, nil
	default:
		panic("never reached code")
	}
}
