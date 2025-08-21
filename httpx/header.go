package httpx

import (
	"net/http"
	"strings"
)

const (
	ContentTypeJson = "application/json"
	ContentType     = "Content-Type"
	ContentTypeText = "text/plain charset=utf-8"

	ContentEncoding     = "Content-Encoding"
	AcceptEncoding      = "Accept-Encoding"
	ContentEncodingGzip = "gzip"

	ContentTypeTextHtml = "text/html"

	AcceptEncodingValue = "gzip, deflate, br"
	GzipEncoding        = "gzip"
	NoneEncoding        = "none"
)

func SetHeader(h http.Header, name, value string) http.Header {
	if h == nil {
		h = make(http.Header)
	}
	h.Set(name, value)
	return h
}

func SetHeaders(w http.ResponseWriter, headers any) {
	if headers == nil {
		return
	}
	if pairs, ok := headers.([]Attr); ok {
		for _, pair := range pairs {
			w.Header().Set(strings.ToLower(pair.Key), pair.Value)
		}
		return
	}
	if h, ok := headers.(http.Header); ok {
		for k, v := range h {
			if len(v) > 0 {
				w.Header().Set(strings.ToLower(k), v[0])
			}
		}
	}
}

func CloneHeader(hdr http.Header) http.Header {
	clone := hdr.Clone()
	if clone == nil {
		clone = make(http.Header)
	}
	return clone
}
