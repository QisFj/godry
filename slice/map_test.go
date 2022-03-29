package slice

import (
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	timeStrings := []string{
		"1970-01-01 00:00:00",
		"1970-01-01 00:00:01",
		"1970-01-01 00:00:02",
	}
	mapF := func(i int, v string) time.Time {
		tm, _ := time.Parse("2006-01-02 15:04:05", v)
		return tm
	}
	for _, tm := range Map(timeStrings, mapF) {
		t.Logf("time: %s", tm)
	}
}
