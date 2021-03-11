# HTTP Request

Example:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/QisFj/godry/http/request"
)

func main() {
	ctx := context.Background()
	var req interface{}  // request body = json.Marshal(req)
	var resp interface{} // json.Unmarshal(response body, &resp)
	err := request.New().With(
		request.Options.Method(http.MethodPost),
		request.Options.CheckResponseBeforeUnmarshal(func(statusCode int, body []byte) error {
			if statusCode != http.StatusOK {
				return fmt.Errorf("http status: %d - %s", statusCode, http.StatusText(statusCode))
			}
			return nil
		}),
	).DoCtx(ctx, "http://example.com", url.Values{
		"params": {"1", "2"},
	}, req, &resp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Resp: %#v\n", resp)
}
```