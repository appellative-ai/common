package core

import (
	"context"
	"fmt"
	"net/http"
)

func ExampleBuildNetwork() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := BuildNetwork([]any{do1ExchangeFn, do2ExchangeFn, do3ExchangeFn})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func ExampleBuildNetworkExchange_Link() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1ExchangeFn, do2ExchangeFn, do3ExchangeFn})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func ExampleBuildNetworkExchange_Chainable() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1Exchange{}, do2Exchange{}, do3Exchange{}})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func ExampleBuildNetworkExchange_Any() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1Exchange{}, do2ExchangeFn, do3Exchange{}, do4ExchangeFn})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do4-Exchange() -> request
	//test: Do4-Exchange() -> response
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func ExampleBuildNetworkExchange_Abbreviated() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1ExchangeFn, do2ExchangeFn, do3ExchangeFailFn, do4ExchangeFn})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange-Fail() -> request
	//test: Do3-Exchange-Fail() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func ExampleBuildNetwork_Exchange() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1ExchangeFn, do2ExchangeFn, do3ExchangeFn})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func ExampleBuildNetworkExchange_Exchangeable() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1ExchangeFn, do2ExchangeFn, do3ExchangeFn})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func do1ExchangeFn(next Exchange) Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		fmt.Printf("test: Do1-Exchange() -> request\n")
		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		fmt.Printf("test: Do1-Exchange() -> response\n")
		return
	}
}

func do2ExchangeFn(next Exchange) Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		fmt.Printf("test: Do2-Exchange() -> request\n")
		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		fmt.Printf("test: Do2-Exchange() -> response\n")
		return
	}
}

func do3ExchangeFn(next Exchange) Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		fmt.Printf("test: Do3-Exchange() -> request\n")
		//fmt.Printf("test: Do3() -> response\n")
		//return &http.Response{StatusCode: http.StatusBadRequest}, nil
		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		fmt.Printf("test: Do3-Exchange() -> response\n")
		return
	}
}

func do3ExchangeFailFn(next Exchange) Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		fmt.Printf("test: Do3-Exchange-Fail() -> request\n")
		fmt.Printf("test: Do3-Exchange-Fail() -> response\n")
		return &http.Response{StatusCode: http.StatusBadRequest}, nil
		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		fmt.Printf("test: Do3-Exchange-Fail() -> response\n")
		return
	}
}

func do4ExchangeFn(next Exchange) Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		fmt.Printf("test: Do4-Exchange() -> request\n")
		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		fmt.Printf("test: Do4-Exchange() -> response\n")
		return
	}
}

type do1Exchange struct{}

func (d do1Exchange) Link(next Exchange) Exchange {
	return do1ExchangeFn(next)
}

type do2Exchange struct{}

func (d do2Exchange) Link(next Exchange) Exchange {
	return do2ExchangeFn(next)
}

type do3Exchange struct{}

func (d do3Exchange) Link(next Exchange) Exchange {
	return do3ExchangeFn(next)
}

type do3FailExchange struct{}

func (d do3FailExchange) Link(next Exchange) Exchange {
	return do3ExchangeFailFn(next)
}

type do4Exchange struct{}

func (d do4Exchange) Link(next Exchange) Exchange {
	return do4ExchangeFn(next)
}

func ExampleBuildNetwork_Chainable() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1Exchange{}, do2Exchange{}, do3Exchange{}})
	ex(req)

	//Output:
	//test: Do1-Exchange() -> request
	//test: Do2-Exchange() -> request
	//test: Do3-Exchange() -> request
	//test: Do3-Exchange() -> response
	//test: Do2-Exchange() -> response
	//test: Do1-Exchange() -> response

}

func _ExampleBuildNetwork_Panic_Nil() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{})
	ex(req)

	//Output:
	//fail
}

func _ExampleBuildNetwork_Panic_Nil_Operative() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1Exchange{}, nil, do3Exchange{}, do4ExchangeFn})
	ex(req)

	//Output:
	//fail
}

func _ExampleBuildNetwork_Panic_Type() {
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)
	ex := buildNetwork[Exchange, Chainable[Exchange]]([]any{do1Exchange{}, req, do3Exchange{}, do4ExchangeFn})
	ex(req)

	//Output:
	//fail
}

/*
type do1Combined struct{}

func (d do1Combined) Link(next messaging.Handler) messaging.Handler {
	return do1HandlerFn(next)
}

func (d do1Combined) doExchange(next Exchange) Exchange {
	return do1ExchangeFn(next)
}

func ExampleBuildNetwork_Combined() {
	rec := BuildNetwork[messaging.Handler, Chainable[messaging.Handler]]([]any{do1Combined{}})
	rec(messaging.ShutdownMessage)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.google.com/search?q=golang", nil)

	// This will panic as do1Combined is not of type Chainable[Exchange]
	//ex := BuildNetwork[Exchange, Chainable[Exchange]]([]any{do1Combined{}})
	//ex(req)

	// This works
	ex := BuildNetwork[Exchange, Chainable[Exchange]]([]any{do1Combined{}.doExchange})
	ex(req)

	//Output:
	//test: Do1-Handler() -> receive
	//test: Do1-Exchange() -> request
	//test: Do1-Exchange() -> response

}

*/
