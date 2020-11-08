package multierr

import (
	"fmt"
	"strings"

	"github.com/QisFj/godry/slice"
)

type Formatter func(errs []error) string

// Default Formatter
func FormatterList(errs []error) string {
	if len(errs) == 0 {
		return "no error occurred"
	}
	if len(errs) == 1 {
		return fmt.Sprintf("1 error occurred: %s", errs[0])
	}
	return fmt.Sprintf("%d errors occurred:%s\n", len(errs),
		strings.Join(slice.MapString(errs, func(i int, v interface{}) string {
			return fmt.Sprintf("\n\t* %s", errs[i])
		}), ""),
	)
}
