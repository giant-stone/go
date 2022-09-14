package ghttp

import "net/http"

const (
	HEAD   = "HEAD"
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

type HttpClient interface {
	Do(rq *http.Request) (rs *http.Response, err error)
}

var Client HttpClient

func UseImpl(impl HttpClient) {
	Client = impl
}
