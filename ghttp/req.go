package ghttp

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/giant-stone/go/gutil"
)

type HttpRequest struct {
	Debug   bool
	Timeout time.Duration

	Method  string
	Uri     string
	Headers map[string]interface{}
	Body    []byte

	UseRandomUserAgent bool
	UserAgent          string

	Proxy string

	RespStatus int
	RespBody   []byte
	RespHeader http.Header
	Elapsed    time.Duration
}

func New() *HttpRequest {
	return &HttpRequest{
		Method:  "GET",
		Headers: map[string]interface{}{},
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

		if its.Debug {
			log.Printf("[debug] proxy=%s", proxyNode)
		}

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

	if its.Debug {
		log.Printf("[debug] %s %s", its.Method, its.Uri)
	}

	req, err := http.NewRequest(its.Method, its.Uri, reqBody)
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

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

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

func (its *HttpRequest) SetDebug(debug bool) *HttpRequest {
	its.Debug = debug
	return its
}
