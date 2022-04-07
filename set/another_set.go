package set

// another_set.go contains thous functions like (s Set)func(another Set)

func (s Set[V]) AddSet(another Set[V]) {
	// NOTE: do not use maps.Copy here
	// when V is a struct which contain more than one filed, use maps.Copy leads compile error
	for v := range another {
		s[v] = struct{}{}
	}
}

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
