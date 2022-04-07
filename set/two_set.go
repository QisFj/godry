package set

// two_set.go contain those function like func(s1, s2 Set) Set

func Union[V comparable](s1, s2 Set[V]) Set[V] {
	s1 = s1.Clone()
	s1.AddSet(s2)
	return s1
}

func Intersect[V comparable](s1, s2 Set[V]) Set[V] {
	s1 = s1.Clone()
	for v := range s1 {
		if !s2.Contains(v) {
			s1.Remove(v)
		}
	}
	return s1
}

func Subtract[V comparable](s1, s2 Set[V]) Set[V] {
	s1 = s1.Clone()
	s1.RemoveSet(s2)
	return s1
}

func Diff[V comparable](s1, s2 Set[V]) (both, only1, only2 Set[V]) {
	both, only1, only2 = Set[V]{}, Set[V]{}, Set[V]{}
	for v := range s1 {
		if s2.Contains(v) {
			both.Add(v)
		} else {
			only1.Add(v)
		}
	}
	for v := range s2 {
		if !s1.Contains(v) {
			only2.Add(v)
		}
		// since all elements in s1 already added, we don't need to add it twice
	}
	return
}
