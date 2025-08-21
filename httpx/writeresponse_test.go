package httpx

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/appellative-ai/common/iox"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

const (
	activityJsonFile = "file://[cwd]/httpxtest/activity.json"
	activityGzipFile = "file://[cwd]/httpxtest/activity.gz"

	testResponseText = "file://[cwd]/httpxtest/test-response.txt"
	jsonContentType  = "application/json"
)

type activity struct {
	ActivityID   string `json:"ActivityID"`
	ActivityType string `json:"ActivityType"`
	Agent        string `json:"Agent"`
	AgentUri     string `json:"AgentUri"`
	Assignment   string `json:"Assignment"`
	Controller   string `json:"Controller"`
	Behavior     string `json:"Behavior"`
	Description  string `json:"Description"`
}

var (
	activityJson []byte
	activityGzip []byte
	activityList []activity
)

func init() {
	var err error
	var buf []byte

	buf, err = os.ReadFile(iox.FileName(activityJsonFile))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
		return
	}
	err = json.Unmarshal(buf, &activityList)
	if err != nil {
		fmt.Printf("test: json.Unmarshal() -> [err:%v]\n", err)
		return
	}
	activityJson, err = json.Marshal(activityList)
	if err != nil {
		fmt.Printf("test: json.Mmarshal() -> [err:%v]\n", err)
		return
	}

	activityGzip, err = os.ReadFile(iox.FileName(activityGzipFile))
	if err != nil {
		if strings.Contains(err.Error(), "open") {
			buff := new(bytes.Buffer)

			// write, flush and close
			zw := gzip.NewWriter(buff)
			cnt, err0 := zw.Write(activityJson)
			ferr := zw.Flush()
			cerr := zw.Close()
			fmt.Printf("test: gzip.Writer() -> [cnt:%v] [write-err:%v] [flush-err:%v] [close_err:%v]\n", cnt, err0, ferr, cerr)
			err = os.WriteFile(iox.FileName(activityGzipFile), buff.Bytes(), 667)
			fmt.Printf("test: os.WriteFile(\"%v\") -> [err:%v]\n", activityGzipFile, err)

		} else {
			fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
		}
	}
}

func ExampleWriteResponse_StatusHeaders() {
	// all nil
	rec := httptest.NewRecorder()
	WriteResponse(rec, nil, 0, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,0,nil) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	// status code
	rec = httptest.NewRecorder()
	WriteResponse(rec, nil, http.StatusTeapot, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,StatusTeapot,nil) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	// status code, headers list
	rec = httptest.NewRecorder()
	WriteResponse(rec, []Attr{{Key: ContentType, Value: ContentTypeTextHtml}}, http.StatusOK, nil, CreateAcceptEncodingHeader())
	fmt.Printf("test: WriteResponse(w,list,StatusOK,nil) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	// status code, httpx.Header
	rec = httptest.NewRecorder()
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)
	h.Add(ContentEncoding, ContentEncodingGzip)
	WriteResponse(rec, h, http.StatusGatewayTimeout, nil, nil)
	fmt.Printf("test: WriteResponse(w,httpx.Header,StatusGatewayTimeout,nil) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	//Output:
	//test: WriteResponse(w,nil,0,nil) -> [status-code:200] [header:map[]]
	//test: WriteResponse(w,nil,StatusTeapot,nil) -> [status-code:418] [header:map[]]
	//test: WriteResponse(w,list,StatusOK,nil) -> [status-code:200] [header:map[Content-Type:[text/html]]]
	//test: WriteResponse(w,httpx.Header,StatusGatewayTimeout,nil) -> [status-code:504] [header:map[Content-Encoding:[gzip] Content-Type:[application/json]]]

}

func ExampleWriteResponse_JSON() {
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)

	// JSON activity list
	rec := httptest.NewRecorder()
	WriteResponse(rec, h, 0, activityList, nil)
	buf, status0 := readAll(rec.Result().Body)
	fmt.Printf("test: WriteResponse(w,httpx.Header,OK,[]activity) -> [read-all:%v] [in:%v] [out:%v]\n", status0, len(activityJson), len(buf))

	// JSON reader
	rec = httptest.NewRecorder()
	reader := bytes.NewReader(activityJson)
	WriteResponse(rec, h, 0, reader, nil)
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: WriteResponse(w,httpx.Header,OK,io.Reader) -> [read-all:%v] [in:%v] [out:%v]\n", status0, len(activityJson), len(buf))

	//Output:
	//test: WriteResponse(w,httpx.Header,OK,[]activity) -> [read-all:<nil>] [in:395] [out:395]
	//test: WriteResponse(w,httpx.Header,OK,io.Reader) -> [read-all:<nil>] [in:395] [out:395]

}

func ExampleWriteResponse_Encoding() {
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)
	//h.Add(AcceptEncoding, AcceptEncodingValue)

	// Should encode
	rec := httptest.NewRecorder()
	WriteResponse(rec, h, 0, activityList, CreateAcceptEncodingHeader())
	buf, status0 := readAll(rec.Result().Body)
	fmt.Printf("test: WriteResponse(w,httpx.Header,0,[]activity) -> [read-all:%v] [buf:%v][header:%v]\n", status0, http.DetectContentType(buf), rec.Result().Header)

	// Should not encode as a ContentEncoding header exists
	h = make(http.Header)
	h.Add(ContentType, ContentTypeJson)
	//h.Add(AcceptEncoding, AcceptEncodingValue)
	//h.Add(ContentEncoding, NoneEncoding)
	rec = httptest.NewRecorder()
	WriteResponse(rec, h, 0, activityList, CreateAcceptEncodingHeader())
	buf, status0 = readAll(rec.Result().Body)
	fmt.Printf("test: WriteResponse(w,httpx.Header,0,[]activity) -> [read-all:%v] [buf:%v][header:%v]\n", status0, http.DetectContentType(buf), rec.Result().Header)

	//Output:
	//test: WriteResponse(w,httpx.Header,0,[]activity) -> [read-all:<nil>] [buf:application/x-gzip][header:map[Content-Encoding:[gzip] Content-Type:[application/json]]]
	//test: WriteResponse(w,httpx.Header,0,[]activity) -> [read-all:<nil>] [buf:text/plain; charset=utf-8][header:map[Content-Encoding:[none] Content-Type:[application/json]]]

}
