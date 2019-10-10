package ghttp_test

import (
	"log"
	"time"

	"github.com/giant-stone/Go/http"
)

func ExampleNew() {
	postData := []byte(`{"msg":"hello"}`)
	req := ghttp.New().
		SetRandomUserAgent(true).
		SetTimeout(time.Second*3).
		SetRequestMethod("POST").
		SetUri("http://httpbin.org/post").
		SetProxy("http://172.21.27.11:8118").
		SetPostBody(&postData)
	err := req.Send()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(req.Response.Status, len(*req.Response.Body), string(*req.Response.Body))
	}

}
