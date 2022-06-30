package concurrent

type Func func()
type FuncMayError func() error
type FuncWithResult func() (result interface{})
type FuncWithResultMayError func() (result interface{}, err error)
