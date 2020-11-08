package assert

import "fmt"

type Violation struct {
	Position  string
	Assertion Assertion
}

func (v Violation) String() string {
	return fmt.Sprintf("<assertion violation>: %s %s", v.Position, v.Assertion)
}

// this can be used to modify Violation object
type ViolationOpt func(violation *Violation)

func CallerOption(deltaDepth int, shortName bool) ViolationOpt {
	return func(violation *Violation) {
		violation.Position = lineNumFmt(lineNum(2+deltaDepth, shortName))
	}
}
