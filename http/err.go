package http

import "log"

// CheckRESTErr print error with stack context and return true for error else false.
func CheckRequestErr(fullurl string, status int, respBody []byte, err error) bool {
	if err != nil {
		log.Println("[error] request", fullurl, err)
		return true
	} else {
		if int(status/100) != 2 {
			log.Printf("[error] request %s resp.status=%d resp.body -%s-", fullurl, status, string(respBody))
			return true
		}
	}

	return false
}
