package httpx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/appellative-ai/common/iox"
	"io"
	"net/http/httptest"
	"os"
)

func ExampleWriteContent_Buffer() {
	content := "<h1>Hello, World!</h1>"
	ct := ""

	// nil
	rec := httptest.NewRecorder()
	cnt, status := writeContent(rec, nil, ct)
	buf, status0 := readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(nil) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	// []byte
	rec = httptest.NewRecorder()
	cnt, status = writeContent(rec, []byte(content), ct)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: writeContent([]byte) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	// empty string
	rec = httptest.NewRecorder()
	cnt, status = writeContent(rec, "", ct)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(\"\") -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	// string
	rec = httptest.NewRecorder()
	cnt, status = writeContent(rec, content, ct)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(string) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	// error message
	rec = httptest.NewRecorder()
	cnt, status = writeContent(rec, errors.New("This is example error message text"), ct)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(error) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	//Output:
	//test: writeContent(nil) -> [cnt:0] [write-status:<nil>] [body:] [read-status:<nil>]
	//test: writeContent([]byte) -> [cnt:22] [write-status:<nil>] [body:<h1>Hello, World!</h1>] [read-status:<nil>]
	//test: writeContent("") -> [cnt:0] [write-status:<nil>] [body:] [read-status:<nil>]
	//test: writeContent(string) -> [cnt:22] [write-status:<nil>] [body:<h1>Hello, World!</h1>] [read-status:<nil>]
	//test: writeContent(error) -> [cnt:34] [write-status:<nil>] [body:This is example error message text] [read-status:<nil>]

}

func ExampleWriteContent_Reader() {
	content, err0 := os.ReadFile(iox.FileName(testResponseText))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	ct := ""

	// io.Reader
	rec := httptest.NewRecorder()
	reader := bytes.NewReader(content)
	cnt, status := writeContent(rec, reader, ct)
	buf, status0 := readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(io.Reader) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, len(buf), status0)

	// io.ReadCloser
	rec = httptest.NewRecorder()
	reader = bytes.NewReader(content)
	cnt, status = writeContent(rec, io.NopCloser(reader), ct)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(io.ReadCloser) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, len(buf), status0)

	//Output:
	//test: writeContent(io.Reader) -> [cnt:188] [write-status:<nil>] [body:188] [read-status:<nil>]
	//test: writeContent(io.ReadCloser) -> [cnt:188] [write-status:<nil>] [body:188] [read-status:<nil>]

}

func ExampleWriteContent_Json() {
	ct := ""
	content := activity{
		ActivityID:   "123456",
		ActivityType: "action",
		Agent:        "Controller",
		AgentUri:     "https://somehost.com/id",
		Assignment:   "case #",
		Controller:   "egress",
		Behavior:     "timeout",
		Description:  "decreased timeout",
	}

	// error - invalid type, no content type
	rec := httptest.NewRecorder()
	cnt, status := writeContent(rec, content, ct)
	buf, status0 := readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(httpx.testActivity) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	// JSON
	rec = httptest.NewRecorder()
	cnt, status = writeContent(rec, content, jsonContentType)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: writeContent(httpx.testActivity) -> [cnt:%v] [write-status:%v] [body:%v] [read-status:%v]\n", cnt, status, string(buf), status0)

	//Output:
	//test: writeContent(httpx.testActivity) -> [cnt:0] [write-status:error: content type is invalid: httpx.activity] [body:] [read-status:<nil>]
	//test: writeContent(httpx.testActivity) -> [cnt:204] [write-status:<nil>] [body:{"ActivityID":"123456","ActivityType":"action","Agent":"Controller","AgentUri":"https://somehost.com/id","Assignment":"case #","Controller":"egress","Behavior":"timeout","Description":"decreased timeout"}] [read-status:<nil>]

}
