package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Request perform HTTP request with custom method, timeout and headers parameters.
func Request(
	timeout time.Duration,
	method string,
	fullurl string,
	headers *map[string]interface{},
	reqBody []byte,
	proxyUrl string,
) (status int, bodyResp []byte, err error) {

	client := &http.Client{
		Timeout: timeout,
	}

	if proxyUrl != "" {
		var tr *http.Transport
		u, errUrl := url.Parse(proxyUrl)
		if errUrl != nil {
			err = errUrl
			return
		}
		tr = &http.Transport{
			Proxy:           http.ProxyURL(u),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client.Transport = tr
		//log.Printf("[debug] proxy=%s", proxyUrl)
	} else {
		var tr *http.Transport
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client.Transport = tr
	}

	start := time.Now()

	var body io.Reader
	if reqBody != nil && len(reqBody) > 0 {
		body = bytes.NewBuffer(reqBody)
	}
	req, err := http.NewRequest(method, fullurl, body)
	if err != nil {
		//log.Println(fmt.Sprintf("[error] http.NewRequest failed, %s %s reqbody.bytes=%d", method, fullurl, len(reqBody)), err)
		return
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	resp, err := client.Do(req)
	elapsed := time.Since(start) / time.Millisecond

	if err != nil {
		//log.Println(fmt.Sprintf("[error] client.Do failed, %s %s reqbody.bytes=%d", method, fullurl, len(reqBody)), err)
		return
	}

	defer resp.Body.Close()

	status = resp.StatusCode

	bodyResp, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Println(fmt.Sprintf("[error] read body failed, url=%s reqbody.bytes=%d", fullurl, len(reqBody)), err)
		return
	}

	warnIfSlowThan := time.Duration(250) * time.Millisecond
	if elapsed > warnIfSlowThan {
		//log.Println(fmt.Sprintf("[debug] request %s %d %dms", fullurl, status, elapsed))
	}
	return
}
