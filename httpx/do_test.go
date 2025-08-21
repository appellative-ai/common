package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	testContent = "this is response write content"
	requestId   = "123-request-id"
	relatesTo   = "test-relates-to"
	statusCode  = http.StatusAccepted
	XRelatesTo  = "X-Relates-To"
	XRequestId  = "X-Request-ID"
)

func deadlineExceededError(t any) bool {
	if t == nil {
		return false
	}
	if r, ok := t.(*http.Request); ok {
		return r.Context() != nil && r.Context().Err() == context.DeadlineExceeded
	}
	if e, ok := t.(error); ok {
		return e == context.DeadlineExceeded
	}
	return false
}

// dial tcp [2607:f8b0:4023:1009::68]:443: i/o timeout]

func _ExampleDo_InvalidArgument() {
	_, s := Do(nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [invalid argument : request is nil]

}

/*
func ExampleDo_ConnectivityError() {
	req, _ := httpx.NewRequest(httpx.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, status := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [errs:%v] [content-type:%v] [body:%v]\n",
		resp != nil, status.Code(), status.Errors(), resp.Header.Get("content-type"), resp.Body != nil)

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:503] [errs:[]] [content-type:text/html] [body:true]

}


*/

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(XRequestId, requestId)
	w.Header().Add(XRelatesTo, relatesTo)
	w.WriteHeader(statusCode)
	w.Write([]byte(testContent))
}

/*
func _ExampleDo_Proxy() {
	uri := "http://localhost:8080/github.com/appellative-ai/core/exchange:Do"
	req, _ := httpx.NewRequest("", uri, nil)

	resp, status := Do(req)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp != nil, status)

	status = RegisterHandler(uri, testHandler)
	fmt.Printf("test: RegisterEndpoint() -> [status:%v]\n", status)

	resp, status = Do(req)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp != nil, status)

	fmt.Printf("test: Do() -> [write-requestId:%v] [response-requestId:%v]\n", requestId, resp.Header.Get(aspect.XRequestId))
	fmt.Printf("test: Do() -> [write-relatesTo:%v] [response-relatesTo:%v]\n", relatesTo, resp.Header.Get(aspect.XRelatesTo))
	fmt.Printf("test: Do() -> [write-statusCode:%v] [response-statusCode:%v]\n", statusCode, resp.StatusCode)

	buf, _ := readAll(resp.Body)
	fmt.Printf("test: Do() -> [write-content:%v] [response-content:%v]\n", testContent, string(buf))

	//Output:
	//test: Do() -> [resp:true] [status:Internal Error [Get "http://localhost:8080/github.com/advanced-go/core/exchange:Do": dial tcp [::1]:8080: connectex: No connection could be made because the target machine actively refused it.]]
	//test: RegisterEndpoint() -> [status:OK]
	//test: Do() -> [resp:true] [status:Accepted]
	//test: Do() -> [write-requestId:123-request-id] [response-requestId:123-request-id]
	//test: Do() -> [write-relatesTo:test-relates-to] [response-relatesTo:test-relates-to]
	//test: Do() -> [write-statusCode:202] [response-statusCode:202]
	//test: Do() -> [write-content:this is response write content] [response-content:this is response write content]

}


*/

func defaultDo(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		//if deadlineExceededError(r) {
		//	return &http.Response{StatusCode: http.StatusGatewayTimeout}, err
		//}
		return resp, err //&http.Response{StatusCode: http.StatusInternalServerError}, err
	}
	//time.Sleep(time.Second * 2)
	buf, err1 := readAll(resp.Body)
	if err1 != nil {
		if deadlineExceededError(err1) {
			return &http.Response{StatusCode: http.StatusGatewayTimeout}, err1
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, err1
	}
	if buf != nil {
	}
	return resp, nil
}

func ExampleDefaultDo_OK() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := defaultDo(req)
	fmt.Printf("test: DefaultDo_OK()-no-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*4)
	defer cancel()
	r2 := req.Clone(ctx)
	resp, err = defaultDo(r2)
	fmt.Printf("test: DefaultDo_OK()-5s-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: DefaultDo_OK()-no-timeout -> [status-code:200] [err:<nil>]
	//test: DefaultDo_OK()-5s-timeout -> [status-code:200] [err:<nil>]

}

func ExampleDefaultDo_Timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := defaultDo(req)
	if err != nil {
		fmt.Printf("test: DefaultDo_Timeout()-Get()-timeout -> [resp:%v] [err:%v]\n", resp, err)
	} else {
		fmt.Printf("test: DefaultDo_Timeout()-Get()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel2()
	req, _ = http.NewRequestWithContext(ctx2, http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, err = defaultDo(req)
	fmt.Printf("test: DefaultDo_Timeout()-ReadAll()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: DefaultDo_Timeout()-Get()-timeout -> [resp:<nil>] [err:Get "https://www.google.com/search?q=golang": context deadline exceeded]
	//test: DefaultDo_Timeout()-ReadAll()-timeout -> [status-code:200] [err:<nil>]

}

func exchangeDo(r *http.Request) (*http.Response, error) {
	resp, err := Do(r)
	if err != nil {
		return resp, err
	}
	buf, err1 := readAll(resp.Body)
	if err1 != nil {
		if deadlineExceededError(err1) {
			return &http.Response{StatusCode: http.StatusGatewayTimeout}, err1
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, err1
	}
	if buf != nil {
	}
	return resp, nil
}

func ExampleExchangeDo_OK() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := exchangeDo(req)
	fmt.Printf("test: ExchangeDo_OK()-no-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*4)
	defer cancel()
	r2 := req.Clone(ctx)
	resp, err = exchangeDo(r2)
	fmt.Printf("test: ExchangeDo_OK()-5s-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: ExchangeDo_OK()-no-timeout -> [status-code:200] [err:<nil>]
	//test: ExchangeDo_OK()-5s-timeout -> [status-code:200] [err:<nil>]

}

func ExampleExchangeDo_Timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := Do(req)
	fmt.Printf("test: Do_Timeout() -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel2()
	req, _ = http.NewRequestWithContext(ctx2, http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, err = Do(req)
	fmt.Printf("test: Do() -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	buf, err1 := readAll(resp.Body)
	fmt.Printf("test: readAll() -> [buff:%v] [err:%v]\n", len(buf), err1)

	//Output:
	//test: Do_Timeout() -> [status-code:504] [err:Get "https://www.google.com/search?q=golang": context deadline exceeded]
	//test: Do() -> [status-code:200] [err:<nil>]
	//test: readAll() -> [buff:83659] [err:<nil>]

}

func ExampleExchangeDo_URLError() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google2345.com/search?q=golang", nil)
	resp, err := exchangeDo(req)
	fmt.Printf("test: ExchangeDo_Timeout()-ReadAll()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: ExchangeDo_Timeout()-ReadAll()-timeout -> [status-code:500] [err:Get "https://www.google2345.com/search?q=golang": tls: first record does not look like a TLS handshake]

}
