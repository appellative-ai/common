package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const (
	ContentTypeText     = "text/plain charset=utf-8"
	ContentTypeBinary   = "application/octet-stream"
	ContentTypeJson     = "application/json"
	ContentTypeTextHtml = "text/html"
)

// Content -
type Content struct {
	Fragment string // returned on a Get
	Type     string // Content-Type
	Value    any
}

func (c Content) String() string {
	return fmt.Sprintf("fragment: %v type: %v value: %v", c.Fragment, c.Type, c.Value != nil)
}

func (c Content) Valid(contentType string) bool {
	return c.Value != nil && c.Type == contentType
}

func New[T any](ct *Content) (t T, err error) {
	if ct == nil {
		return t, errors.New(fmt.Sprintf("content is nil"))
	}
	if ct.Type == "" || ct.Value == nil {
		return t, errors.New(fmt.Sprintf("content type is empty, or content value is nil"))
	}
	// Check for binary and unmarshal
	if _, ok := ct.Value.([]byte); ok {
		return Unmarshal[T](ct)
	}
	var ok bool

	if t, ok = ct.Value.(T); ok {
		return t, nil
	}
	return t, errors.New(fmt.Sprintf("content value type: %v is not of generic type: %v", reflect.TypeOf(ct.Value), reflect.TypeOf(t)))
}

// Unmarshal - []byte -> string, []byte, io.Reader, type via json.Unmarshal
func Unmarshal[T any](ct *Content) (t T, err error) {
	var body []byte
	var ok bool

	if ct == nil {
		return t, errors.New(fmt.Sprintf("content is nil"))
	}
	if ct.Type == "" || ct.Value == nil {
		return t, errors.New(fmt.Sprintf("content type is empty, or content value is nil"))
	}
	if body, ok = ct.Value.([]byte); !ok {
		return t, errors.New(fmt.Sprintf("content value type: %v is not of type: []byte", reflect.TypeOf(ct.Value)))
	}
	if len(body) == 0 {
		return t, nil
	}
	switch ptr := any(&t).(type) {
	case *string:
		if ct.Type != ContentTypeText && ct.Type != ContentTypeTextHtml {
			return t, errors.New(fmt.Sprintf("content type: %v is invalid for string", ct.Type))
		}
		*ptr = string(body)
	case *[]byte:
		if ct.Type != ContentTypeBinary {
			return t, errors.New(fmt.Sprintf("content type: %v is invalid for []byte", ct.Type))
		}
		*ptr = body
	default:
		if ct.Type != ContentTypeJson {
			return t, errors.New(fmt.Sprintf("content type: %v is invalid for json.Unmarshal()", ct.Type))
		}
		err := json.Unmarshal(body, ptr)
		if err != nil {
			return t, errors.New(fmt.Sprintf("JSON unmarshalling %v", err))
		}
	}
	return t, nil
}

// Marshal -  type -> []byte | io.Reader
func Marshal[T any](ct *Content) (t T, err error) {
	var buf []byte

	if ct == nil {
		return t, errors.New(fmt.Sprintf("content is nil"))
	}
	if ct.Type == "" || ct.Value == nil {
		return t, errors.New(fmt.Sprintf("error: content type is empty, or content value is nil"))
	}
	switch ptr := ct.Value.(type) {
	case string:
		buf = []byte(ptr)
	case []byte:
		buf = ptr
	default:
		var err1 error
		buf, err1 = json.Marshal(ptr)
		if err1 != nil {
			return t, err1
		}
	}
	if len(buf) == 0 {
		return t, errors.New("content value is empty")
	}
	switch ptr := any(&t).(type) {
	case *[]byte:
		*ptr = buf
		return t, nil
	case *io.Reader:
		*ptr = bytes.NewReader(buf)
		return t, nil
	default:
	}
	return t, errors.New(fmt.Sprintf("error: generic type: %v is not supported for marshalling", reflect.TypeOf(t)))
}
