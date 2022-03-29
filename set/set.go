package set

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
)

type Set[V constraints.Ordered] map[V]struct{}

func Of[V constraints.Ordered](values ...V) Set[V] {
	set := Set[V]{}
	set.Add(values...)
	return set
}

func (s *Set[V]) Add(values ...V) {
	for _, v := range values {
		(*s)[v] = struct{}{}
	}
}

func (s *Set[V]) Remove(values ...V) {
	for _, v := range values {
		delete(*s, v)
	}
}

func (s Set[V]) Contains(value V) bool {
	_, ok := s[value]
	return ok
}

func (s Set[V]) ContainsAll(values ...V) bool {
	for _, v := range values {
		if !s.Contains(v) {
			return false
		}
	}
	return true
}

func (s Set[V]) ContainsAny(values ...V) bool {
	for _, v := range values {
		if s.Contains(v) {
			return true
		}
	}
	return false
}

func (s Set[V]) List() []V { return maps.Keys(s) }

func (s Set[V]) Clone() Set[V] { return maps.Clone(s) }
