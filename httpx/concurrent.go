package httpx

import (
	"context"
	"github.com/appellative-ai/common/core"
	"net/http"
	"sync"
	"time"
)

type Logger interface {
	Log(start time.Time, duration time.Duration, routeName string, req *http.Request, resp *http.Response, timeout time.Duration)
}

type Params struct {
	Name   string
	Req    *http.Request
	cancel func()
}

type Result struct {
	Name   string
	Req    *http.Request
	Resp   *http.Response
	Err    error
	cancel func()
}

func newResult(p Params) *Result {
	r := new(Result)
	r.Name = p.Name
	r.Req = p.Req
	r.cancel = p.cancel
	return r
}

func DoConcurrent[T Logger](do func(req *http.Request) (*http.Response, error), params ...Params) *core.MapT[string, *Result] {
	var wg sync.WaitGroup
	var t T

	m := core.NewSyncMap[string, *Result]()
	for _, p := range params {
		wg.Add(1)
		go func(r *Result) {
			defer wg.Done()
			start := time.Now().UTC()
			r.Resp, r.Err = do(r.Req)
			t.Log(start, time.Since(start), r.Name, r.Req, r.Resp, timeout(r.Req.Context()))
			m.Store(r.Name, r)
		}(newResult(p))
	}
	wg.Wait()
	return m
}

func timeout(ctx context.Context) time.Duration {
	var t time.Duration
	if d, ok := ctx.Deadline(); ok {
		t = time.Until(d)
	}
	return t
}
