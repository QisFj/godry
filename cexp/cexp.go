package cexp

//go:generate go run cexp.gen.go

type InterfaceGetter func() interface{}

func Interface(condition bool, v1, v2 interface{}) interface{} {
	if condition {
		return v1
	}
	return v2
}

func InterfaceShortCircuit(condition bool, g1, g2 InterfaceGetter) interface{} {
	if condition {
		return g1()
	}
	return g2()
}
