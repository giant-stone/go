# About

![Go](https://github.com/giant-stone/go/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/giant-stone/go)](https://goreportcard.com/report/github.com/giant-stone/go)
[![LICENSE](https://img.shields.io/github/license/giant-stone/go.svg?style=flat-square)](https://github.com/giant-stone/go/blob/master/LICENSE)

giant-stone/go is a Go library which provides utility functions for common programming tasks.

giant-stone/go 是一个将多个生产环境项目高频使用函数整合一起，避免在每个项目中不断重复。

安装最新版本

    go list -m -versions github.com/giant-stone/go
	# 在输出中选择最新一个版本，比如 `v0.0.14`
    go get -u github.com/giant-stone/go@v0.0.14

使用示例见 https://github.com/giant-stone/go/wiki


更新 ghttp mock 代码

    mockgen -source=ghttp/ghttp.go -destination=ghttp/httpclitest.go -package=ghttp
