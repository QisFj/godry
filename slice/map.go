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
