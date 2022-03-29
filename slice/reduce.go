package slice

func Reduce[I, V any](slice []I, initReduceValue V, f func(reduceValue V, index int, value I) V) (reduceValue V) {
	reduceValue = initReduceValue
	for i, v := range slice {
		reduceValue = f(reduceValue, i, v)
	}
	return
}
