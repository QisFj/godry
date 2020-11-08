package assert

import (
	"fmt"
	"runtime"
	"strings"
)

// stack info, filename & line number
func lineNum(depth int, shortName bool) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(1 + depth)
	if !ok {
		file = "???"
	}
	if shortName {
		lastIndex := strings.LastIndex(file, "/")
		if lastIndex != -1 {
			file = file[lastIndex+1:]
		}
	}
	return file, line
}

func lineNumFmt(file string, line int) string {
	return fmt.Sprintf("%s:%d", file, line)
}
