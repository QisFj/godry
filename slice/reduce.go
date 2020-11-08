package slice

func Reduce(slice interface{}, initReduceValue interface{}, f func(reduceValue interface{}, i int, v interface{}) interface{}) (reduceValue interface{}) {
	reduceValue = initReduceValue
	Foreach(slice, func(i int, v interface{}) {
		reduceValue = f(reduceValue, i, v)
	})
	return
}
