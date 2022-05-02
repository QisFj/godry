package slice

func Map[I, O any](input []I, f func(index int, value I) O) []O {
	if input == nil {
		return nil
	}
	output := make([]O, len(input))
	for i, v := range input {
		output[i] = f(i, v)
	}
	return output
}

type BidirectionalManyToManyMap[K, V comparable] struct {
	k2v map[K]map[V]struct{} // avoid import set
	v2k map[V]map[K]struct{} // avoid import set
}

func (m *BidirectionalManyToManyMap[K, V]) Insert(k K, v V) {
	if m.k2v == nil {
		m.k2v = map[K]map[V]struct{}{}
	}
	if m.k2v[k] == nil {
		m.k2v[k] = map[V]struct{}{}
	}
	m.k2v[k][v] = struct{}{}
	if m.v2k == nil {
		m.v2k = map[V]map[K]struct{}{}
	}
	if m.v2k[v] == nil {
		m.v2k[v] = map[K]struct{}{}
	}
	m.v2k[v][k] = struct{}{}
}

func (m *BidirectionalManyToManyMap[K, V]) Keys() []K {
	if m.k2v == nil {
		return nil
	}
	keys := make([]K, 0, len(m.k2v))
	for k := range m.k2v {
		keys = append(keys, k)
	}
	return keys
}

func (m *BidirectionalManyToManyMap[K, V]) Values(k K) []V {
	if m.k2v == nil {
		return nil
	}
	values := make([]V, 0, len(m.v2k))
	for v := range m.v2k {
		values = append(values, v)
	}
	return values
}

// should not write to returned map
func (m *BidirectionalManyToManyMap[K, V]) KeysOf(v V) map[K]struct{} { return m.v2k[v] }

// should not write to returned map
func (m *BidirectionalManyToManyMap[K, V]) ValuesOf(k K) map[V]struct{} { return m.k2v[k] }

func (m *BidirectionalManyToManyMap[K, V]) DeleteKey(k K) {
	values := m.k2v[k]
	delete(m.k2v, k)
	for v := range values {
		delete(m.v2k[v], k)
		if len(m.v2k[v]) == 0 {
			delete(m.v2k, v)
		}
	}
}

func (m *BidirectionalManyToManyMap[K, V]) DeleteValue(v V) {
	keys := m.v2k[v]
	delete(m.v2k, v)
	for k := range keys {
		delete(m.k2v[k], v)
		if len(m.k2v[k]) == 0 {
			delete(m.k2v, k)
		}
	}
}
