# About

![Go](https://github.com/giant-stone/go/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/giant-stone/go)](https://goreportcard.com/report/github.com/giant-stone/go)
[![LICENSE](https://img.shields.io/github/license/giant-stone/go.svg?style=flat-square)](https://github.com/giant-stone/go/blob/master/LICENSE)

giant-stone/go is a Go library which provides utility functions for common programming tasks.
Confirm to https://en.wikipedia.org/wiki/Don%27t_repeat_yourself

_Life is short, don't repeat yourself._

giant-stone/go 是一个将多个生产环境项目高频使用函数整合一起，避免再每个项目中不断重复。
人生苦短。

## Modules

    ghttp - HTTP client wrapper in Method chaining.
    gtime - parse timestamp into YYYY-MM-DD in UTC and vise versa
    gstr  - strconv, crypto and unicode shortcut functions
    logger - custom logging level and logrotate

## Examples

### Custom logging

custom logging level and logrotate
自定义日志级别和日志切割（默认 100 MB 一个、保留 30 天，最多保留 15 个）

```
package main

import (
	"github.com/giant-stone/go/logger"
)

func main() {
	logger.Init([]string{"stderr"}, "debug")
	// or logger.Init([]string{"/data/foo/main.log"}, "warn")

	logger.Sugared.Debug("hello")
	logger.Sugared.Error("hello")
	logger.Sugared.Warnf("hello %s", "world")
}
```

### Send a HTTP GET request in ghttp

Custom HTTP request timeout, method, proxy and body in [Method chaining](https://en.wikipedia.org/wiki/Method_chaining)

自定义 HTTP 请求超的时、方法和 HTTP 代理

```
package main


import (
	"log"
	"os"
	"time"

	"github.com/giant-stone/go/logger"
	"github.com/giant-stone/go/ghttp"
)

func main() {
	logger.Init(nil, "")

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

### Send a POST multipart/form-data request in ghttp

自定义带文件表单的 HTTP 请求

```
	var err error

	rq := ghttp.New().
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
```
