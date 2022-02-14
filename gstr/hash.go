package gstr

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// Md5 hashes using md5 algorithm
func Md5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Sha1 hashes a string in sha1 algorithm.
func Sha1(s string) (rs string) {
	algorithm := sha1.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}
