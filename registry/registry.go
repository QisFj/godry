package registry

import (
	"fmt"
	"strings"
)

// Registry provider Register and Get
// zero value is not ready for use, New it
//
// this is a map, but
// - can't Register a name twice, panic if it happens
// - can get a hook on Get an unregistered name
//
// after go support generics type, it seems not need to use this package.
type Registry[T any] struct {
	name             string
	unregisterReturn func(name string) T
	items            map[string]T
}

// AlwaysReturn the given value, it's an implementation of unregisterReturn
func AlwaysReturn[T any](v T) func(name string) T {
	return func(name string) T {
		return v
	}
}

func New[T any](
	name string,
	unregisterReturn func(name string) T,
) *Registry[T] {
	return &Registry[T]{
		name:             name,
		unregisterReturn: unregisterReturn,
		items:            map[string]T{},
	}
}

// Register item
func (r *Registry[T]) Register(name string, item T) {
	_, exist := r.items[name]
	if exist {
		msg := strings.Builder{}
		msg.WriteString("failed to register")
		if r.name != "" {
			msg.WriteString(fmt.Sprintf("(registry=%s)", r.name))
		}
		msg.WriteString(fmt.Sprintf(" %q", name))
		msg.WriteString(": duplicate name")
		panic(msg.String())
	}
	r.items[name] = item
}

func (r *Registry[T]) Get(name string) T {
	i, exist := r.items[name]
	if exist {
		return i
	}
	if r.unregisterReturn != nil {
		return r.unregisterReturn(name)
	}
	return i // i should be zero value to T
}
