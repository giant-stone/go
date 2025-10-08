# About

[![Go](https://github.com/giant-stone/go/actions/workflows/go.yml/badge.svg)](https://github.com/giant-stone/go/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/giant-stone/go)](https://goreportcard.com/report/github.com/giant-stone/go)
[![LICENSE](https://img.shields.io/github/license/giant-stone/go.svg?style=flat-square)](https://github.com/giant-stone/go/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/giant-stone/go?status.svg)](https://godoc.org/github.com/giant-stone/go)

giant-stone/go is a library that integrates frequently used functions from multiple production environment projects,
avoiding repetition in each project.

Build requirement:

- Go >= 1.20

Installing the latest version

    go list -m -versions github.com/giant-stone/go
    # Choose the latest version from the output, such as `v1.1.0`
    go get -u github.com/giant-stone/go@v1.1.0

For usage examples, see <https://github.com/giant-stone/go/wiki>

Updating ghttp mock code

    go install go.uber.org/mock/mockgen@v0.4.0
    mockgen -source=ghttp/ghttp.go -destination=ghttp/impl_mock.go -package=ghttp -mock_names Interface=ImplMock
