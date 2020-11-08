# Package concurrent

run function concurrently.

## Usage

```go
package main

import (
	"errors"
	"log"
	"time"

	"github.com/QisFj/godry/concurrent"
)

func main() {
	log.Println("with error:")
	withErr()
	log.Println("no error:")
	noErr()
}

func withErr() {
	if err := concurrent.Do(
		func() error {
			time.Sleep(20 * time.Millisecond)
			return errors.New("error-1")
		},
		func() error {
			time.Sleep(20 * time.Millisecond)
			return errors.New("error-1")
		},
		func() error {
			time.Sleep(20 * time.Millisecond)
			return errors.New("error-2")
		},
	); err != nil {
		log.Printf("err: %s\n", err)
	} else {
		log.Println("no error")
	}
}

func noErr() {
	if err := concurrent.Do(
		func() error {
			time.Sleep(20 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(20 * time.Millisecond)
			return nil
		},
		func() error {
			time.Sleep(20 * time.Millisecond)
			return nil
		},
	); err != nil {
		log.Printf("err: %s\n", err)
	} else {
		log.Println("no error")
	}
}
```