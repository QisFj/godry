package set

import "golang.org/x/exp/maps"

// another_set.go contains thous functions like (s Set)func(another Set)

func (s Set[V]) AddSet(another Set[V]) { maps.Copy(s, another) }

func (s Set[V]) RemoveSet(another Set[V]) {
	for v := range another {
		delete(s, v)
	}
}

func (s Set[V]) ContainsAnyOfSet(another Set[V]) bool {
	for v := range another {
		if _, ok := s[v]; ok {
			return true
		}
	}
	return false
}

func (s Set[V]) ContainsAllOfSet(another Set[V]) bool {
	for v := range another {
		if _, ok := s[v]; !ok {
			return false
		}
	}
	return true
}
