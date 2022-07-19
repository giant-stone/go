package gstr

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

// Md5 hashes using md5 algorithm.
func Md5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Md5FromBytes hashes a slice of bytes in md5 algorithm.
func Md5FromBytes(data []byte) string {
	algorithm := md5.New()
	algorithm.Write(data)
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Sha1 hashes a string in sha1 algorithm.
func Sha1(s string) (rs string) {
	algorithm := sha1.New()
	algorithm.Write([]byte(s))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// Sha1FromBytes hashes a slice of bytes in sha1 algorithm.
func Sha1FromBytes(data []byte) (rs string) {
	algorithm := sha1.New()
	algorithm.Write(data)
	return hex.EncodeToString(algorithm.Sum(nil))
}
