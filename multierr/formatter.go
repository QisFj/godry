package multierr

import (
	"fmt"
	"strings"

	"github.com/QisFj/godry/slice"
)

type Formatter func(errs []error) string

// FormatterList is the default Formatter
func FormatterList(errs []error) string {
	if len(errs) == 0 {
		return "no error occurred"
	}
	if len(errs) == 1 {
		return fmt.Sprintf("1 error occurred: %s", errs[0])
	}
	return fmt.Sprintf("%d errors occurred:%s\n", len(errs),
		strings.Join(slice.Map(errs, func(_ int, err error) string {
			return fmt.Sprintf("\n\t* %s", err)
		}), ""),
	)
}
