package uiter

import (
	"iter"
)

// Any is like any in python, return if any element in the slice satisfy the predicate.
func Any[T any](it iter.Seq[T], predicate func(T) bool) bool {
	for item := range it {
		if predicate(item) {
			return true
		}
	}
	return false
}

// All is like all in python, return if all elements in the slice satisfy the predicate.
func All[T any](it iter.Seq[T], predicate func(T) bool) bool {
	for item := range it {
		if !predicate(item) {
			return false
		}
	}
	return true
}
