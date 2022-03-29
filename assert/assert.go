package assert

import (
	"fmt"
	"strings"
)

//go:generate go run type_assert.gen.go

type Assertion struct {
	Message string
	KVs     map[string]interface{}
}

func (a Assertion) String() string {
	sb := strings.Builder{}
	sb.WriteString(a.Message)
	for k, v := range a.KVs {
		sb.WriteString("\n\t* ")
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	return sb.String()
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
