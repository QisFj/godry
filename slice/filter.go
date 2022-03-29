package slice

func Filter[I any](slice []I, f func(index int, value I) (keep bool)) []I {
	if slice == nil { // do nothing
		return nil
	}
	filtered := make([]I, 0, len(slice))
	for i := 0; i < len(slice); i++ {
		if f(i, slice[i]) {
			filtered = append(filtered, slice[i])
		}
	}
	return filtered
}

// FilterOn is in-place filter of a slice.
func FilterOn[I any](slice *[]I, f func(index int) (keep bool)) {
	if slice == nil { // do nothing
		return
	}
	var newLength int
	for i := 0; i < len(*slice); i++ {
		if f(i) {
			if i != newLength {
				// move slice[i] to slice[newLength]
				(*slice)[newLength] = (*slice)[i]
			}
			newLength++
		}
	}
	*slice = (*slice)[:newLength]
}
