package set

import (
	"reflect"
)

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

func (set *T) remove(e reflect.Value) {
	set.m.SetMapIndex(e, reflect.Value{})
}

func (set *T) Remove(es ...interface{}) {
	for _, e := range es {
		set.remove(reflect.ValueOf(e))
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

func (set T) Clone() T {
	clone := NewT(set.t)
	for it := set.m.MapRange(); it.Next(); {
		clone.add(it.Key())
	}
	return clone
}

func (set *T) AddSet(another T) {
	if set.t != another.t {
		panic("two sets must be same type")
	}
	for it := another.m.MapRange(); it.Next(); {
		set.add(it.Key())
	}
}

func (set *T) RemoveSet(another T) {
	if set.t != another.t {
		panic("two sets must be same type")
	}
	for it := another.m.MapRange(); it.Next(); {
		set.remove(it.Key())
	}
}

func TUnion(s1, s2 T) T {
	s1 = s1.Clone()
	s1.AddSet(s2)
	return s1
}

func TIntersect(s1, s2 T) T {
	s1 = s1.Clone()
	for _, v := range s1.m.MapKeys() {
		if !s2.contains(v) {
			s1.remove(v)
		}
	}
	return s1
}

func TSubtract(s1, s2 T) T {
	s1 = s1.Clone()
	s1.RemoveSet(s2)
	return s1
}

func TDiff(s1, s2 T) (both, only1, only2 T) {
	if s1.t != s2.t {
		panic("two sets must be same type")
	}
	both, only1, only2 = NewT(s1.t), NewT(s1.t), NewT(s1.t)
	for _, v := range s1.m.MapKeys() {
		if s2.contains(v) {
			both.add(v)
		} else {
			only1.add(v)
		}
	}
	for _, v := range s2.m.MapKeys() {
		if !s1.contains(v) {
			only2.add(v)
		}
		// since all elements in s1 already added, we don't need to add it twice
	}
	return
}
