# HTTP Binding

## QueryBiding

Example:

```go
package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/QisFj/godry/http/binding"
)

func main() {
	req := &http.Request{URL: &url.URL{RawQuery: "a=1&b.c=2"}}
	var obj struct {
		A int
		B struct {
			C int
		}
	}
	err := binding.QueryBinding{}.Bind(req, &obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Resp: %#v\n", obj)
}
```