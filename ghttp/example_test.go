package ghttp_test

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/giant-stone/go/ghttp"
	"github.com/giant-stone/go/gutil"
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


func ExamplePostInUrlencoded() {
	rq := ghttp.New().
		SetDebug(true).
		SetRequestMethod("POST").
		SetUri("https://httpbin.org/post").
		SetTimeout(time.Second*5)

	form := url.Values{}
	form.Add("id", fmt.Sprintf("%d", 123))
	form.Add("name", "foo")
	rq.SetHeader("Content-Type", "application/x-www-form-urlencoded")

	rqBody := []byte(form.Encode())

	rq.SetPostBody(&rqBody)
	err := rq.Send()
	gutil.CheckErr(err)

	log.Println(
		rq.RespStatus,
		string(rq.RespBody),
	)
}
