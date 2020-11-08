# Package multierr

multiple error.

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/QisFj/godry/multierr"
)

func main() {
	errs := multierr.New(func(err error) error {
		return fmt.Errorf("-> %w", err)
	}, nil)
	for i := 0; i < 6; i++ {
		var err error
		if i%2 == 0 {
			err = fmt.Errorf("error-%d", i+1)
		}
		errs.AppendOnlyNotNil(err)
	}
	log.Printf("err: %s\n", errs.Error())
}
```