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
func FilterOn[I any](slice *[]I, f func(index int, value I) (keep bool)) {
	if slice == nil { // do nothing
		return
	}
	// see: https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	newSlice := (*slice)[:0]
	for i, v := range *slice {
		if f(i, v) {
			newSlice = append(newSlice, v)
		}
	}
	*slice = newSlice
}
