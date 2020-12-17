package ghttp

import (
	"io"
	"mime/multipart"
)

func AppendMultipartFormData(w *multipart.Writer, fieldName, fileName string, fileData []byte) (err error) {
	var fw io.Writer
	if fw, err = w.CreateFormFile(fieldName, fileName); err != nil {
		return
	}

	if _, err = fw.Write(fileData); err != nil {
		return
	}
	return
}
