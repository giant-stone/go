package ghttp_test

import (
	"os"
	"time"

	"github.com/giant-stone/go/ghttp"
)

func ExampleNew() {
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
	ghttp.CheckRequestErr(fullurl, req.RespStatus, req.RespBody, err)
}
