package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	jsonToken = "json"
)

func writeContent(w io.Writer, content any, contentType string) (length int64, err error) {
	var cnt int

	if content == nil {
		return 0, nil
	}
	switch ptr := (content).(type) {
	case []byte:
		cnt, err = w.Write(ptr)
	case string:
		cnt, err = w.Write([]byte(ptr))
	case error:
		cnt, err = w.Write([]byte(ptr.Error()))
	case io.Reader:
		var buf []byte
		var err1 error

		buf, err1 = readAll(ptr)
		if err1 != nil {
			return 0, err1
		}
		cnt, err = w.Write(buf)
	case io.ReadCloser:
		var buf []byte
		var err1 error

		buf, err1 = readAll(ptr)
		_ = ptr.Close()
		if err1 != nil {
			return 0, err1
		}
		cnt, err = w.Write(buf)
	default:
		if strings.Contains(contentType, jsonToken) {
			var buf []byte

			buf, err = json.Marshal(content)
			if err != nil {
				return
			}
			cnt, err = w.Write(buf)
		} else {
			return 0, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr)))
		}
	}
	return int64(cnt), err
}
