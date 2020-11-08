package assert

import (
	"fmt"

	"github.com/QisFj/godry/slice"
)

//go:generate go run type_assert.gen.go

type Assertion struct {
	Message string
	KVs     map[string]interface{}
}

func (a Assertion) String() string {
	return a.Message + slice.Reduce(slice.KVsOfMap(a.KVs), "", func(reduceValue interface{}, i int, v interface{}) interface{} {
		kv := v.(slice.KV)
		return reduceValue.(string) + fmt.Sprintf("\n\t* %s=%v", kv.Key, kv.Value)
	}).(string)
}

func Assert(assert bool, message string, kvs map[string]interface{}, opts ...ViolationOpt) {
	if assert {
		return
	}
	violation := Violation{
		Position: lineNumFmt(lineNum(1, false)),
		Assertion: Assertion{
			Message: message,
			KVs:     kvs,
		},
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&violation)
	}
	panic(violation)
}

func Catch(handlers ...func(Violation)) {
	if p := recover(); p != nil {
		if violation, ok := p.(Violation); ok {
			for _, handler := range handlers {
				if handler == nil {
					continue
				}
				handler(violation)
			}
			return
		}
		panic(p)
	}
}
