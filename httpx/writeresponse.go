package httpx

import (
	"github.com/appellative-ai/common/iox"
	"net/http"
)

// WriteResponse - write a httpx.Response, utilizing the content, status code, and headers
// Content types supported: []byte, string, error, io.Reader, io.ReadCloser. Other types will be treated as JSON and serialized, if
// the headers content type is JSON. If not JSON, then an error will be raised.
func WriteResponse(w http.ResponseWriter, headers any, statusCode int, content any, reqHeader http.Header) (contentLength int64) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	SetHeaders(w, headers)
	if content == nil {
		w.WriteHeader(statusCode)
		return 0
	}
	clone := CloneHeader(reqHeader)
	if len(w.Header().Get(ContentEncoding)) != 0 {
		clone.Del(AcceptEncoding)
	}
	writer, err := iox.NewEncodingWriter(w, clone)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return 0
	}
	if writer.ContentEncoding() != NoneEncoding {
		w.Header().Add(ContentEncoding, writer.ContentEncoding())
	}
	w.WriteHeader(statusCode)
	contentLength, err = writeContent(writer, content, w.Header().Get(ContentType))
	_ = writer.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	return contentLength
}

func CreateAcceptEncodingHeader() http.Header {
	out := make(http.Header)
	out.Add(AcceptEncoding, AcceptEncodingValue)
	return out
}
