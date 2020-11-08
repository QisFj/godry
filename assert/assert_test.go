package assert

import (
	"testing"
)

func TestAssert(t *testing.T) {
	var i interface{}
	i = "1"
	defer Catch(func(violation Violation) {
		t.Logf("Cacth Violation: %s", violation)
	})
	_ = Int(i)
	t.Errorf("not panic")
}
