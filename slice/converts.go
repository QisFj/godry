package slice

// The functions in this file(converts.go) are related to the type conversion between map and slice

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

// KVsOfMap return a kv slice from map
func KVsOfMap[K comparable, V any](m map[K]V) []KV[K, V] {
	// m must be a map, result's Key is map's key, Value is map's value
	kvs := make([]KV[K, V], 0, len(m))
	for k, v := range m {
		kvs = append(kvs, KV[K, V]{Key: k, Value: v})
	}
	return kvs
}

// ToMap convert a slice to a map
// if vf is nil, will try to use elem of slice as value
func ToMap[I any, K comparable, V any](slice []I, kf func(I) K, vf func(I) V) map[K]V {
	m := make(map[K]V)
	for _, i := range slice {
		m[kf(i)] = vf(i)
	}
	return m
}

// ToMapItself convert a slice to a map
// like ToMap, but the value of map is the elem of slice itself
func ToMapItself[I any, K comparable](slice []I, kf func(I) K) map[K]I {
	m := make(map[K]I)
	for _, i := range slice {
		m[kf(i)] = i
	}
	return m
}
