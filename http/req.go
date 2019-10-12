package ghttp

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/giant-stone/go/util"
)

type HttpResponse struct {
	http.Response

	Body    *[]byte
	Elapsed time.Duration
}

type HttpRequest struct {
	Timeout time.Duration

	Method  string
	Uri     string
	Headers *map[string]interface{}
	Body    *[]byte

	UseRandomUserAgent bool
	UserAgent          string

	Proxy string

	Response *HttpResponse
}

func New() *HttpRequest {
	return &HttpRequest{
		Headers: &map[string]interface{}{},
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
	its.Body = body
	return its
}

func (its *HttpRequest) SetHttpAuth(username, password string) *HttpRequest {
	plain := fmt.Sprintf("%s:%s", username, password)
	value := "Basic " + base64.StdEncoding.EncodeToString([]byte(plain))
	its.SetHeader("Authorization", value)
	return its
}

func (its *HttpRequest) SetHeader(key string, value interface{}) *HttpRequest {
	(*its.Headers)[key] = value
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
	if its.Body != nil && len(*its.Body) > 0 {
		reqBody = bytes.NewBuffer(*its.Body)
	}
	req, err := http.NewRequest(its.Method, its.Uri, reqBody)
	if err != nil {
		return
	}

	if its.Headers != nil {
		for k, v := range *its.Headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	if its.UseRandomUserAgent {
		req.Header.Set("User-Agent", gutil.RandChoice(UserAgents).(string))
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		its.Response = &HttpResponse{
			Body:    nil,
			Elapsed: elapsed,
		}
		return
	}

	its.Response = &HttpResponse{
		Response: *resp,
		Body:     nil,
		Elapsed:  elapsed,
	}

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		its.Response.Body = &bodyResp
	}

	return
}
