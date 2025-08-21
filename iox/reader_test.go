package iox

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func ExampleIdentityReader() {
	s := "identity encoding"

	br := bytes.NewReader([]byte(s))

	er, status := NewEncodingReader(br, nil)
	fmt.Printf("test: NewEncodingReader(none) -> [er:%v] [status:%v]\n", reflect.TypeOf(er).String(), status)

	buf, err := io.ReadAll(er)
	fmt.Printf("test: Read() -> [err:%v] [content:\"%v\"]\n", err, string(buf))

	//h := make(httpx.Header)
	//h.Add(ContentEncoding, GzipEncoding)
	//er, status = NewEncodingReader(br, h)
	//fmt.Printf("test: NewEncodingReader(gzip) -> [er:%v] [status:%v]\n", reflect.TypeOf(er).String(), status)

	//Output:
	//test: NewEncodingReader(none) -> [er:*iox.identityReader] [status:<nil>]
	//test: Read() -> [err:<nil>] [content:"identity encoding"]

}

/*
func ExampleEncodingReader_Error() {
	s := address3Url
	buf0, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := strings.NewReader(string(buf0))

	h := make(httpx.Header)
	h.Set(ContentEncoding, BrotliEncoding)
	reader, status := EncodingReader(r, h)
	fmt.Printf("test: EncodingReader() -> [reader:%v] [status:%v]\n", reader, status)

	h.Set(ContentEncoding, DeflateEncoding)
	reader, status = EncodingReader(r, h)
	fmt.Printf("test: EncodingReader() -> [reader:%v] [status:%v]\n", reader, status)

	h.Set(ContentEncoding, CompressEncoding)
	reader, status = EncodingReader(r, h)
	fmt.Printf("test: EncodingReader() -> [reader:%v] [status:%v]\n", reader, status)

	//Output:
	//test: EncodingReader() -> [reader:<nil>] [status:Content Decoding Failure [error: content encoding not supported [br]]]
	//test: EncodingReader() -> [reader:<nil>] [status:Content Decoding Failure [error: content encoding not supported [deflate]]]
	//test: EncodingReader() -> [reader:<nil>] [status:Content Decoding Failure [error: content encoding not supported [compress]]]

}


*/

/*
func ExampleEncodingReader_Gzip() {
	s := searchResultsGzip
	buf0, err0 := os.ReadFile(FileName(s))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	r := bytes.NewReader(buf0)

	h := make(httpx.Header)
	h.Set(ContentEncoding, GzipEncoding)
	zr, _ := EncodingReader(r, h)
	buf, err := iox.ReadAll(zr)
	fmt.Printf("test: iox.ReadAll() -> [input:%v] [output:%v] [err:%v]\n", httpx.DetectContentType(buf0), httpx.DetectContentType(buf), err)

	//Output:
	//test: iox.ReadAll() -> [input:application/x-gzip] [output:text/html; charset=utf-8] [err:<nil>]

}

*/
