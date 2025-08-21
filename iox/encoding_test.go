package iox

import (
	"fmt"
	"net/http"
	"os"
)

const (
	contentType = "Content-Type"
)

type newAddress struct {
	City    string
	State   string
	ZipCode string
}

func ExampleDecode_TextPlain() {
	buf, status := Decode(nil, nil)
	fmt.Printf("test: Decode(nil,nil) -> [buf:%v] [status:%v]\n", len(buf), status)

	content, err0 := os.ReadFile(FileName(testResponseTxt))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buf, status = Decode(content, nil)
	fmt.Printf("test: Decode(content,nil) -> [buf:%v] [status:%v] [content-type:%v] [buf-type:%v]\n", len(buf), status, http.DetectContentType(content), http.DetectContentType(buf))

	//Output:
	//test: Decode(nil,nil) -> [buf:0] [status:<nil>]
	//test: Decode(content,nil) -> [buf:188] [status:<nil>] [content-type:text/plain; charset=utf-8] [buf-type:text/plain; charset=utf-8]

}

func ExampleDecode_Gzip() {
	content, err0 := os.ReadFile(FileName(testResponseGzip))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buf, status := Decode(content, nil)
	fmt.Printf("test: Decode(content,nil) -> [buf:%v] [status:%v] [content-type:%v] [buf-type:%v]\n", len(buf), status, http.DetectContentType(content), http.DetectContentType(buf))

	h := make(http.Header)
	h.Set(ContentEncoding, GzipEncoding)
	buf, status = Decode(content, h)
	fmt.Printf("test: Decode(content,h) -> [buf:%v] [status:%v] [content-type:%v] [buf-type:%v]\n", len(buf), status, http.DetectContentType(content), http.DetectContentType(buf))

	//Output:
	//test: Decode(content,nil) -> [buf:188] [status:<nil>] [content-type:application/x-gzip] [buf-type:text/plain; charset=utf-8]
	//test: Decode(content,h) -> [buf:188] [status:<nil>] [content-type:application/x-gzip] [buf-type:text/plain; charset=utf-8]

}

func ExampleDecode_Error() {
	h := make(http.Header)
	h.Set(ContentEncoding, DeflateEncoding)

	content, err0 := os.ReadFile(FileName(testResponseTxt))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buf, status := Decode(content, h)
	fmt.Printf("test: Decode(content,h) -> [buf:%v] [status:%v] [content-type:%v] [buf-type:%v]\n", len(buf), status, http.DetectContentType(content), http.DetectContentType(buf))

	content, err0 = os.ReadFile(FileName(testResponseGzip))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buf, status = Decode(content, h)
	fmt.Printf("test: Decode(content,h) -> [buf:%v] [status:%v] [content-type:%v] [buf-type:%v]\n", len(buf), status, http.DetectContentType(content), http.DetectContentType(buf))

	//Output:
	//test: Decode(content,h) -> [buf:188] [status:error: content encoding not supported [deflate]] [content-type:text/plain; charset=utf-8] [buf-type:text/plain; charset=utf-8]
	//test: Decode(content,h) -> [buf:188] [status:error: content encoding not supported [deflate]] [content-type:application/x-gzip] [buf-type:application/x-gzip]

}

func ExampleZipFile() {
	status := ZipFile(helloWorldTxt)

	fmt.Printf("test: ZipFile(\"\") -> [status:%v]\n", status)

	//Output:
	//test: ZipFile("") -> [status:<nil>]

}
