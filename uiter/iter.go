package uiter

import (
	"github.com/samber/lo"
	"iter"
)

// Static return an iter.Seq return certain value
func Static[T any](v T) iter.Seq[T] {
	return func(yield func(T) bool) {
		yield(v)
	}
}

// Static2 return an iter.Seq2 return certain value
func Static2[K, V any](k K, v V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		yield(k, v)
	}
}

// Noop is an iter.Seq do nothing
func Noop[T any](_ func(T) bool) {}

// Noop2 is an iter.Seq2 do nothing
func Noop2[K, V any](_ func(K, V) bool) {}

// Chain chains multiple iter.Seq into one iter.Seq
func Chain[T any](its ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, it := range its {
			for v := range it {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Chain2 chains multiple iter.Seq2 into one iter.Seq2
func Chain2[K, V any](its ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, it := range its {
			for k, v := range it {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Keys converts an iter.Seq2 into an iter.Seq
func Keys[K, V any](it iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range it {
			if !yield(k) {
				return
			}
		}
	}
}

// Values converts an iter.Seq2 into an iter.Seq
func Values[K, V any](it iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range it {
			if !yield(v) {
				return
			}
		}
	}
}

// Enumerate is like enumerate in python
// convert's an iter.Seq into an iter.Seq2 by its index
func Enumerate[T any](it iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range it {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

// Filter filter out if pred(v) is false
func Filter[T any](it iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range it {
			if !pred(v) {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Filter2 filter out if pred(k, v) is false
func Filter2[K, V any](it iter.Seq2[K, V], pred func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range it {
			if !pred(k, v) {
				continue
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

func makeSlice[T any](capHint ...int) []T {
	if len(capHint) > 0 && capHint[0] > 0 {
		return make([]T, 0, capHint[0])
	}
	return nil
}

// Dump dumps an iter.Seq into a slice
func Dump[T any](it iter.Seq[T], capHint ...int) []T {
	ret := makeSlice[T](capHint...)
	for v := range it {
		ret = append(ret, v)
	}
	return ret
}

// Dump2 dumps an iter.Seq2 into a lo.Tuple2 slice
func Dump2[K, V any](it iter.Seq2[K, V], capHint ...int) []lo.Tuple2[K, V] {
	ret := makeSlice[lo.Tuple2[K, V]](capHint...)
	for k, v := range it {
		ret = append(ret, lo.T2(k, v))
	}
	return ret
}
