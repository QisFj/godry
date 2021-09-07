package registry

import (
	"fmt"
	"reflect"
)

// Registry provider Register and Get
// zero value is not ready for use, New it
type Registry struct {
	itemType         reflect.Type
	unregisterReturn func(name string) interface{}

	items map[string]interface{}
}

func New(
	itemType reflect.Type,
	unregisterReturn func(name string) interface{},
) *Registry {
	return &Registry{
		itemType:         itemType,
		unregisterReturn: unregisterReturn,

		items: map[string]interface{}{},
	}
}

// Register item
// panic for unacceptable type
// panic for duplicate name
func (r *Registry) Register(name string, item interface{}) {
	panicf := func(format string, v ...interface{}) string {
		panic(fmt.Sprintf("failed to register %s: %q, %s", r.itemType.Name(), name, fmt.Sprintf(format, v...)))
	}
	itemType := reflect.TypeOf(item)
	if r.itemType.Kind() == reflect.Interface {
		if !itemType.Implements(r.itemType) {
			panicf("unacceptable type: %s, not implement", r.itemType.Name())
		}
	} else {
		if !itemType.ConvertibleTo(r.itemType) {
			panicf("unacceptable type: %s, not convertable", r.itemType.Name())
		}
		item = reflect.ValueOf(item).Convert(r.itemType).Interface()
	}

	_, exist := r.items[name]
	if exist {
		panicf("duplicate name: %q", name)
	}
	r.items[name] = item
}

func (r *Registry) Get(name string) interface{} {
	i, exist := r.items[name]
	if exist {
		return i
	}
	if r.unregisterReturn != nil {
		return r.unregisterReturn(name)
	}
	return nil
}
