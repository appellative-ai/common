package iox

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type echo struct {
	Method string      `json:"method"`
	Host   string      `json:"host"`
	Url    string      `json:"url"`
	Header http.Header `json:"header"`
}

func ExampleEncodeContent() {
	h := make(http.Header)
	h.Add("x-request-id", "1234-56-7890")
	e := echo{
		Method: http.MethodGet,
		Host:   "localhost",
		Url:    "https://www.google.com/search?q=golang",
		Header: h,
	}
	content, err1 := json.Marshal(e)
	if err1 != nil {
		fmt.Printf("test: json.Marshal() -> [err:%v]\n", err1)
	}
	//buf, encoding, err := EncodeContent(nil, content)
	//fmt.Printf("test: EncodeContent-Nil-Header() -> [buf:%v] [encoding:%v] [err:%v]\n", len(buf), encoding, err)

	h2 := make(http.Header)
	buf, encoding, err := EncodeContent(h2, content)
	fmt.Printf("test: EncodeContent-No-Accept-Encoding() -> [buf:%v] [encoding:%v] [err:%v]\n", len(buf), encoding, err)

	h2.Add(AcceptEncoding, AcceptEncodingValue)
	buf, encoding, err = EncodeContent(h2, content)
	ct := http.DetectContentType(buf)
	fmt.Printf("test: EncodeContent() -> [buf:%v] [encoding:%v] [content-type:%v] [err:%v]\n", len(buf), encoding, ct, err)

	//Output:
	//test: EncodeContent-No-Accept-Encoding() -> [buf:0] [encoding:] [err:<nil>]
	//test: EncodeContent() -> [buf:144] [encoding:gzip] [content-type:application/x-gzip] [err:<nil>]

}
