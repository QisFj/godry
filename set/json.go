package set

import (
	"encoding/json"
	"sort"
)

// MarshalOrder can change the order of the elements when marshaling.
// 1 for desc
// 2 for asc
// 0 or other value for no certain order
var MarshalOrder = 0

func (s Set[V]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	return json.Marshal(s.List())
}

func (s *Set[V]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		(*s) = nil // set the map to nil, not pointer
		return nil
	}
	var l []V
	if err := json.Unmarshal(data, &l); err != nil {
		return err
	}
	(*s) = Of(l...)
	return nil
}

func (s SortableSet[V]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	list := Set[V](s).List()
	switch MarshalOrder {
	case 1:
		sort.Slice(list, func(i, j int) bool { return list[i] < list[j] }) // ensure return with certain order
	case 2:
		sort.Slice(list, func(i, j int) bool { return list[i] > list[j] }) // ensure return with certain order
	}
	return json.Marshal(list)
}

func (s *SortableSet[V]) UnmarshalJSON(data []byte) error {
	if len(data) == 4 && string(data) == "null" {
		(*s) = nil // set the map to nil, not pointer
		return nil
	}
	var l []V
	if err := json.Unmarshal(data, &l); err != nil {
		return err
	}
	(*s) = SortableSet[V](Of(l...))
	return nil
}
