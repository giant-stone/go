package ghttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/giant-stone/go/ghttp"
	"github.com/giant-stone/go/glogging"
	"github.com/giant-stone/go/gutil"
)

func ExampleNew() {
	glogging.Init([]string{"stderr"}, "debug")

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
	fmt.Println(req.RespStatus)
	// Output: 200
}

type FieldForm struct {
	Id   string
	Name string
}
type Rs struct {
	Form *FieldForm `json:"form"`
}

// ExampleHttpRequest_SetPostBody show howto POST in application/x-www-form-urlencoded
func ExampleHttpRequest_SetPostBody() {
	glogging.Init([]string{"stderr"}, "debug")

	rq := ghttp.New().
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
	gutil.ExitOnErr(err)

	want := &FieldForm{
		Id:   form.Get("id"),
		Name: form.Get("name"),
	}
	var rs Rs
	err = json.Unmarshal(rq.RespBody, &rs)
	gutil.ExitOnErr(err)

	fmt.Println(rq.RespStatus)
	fmt.Println(want.Id == rs.Form.Id && want.Name == rs.Form.Name)
	// Output:
	// 200
	// true
}

// ExampleHttpRequest_SetPostBody show howto POST in multipart/form-data
func ExampleHttpRequest_SetPostBody_multipart() {
	glogging.Init([]string{"stderr"}, "debug")
	// use default implement avoid conflict with ghttp/ghttp_test.go
	ghttp.UseImpl(nil)

	var err error

	gh := ghttp.New().SetTimeout(time.Second * 5)

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

	rq, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewReader(b.Bytes()))
	gutil.ExitOnErr(err)

	rq.Header.Add("Content-Type", w.FormDataContentType())
	rs, err := gh.Do(rq)
	gutil.ExitOnErr(err)

	rsBody, err := ghttp.ReadBody(rs)
	gutil.ExitOnErr(err)

	fmt.Println(rs.StatusCode)
	fmt.Println(len(rsBody) > 0)
	// Output:
	// 200
	// true
}
