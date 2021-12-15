package ghttp_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"time"

	"github.com/giant-stone/go/ghttp"
	"github.com/giant-stone/go/gutil"
)

func ExampleNew() {
	fullurl := "https://httpbin.org/post"
	postData := []byte(`{"msg":"hello"}`)
	req := ghttp.New(context.Background()).
		SetRandomUserAgent(true).
		SetTimeout(time.Second * 3).
		SetRequestMethod("POST").
		SetUri(fullurl).
		SetProxy(os.Getenv("HTTPS_PROXY")).
		SetPostBody(&postData)
	err := req.Send()
	ghttp.CheckRequestErr(fullurl, req.RespStatus, req.RespBody, err)
}

// ExampleHttpRequest_SetPostBody show howto POST in application/x-www-form-urlencoded
func ExampleHttpRequest_SetPostBody() {
	rq := ghttp.New(context.Background()).
		SetDebug(true).
		SetRequestMethod("POST").
		SetUri("https://httpbin.org/post").
		SetTimeout(time.Second * 5)

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

// ExampleHttpRequest_SetPostBody2 show howto POST in multipart/form-data
func ExampleHttpRequest_SetPostBody2() {
	var err error

	rq := ghttp.New(context.Background()).
		SetDebug(true).
		SetRequestMethod("POST").
		SetUri("https://httpbin.org/post").
		SetTimeout(time.Second * 5)

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	err = ghttp.AppendMultipartFormData(w, "myfile", "myfile.data", []byte(`hello 中文`))
	gutil.ExitOnErr(err)

	err = ghttp.AppendMultipartFormData(w, "myfile2", "myfile2.data", []byte(`foo\nbar`))
	gutil.ExitOnErr(err)

	err = w.WriteField("id", "123")
	gutil.ExitOnErr(err)

	err = w.Close()
	gutil.ExitOnErr(err)

	rqBody := b.Bytes()

	rq.SetPostBody(&rqBody)
	rq.SetHeader("Content-Type", w.FormDataContentType())
	err = rq.Send()
	gutil.ExitOnErr(err)

	log.Println(
		rq.RespStatus,
		string(rq.RespBody),
	)
}
