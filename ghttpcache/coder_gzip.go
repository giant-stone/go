package ghttpcache

import (
	"bytes"
	"compress/gzip"
	"io"
)

type coderGzip struct {
}

func NewCoderGzip() (rs *coderGzip) {
	return &coderGzip{}
}

func (its *coderGzip) Compress(data []byte) (rs []byte, err error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)

	_, err = w.Write(data)
	if err != nil {
		return
	}

	if err = w.Close(); err != nil {
		return
	}

	rs = buf.Bytes()
	return
}

func (its *coderGzip) Decompress(data []byte) (rs []byte, err error) {
	b := bytes.NewBuffer(data)

	var r io.Reader
	r, err = gzip.NewReader(b)
	if err != nil {
		return
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	rs = resB.Bytes()

	return
}
