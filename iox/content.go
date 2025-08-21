package iox

import (
	"bytes"
	"net/http"
	"strings"
)

func EncodeContent(h http.Header, content []byte) (buff []byte, encoding string, err error) {
	enc := acceptEncoding(h)
	if !strings.Contains(enc, GzipEncoding) {
		return nil, "", nil
	}
	buf := new(bytes.Buffer)
	w, err1 := NewEncodingWriter(buf, h)
	if err1 != nil {
		return nil, "", err
	}
	cnt, err2 := w.Write(content)
	err3 := w.Close()
	if err2 != nil || cnt != len(content) {
		return nil, "", err2
	}
	if err3 != nil {
		return nil, "", err3
	}
	return buf.Bytes(), GzipEncoding, nil
}
