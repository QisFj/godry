# HTTP 

## QueryBinding

## Request

Example:

```go
package main

import (
    "fmt"
    "net/url"
    
    "github.com/QisFj/godry/http"
)

func main(){
    var req interface{}     // request body = json.Marshal(req)
    var resp interface{}    // json.Unmarshal(response body, &resp)
    err := new(http.Request).With(
        http.POST{},
        http.StatusCodeCheck{
            // if status code != 200 ; return error
            Equal:200,
        },
    ).Do("http://example.com", url.Values{
        "params": {"1", "2"},
    }, req, &resp)
    if err != nil {
        fmt.Println(err)
        return 
    }
    fmt.Printf("Resp: %#v\n", resp)
}
```