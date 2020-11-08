# Package retry

retry

## Usage

```go
package main

import (
	"errors"
	"log"
	"time"

	"github.com/QisFj/godry/retry"
)

func main() {
	r := retry.New(retry.Option{
		MaxRunTime:    5,
		RetryInterval: time.Second,
		F: func(r *retry.Retry) error {
			log.Printf("Run: %d\n", r.RunCount()+1)
			return errors.New("always run")
		},
		ResultSize: 3,
		ErrorWrapper: func(r *retry.Retry, err error) error {
			if err != nil { // error may be nil
				log.Printf("error: %s\n", err)
			}
			return err
		},
	})
	r.Start().Wait() // Start & Wait
	for i := 1; i <= 5; i++ {
		result := r.Result(i) // Start From 1
		if !result.Valid {
			log.Printf("Run %d: No Valid Result\n", i)
		} else if result.Error == nil {
			log.Printf("Run %d: Result: success\n", i)
		} else {
			log.Printf("Run %d: Result: error: %s\n", i, result.Error)
		}
	}
}
```