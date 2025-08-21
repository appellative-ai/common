package core

import (
	"fmt"
	"net/http"
)

func _ExampleNewEndpoint() {
	e := NewEndpoint("/resource/test", nil, nil, nil)

	_, ok := any(e).(http.Handler)
	fmt.Printf("test: NewEndpoint() -> [%v] [%v] [ServeHTTP:%v]\n", e, e.Pattern(), ok)

	//Output:
	//test: NewEndpoint() -> [&{/resource/test <nil> <nil> <nil>}] [/resource/test] [ServeHTTP:true]

}
