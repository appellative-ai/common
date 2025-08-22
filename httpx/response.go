package httpx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/appellative-ai/common/iox"
	"io"
	"net/http"
	"reflect"
)

const (
	fileExistsError = "The system cannot find the file specified"
)

// TransformBody - read the body and create a new []byte buffer reader
func TransformBody(resp *http.Response) error {
	if resp == nil || resp.Body == nil {
		return nil
	}
	var ce string
	if resp.Header != nil {
		ce = resp.Header.Get(ContentEncoding)
	}
	if ce == "" || ce == iox.NoneEncoding {
		buf, err := readAll(resp.Body)
		if err == nil {
			resp.ContentLength = int64(len(buf))
			resp.Body = io.NopCloser(bytes.NewReader(buf))
		}
		return err
	}
	r, err := iox.NewEncodingReader(resp.Body, resp.Header)
	if err != nil {
		return err
	}
	var buf []byte
	cnt, err2 := r.Read(buf)
	if err2 != nil {
		return err2
	}
	resp.Header.Del(ContentEncoding)
	resp.ContentLength = int64(cnt)
	resp.Body = io.NopCloser(bytes.NewReader(buf))
	return nil

}

func NewResponse(statusCode int, h http.Header, content any) (resp *http.Response) {
	resp = &http.Response{StatusCode: statusCode, ContentLength: -1, Header: h, Body: EmptyReader}
	if h == nil {
		resp.Header = make(http.Header)
	}
	if content == nil {
		return resp
	}
	switch ptr := (content).(type) {
	case error:
		if ptr.Error() != "" {
			resp.ContentLength = int64(len(ptr.Error()))
			resp.Body = io.NopCloser(bytes.NewReader([]byte(ptr.Error())))
		}
	case string:
		if ptr != "" {
			resp.ContentLength = int64(len(ptr))
			resp.Body = io.NopCloser(bytes.NewReader([]byte(ptr)))
		}
	case []byte:
		if len(ptr) > 0 {
			resp.ContentLength = int64(len(ptr))
			resp.Body = io.NopCloser(bytes.NewReader(ptr))
		}
	case bytes.Buffer:
		buf := ptr.Bytes()
		resp.ContentLength = int64(len(buf))
		resp.Body = io.NopCloser(bytes.NewReader(buf))
	default:
		err := errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr)))
		return &http.Response{StatusCode: http.StatusInternalServerError, Header: SetHeader(nil, ContentType, ContentTypeText), ContentLength: int64(len(err.Error())), Body: io.NopCloser(bytes.NewReader([]byte(err.Error())))}
	}
	return resp
}

/*
// NewResponseFromUri - read a Http response given a URL
func NewResponseFromUri(uri any) (*http.Response, error) {
	serverErr := &http.Response{StatusCode: http.StatusInternalServerError, Status: internalError, Header: make(http.Header)}
	if uri == nil {
		return serverErr, errors.New("error: URL is nil")
	}
	buf, err := readFile(uri)
	if err != nil {
		if strings.Contains(err.Error(), fileExistsError) {
			return &http.Response{StatusCode: http.StatusNotFound, Status: "Not Found", Header: make(http.Header)}, err
		}
		return serverErr, err
	}
	resp1, err2 := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
	if err2 != nil {
		return serverErr, err2
	}
	return resp1, nil

}


*/
