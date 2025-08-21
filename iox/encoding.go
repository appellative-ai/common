package iox

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	AcceptEncoding      = "Accept-Encoding"
	AcceptEncodingValue = "gzip, deflate, br"
	ContentEncoding     = "Content-Encoding"

	GzipEncoding     = "gzip"
	BrotliEncoding   = "br"
	DeflateEncoding  = "deflate"
	CompressEncoding = "compress"
	NoneEncoding     = "none"

	ApplicationGzip    = "application/x-gzip"
	ApplicationBrotli  = "application/x-br"
	ApplicationDeflate = "application/x-deflate"

	encodingErrorFmt = "error: content encoding not supported [%v]"
)

func acceptEncoding(h http.Header) string {
	if h == nil {
		return NoneEncoding
	}
	enc := h.Get(AcceptEncoding)
	if len(enc) > 0 {
		return strings.ToLower(enc)
	}
	return NoneEncoding
}

func contentEncoding(h http.Header) string {
	if h == nil {
		return NoneEncoding
	}
	enc := h.Get(ContentEncoding)
	if len(enc) > 0 {
		return strings.ToLower(enc)
	}
	return NoneEncoding
}

func newStatusContentEncodingError(ct string) error {
	return errors.New(fmt.Sprintf(encodingErrorFmt, ct))

}

// Decode - decode a []byte
func Decode(buf []byte, h http.Header) ([]byte, error) {
	if len(buf) == 0 {
		return buf, nil
	}
	ct := NoneEncoding
	if h == nil {
		ct = http.DetectContentType(buf)
	} else {
		ct = contentEncoding(h)
	}
	switch ct {
	case ApplicationGzip, GzipEncoding:
		zr, status := NewGzipReader(bytes.NewReader(buf))
		if status != nil {
			return nil, status
		}
		buf2, err1 := io.ReadAll(zr)
		err2 := zr.Close()
		if err1 != nil {
			return nil, err1
		}
		if err2 != nil {
			return nil, err2
		}
		return buf2, nil
	case ApplicationBrotli, BrotliEncoding:
		return buf, newStatusContentEncodingError(ct)
	case ApplicationDeflate, DeflateEncoding:
		return buf, newStatusContentEncodingError(ct)
	default:
		return buf, nil
	}
}

func ZipFile(uri string) error {
	if len(uri) == 0 {
		return errors.New("error: file path is empty")
	}
	path := FileName(uri)
	content, err0 := os.ReadFile(path)
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return err0
	}
	// write content
	buff := new(bytes.Buffer)
	zw := NewGzipWriter(buff)
	cnt, err := zw.Write(content)
	err1 := zw.Close()
	if err != nil {
		return err
	}
	if cnt == 0 || err1 != nil {
		fmt.Printf("error: count %v err %v", cnt, err1)
	}
	i := strings.LastIndex(path, ".")
	path2 := ""
	if i > 0 {
		path2 = path[:i]
		path2 += ".gz"
	} else {
		path2 = path + ".gz"
	}
	err = os.WriteFile(path2, buff.Bytes(), 667)
	if err != nil {
		return err
	}
	return nil
}
