package ghttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/giant-stone/go/gutil"
	"github.com/giant-stone/go/logger"
	"github.com/giant-stone/go/utilhuman"
)

const (
	defaultMethod = "GET"
)

type HttpRequest struct {
	Ctx     context.Context
	Timeout time.Duration

	Method  string
	Uri     string
	Headers map[string]interface{}
	Body    []byte

	UseRandomUserAgent bool
	UserAgent          string

	Proxy string

	// request send at timestamp in unix(milliseconds)
	Rqts int64

	RespStatus int
	RespBody   []byte
	RespHeader http.Header
	Elapsed    time.Duration
}

func New() *HttpRequest {
	return &HttpRequest{
		Ctx:     context.Background(),
		Method:  defaultMethod,
		Headers: map[string]interface{}{},
		Timeout: time.Duration(10) * time.Second,
	}
}

func NewWithCtx(ctx context.Context) *HttpRequest {
	return &HttpRequest{
		Ctx:     ctx,
		Method:  defaultMethod,
		Headers: map[string]interface{}{},
		Timeout: time.Duration(10) * time.Second,
	}
}

func (its *HttpRequest) SetRandomUserAgent(flag bool) *HttpRequest {
	its.UseRandomUserAgent = flag
	return its
}

func (its *HttpRequest) SetTimeout(duration time.Duration) *HttpRequest {
	its.Timeout = duration
	return its
}

func (its *HttpRequest) SetRequestMethod(method string) *HttpRequest {
	its.Method = method
	return its
}

func (its *HttpRequest) SetUri(uri string) *HttpRequest {
	its.Uri = uri

	return its
}

func (its *HttpRequest) SetProxy(addr string) *HttpRequest {
	its.Proxy = addr
	return its
}

func (its *HttpRequest) SetPostBody(body *[]byte) *HttpRequest {
	if body != nil && len(*body) > 0 {
		its.Body = make([]byte, len(*body))
		copy(its.Body, *body)
	}
	return its
}

func (its *HttpRequest) SetHttpAuth(username, password string) *HttpRequest {
	plain := fmt.Sprintf("%s:%s", username, password)
	value := "Basic " + base64.StdEncoding.EncodeToString([]byte(plain))
	its.SetHeader("Authorization", value)
	return its
}

func (its *HttpRequest) SetHeader(key string, value interface{}) *HttpRequest {
	its.Headers[key] = value
	return its
}

func (its *HttpRequest) SetHeaders(headers map[string]interface{}) *HttpRequest {
	for key, value := range headers {
		its.Headers[key] = value
	}
	return its
}

func (its *HttpRequest) Send() (err error) {
	client := &http.Client{
		Timeout: its.Timeout,
	}

	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if its.Proxy != "" {
		proxyNode := its.Proxy
		if !strings.HasPrefix(proxyNode, "http") {
			proxyNode = fmt.Sprintf("http://%s", proxyNode)
		}
		logger.Sugared.Infof("proxy=%s", proxyNode)

		u, errUrl := url.Parse(proxyNode)
		if errUrl != nil {
			err = errUrl
			return
		}
		tr.Proxy = http.ProxyURL(u)
		client.Transport = &tr
	}
	client.Transport = &tr

	var reqBody io.Reader
	if len(its.Body) > 0 {
		reqBody = bytes.NewBuffer(its.Body)
	}

	req, err := http.NewRequestWithContext(its.Ctx, its.Method, its.Uri, reqBody)
	if err != nil {
		return
	}

	for k, v := range its.Headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	if req.Header.Get("User-Agent") == "" {
		if its.UseRandomUserAgent {
			req.Header.Set("User-Agent", gutil.RandChoice(UserAgents).(string))
		}
	}

	now := time.Now()
	its.Rqts = now.UnixNano() / 1000000
	resp, err := client.Do(req)
	elapsed := time.Since(now)

	logger.Sugared.Infof("%s %s elapsed=%v err=%v", its.Method, its.Uri, utilhuman.FmtDuration(elapsed), err)

	if err != nil {
		its.Elapsed = elapsed
		return
	}
	defer resp.Body.Close()

	its.RespStatus = resp.StatusCode
	its.RespHeader = resp.Header

	RespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	its.RespBody = RespBody
	return
}
