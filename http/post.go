package http

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Post perform HTTP POST request with timeout and headers parameters.
func Post(u string, timeout time.Duration, headers *map[string]string, reqBody []byte) (body []byte, status int, err error) {
	client := &http.Client{
		Timeout: timeout,
	}
	start := time.Now()

	req, err := http.NewRequest("POST", u, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(fmt.Sprintf("[error] http.NewRequest failed, POST %s reqbody.bytes=%d", u, len(reqBody)), err)
		return
	}

	if headers != nil {
		for k, v := range *headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	elapsed := time.Since(start) / time.Millisecond

	if err != nil {
		log.Println(fmt.Sprintf("[error] client.Do failed, POST %s reqbody.bytes=%d", u, len(reqBody)), err)
		return
	}

	defer resp.Body.Close()

	status = resp.StatusCode

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("[error] read body failed, url=%s reqbody.bytes=%d", u, len(reqBody)), err)
		return
	}

	warnIfSlowThan := time.Duration(250) * time.Millisecond
	if elapsed > warnIfSlowThan {
		log.Println(fmt.Sprintf("[debug] request %s %d %dms", u, status, elapsed))
	}
	return
}
