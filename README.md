# About

![Go](https://github.com/giant-stone/go/workflows/Go/badge.svg?branch=master)

giant-stone/go is a Go library which provides utility functions for common programming tasks.
Confirm to https://en.wikipedia.org/wiki/Don%27t_repeat_yourself

*Life is short, don't repeat yourself.*



ghttp - HTTP client wrapper in Method chaining   
gsql - SQL CRUD and search wrapper


## Examples

### Send a HTTP reqeust in ghttp 

Custom HTTP request timeout, method, proxy and body in [Method chaining](https://en.wikipedia.org/wiki/Method_chaining)

```
package main


import (
	"log"
	"os"
	"time"

	"github.com/giant-stone/go/ghttp"
)

func main() {
	fullurl := "https://httpbin.org/post"
	postData := []byte(`{"msg":"hello"}`)
	req := ghttp.New().
		SetRandomUserAgent(true).
		SetTimeout(time.Second * 3).
		SetRequestMethod("POST").
		SetUri(fullurl).
		SetProxy(os.Getenv("HTTPS_PROXY")).
		SetPostBody(&postData)
	err := req.Send()
	if ghttp.CheckRequestErr(fullurl, req.RespStatus, req.RespBody, err) {
		log.Println("handler error ...")
	}

	log.Println("process response body ...", len(req.RespBody))
}
```
