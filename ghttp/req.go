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

	"github.com/giant-stone/go/ghuman"
	"github.com/giant-stone/go/glogging"
	"github.com/giant-stone/go/gutil"
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
	Rqts    int64
	Elapsed time.Duration

	// DEPRECATED.
	RespStatus int
	// DEPRECATED.
	RespBody []byte
	// DEPRECATED.
	RespHeader http.Header
}

func New() *HttpRequest {
	return &HttpRequest{
		Ctx:     context.Background(),
		Method:  "",
		Headers: map[string]interface{}{},
		Timeout: time.Duration(10) * time.Second,
	}
}

func NewWithCtx(ctx context.Context) *HttpRequest {
	return &HttpRequest{
		Ctx:     ctx,
		Method:  "",
		Headers: map[string]interface{}{},
		Timeout: time.Duration(10) * time.Second,
	}
}

func (it *HttpRequest) SetRandomUserAgent(flag bool) *HttpRequest {
	it.UseRandomUserAgent = flag
	return it
}

func (it *HttpRequest) SetTimeout(duration time.Duration) *HttpRequest {
	it.Timeout = duration
	return it
}

func (it *HttpRequest) SetRequestMethod(method string) *HttpRequest {
	it.Method = method
	return it
}

func (it *HttpRequest) SetUri(uri string) *HttpRequest {
	it.Uri = uri

	return it
}

func (it *HttpRequest) SetProxy(addr string) *HttpRequest {
	it.Proxy = addr
	return it
}

func (it *HttpRequest) SetPostBody(body *[]byte) *HttpRequest {
	if body != nil && len(*body) > 0 {
		it.Body = make([]byte, len(*body))
		copy(it.Body, *body)
	}
	return it
}

func (it *HttpRequest) SetHttpAuth(username, password string) *HttpRequest {
	plain := fmt.Sprintf("%s:%s", username, password)
	value := "Basic " + base64.StdEncoding.EncodeToString([]byte(plain))
	it.SetHeader("Authorization", value)
	return it
}

func (it *HttpRequest) SetHeader(key string, value interface{}) *HttpRequest {
	it.Headers[key] = value
	return it
}

func (it *HttpRequest) SetHeaders(headers map[string]interface{}) *HttpRequest {
	for key, value := range headers {
		it.Headers[key] = value
	}
	return it
}

// DEPRECATED.
//   `Send()` does not supports gomock, use `Do()` instead.
func (it *HttpRequest) Send() (err error) {
	client := &http.Client{
		Timeout: it.Timeout,
	}

	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if it.Proxy != "" {
		proxyNode := it.Proxy
		if !strings.HasPrefix(proxyNode, "http") {
			proxyNode = fmt.Sprintf("http://%s", proxyNode)
		}
		glogging.Sugared.Infof("proxy=%s", proxyNode)

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
	if len(it.Body) > 0 {
		reqBody = bytes.NewBuffer(it.Body)
	}

	if it.Method == "" {
		it.Method = defaultMethod
	}

	req, err := http.NewRequestWithContext(it.Ctx, it.Method, it.Uri, reqBody)
	if err != nil {
		return
	}

	for k, v := range it.Headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	if req.Header.Get("User-Agent") == "" {
		if it.UseRandomUserAgent {
			req.Header.Set("User-Agent", gutil.RandChoice(UserAgents).(string))
		}
	}

	now := time.Now()
	it.Rqts = now.UnixNano() / 1000000
	resp, err := client.Do(req)
	elapsed := time.Since(now)

	glogging.Sugared.Infof("%s %s elapsed=%v err=%v", it.Method, it.Uri, ghuman.FmtDuration(elapsed), err)

	if err != nil {
		it.Elapsed = elapsed
		return
	}
	defer resp.Body.Close()

	it.RespStatus = resp.StatusCode
	it.RespHeader = resp.Header

	RespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	it.RespBody = RespBody
	return
}

func (it *HttpRequest) GenerateRequest() (rs *http.Request) {
	var rqBody io.Reader
	if len(it.Body) > 0 {
		rqBody = io.NopCloser(bytes.NewReader(it.Body))
	}

	rs, _ = http.NewRequestWithContext(it.Ctx, it.Method, it.Uri, rqBody)

	for k, v := range it.Headers {
		value := fmt.Sprintf("%v", v)
		rs.Header.Set(k, value)
		glogging.Sugared.Debugf("Header.Set %s=%s", k, value)
	}

	if rs.Header.Get("User-Agent") == "" {
		if it.UseRandomUserAgent {
			value := gutil.RandChoice(UserAgents).(string)
			rs.Header.Set("User-Agent", value)

			glogging.Sugared.Debugf("Header.Set User-Agent=%s", value)
		}
	}

	return
}

// Do implements HttpClient
func (it *HttpRequest) Do(rq *http.Request) (rs *http.Response, err error) {
	if Client == nil {
		client := &http.Client{
			Timeout: it.Timeout,
		}

		tr := http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		if it.Proxy != "" {
			proxyNode := it.Proxy
			if !strings.HasPrefix(proxyNode, "http") {
				proxyNode = fmt.Sprintf("http://%s", proxyNode)
			}
			glogging.Sugared.Debugf("proxy=%s", proxyNode)

			u, errUrl := url.Parse(proxyNode)
			if errUrl != nil {
				err = errUrl
				return
			}
			tr.Proxy = http.ProxyURL(u)
			client.Transport = &tr
		}
		client.Transport = &tr

		Client = client
	}

	now := time.Now()
	it.Rqts = now.UnixNano() / 1000000
	rs, err = Client.Do(rq)
	elapsed := time.Since(now)

	glogging.Sugared.Debugf("%s %s elapsed=%v err=%v", it.Method, it.Uri, ghuman.FmtDuration(elapsed), err)

	if err != nil {
		it.Elapsed = elapsed
		return
	}

	return
}

// read body data from http.Response, it allows call multiple times
func ReadBody(httpRs *http.Response) (rs []byte, err error) {
	defer httpRs.Body.Close()
	rs, err = ioutil.ReadAll(httpRs.Body)
	if err == nil {
		httpRs.Body = io.NopCloser(bytes.NewReader(rs))
	}
	return rs, err
}
