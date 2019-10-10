package gstr

// Shorten it cuts and concatenates a part of string with ellipsis if its length long than 500 bytes.
func Shorten(b []byte) string {
	var bodyChunk string
	maxBodyChunk := 500
	bodyStr := string(b)
	if len(bodyStr) > maxBodyChunk {
		bodyChunk = bodyStr[:maxBodyChunk] + "..."
	} else {
		bodyChunk = bodyStr
	}
	return bodyChunk
}
