package gstr

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 hashes using md5 algorithm
func Md5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
