package messaging

import (
	"fmt"
	"net/http"
)

func ExampleConfigContent() {
	cfg := make(map[string]string)
	cfg["hostname"] = "localhost"
	m := NewConfigMessage(cfg)

	cfg2, ok := ConfigContent[map[string]string](m)
	fmt.Printf("test: ConfigContent() -> [%v] [ok:%v]\n", cfg2, ok)

	cfg3, ok2 := ConfigContent[string](m)
	fmt.Printf("test: ConfigContent() -> [%v] [ok:%v]\n", cfg3, ok2)

	var cfg4 map[string]string

	ok2 = UpdateContent[map[string]string](m, &cfg4)
	fmt.Printf("test: UpdateContent() -> [%v] [ok:%v]\n", cfg4, ok2)

	//Output:
	//test: ConfigContent() -> [map[hostname:localhost]] [ok:true]
	//test: ConfigContent() -> [] [ok:false]
	//test: UpdateContent() -> [map[hostname:localhost]] [ok:true]

}

/*
func ExampleMapMessage() {
	cfg := make(map[string]string)
	cfg["key1"] = "value1"
	cfg["key2"] = "value2"
	m := NewMapMessage(cfg)

	t, status := MapContent(m)
	fmt.Printf("test: MapContent() -> [%v] [status:%v]\n", t, status)

	m.SetContent(ContentTypeMap, "invalid content")
	t2, status2 := MapContent(m)
	fmt.Printf("test: MapContent() -> [%v] [status:%v]\n", t2, status2)

	//Output:
	//test: MapContent() -> [map[key1:value1 key2:value2]] [status:OK]
	//test: MapContent() -> [map[]] [status:Invalid Content - error: content value type: string is not of generic type: map[string]string]

}


*/

func ExampleStatusMessage() {
	m := NewStatusMessage(http.StatusTeapot, "")

	status, event, result := StatusContent(m)
	fmt.Printf("test: StatusContent() -> [%v] [%v] [status:%v]\n", status, event, result)

	m.SetContent(ContentTypeStatus, "invalid content")
	status2, event2, result2 := StatusContent(m)
	fmt.Printf("test: StatusContent() -> [%v] [%v] [status:%v]\n", status2, event2, result2)

	//Output:
	//test: StatusContent() -> [418] [] [status:<nil>]
	//test: StatusContent() -> [0] [] [status:content value type: string is not of generic type: int]

}

/*
func ExampleAgentMessage() {
	a := newControlAgent("test:agent/example", nil)
	m := NewAgentMessage(a)

	a1, status := AgentContent(m)
	fmt.Printf("test: AgentContent() -> [%v] [status:%v]\n", a1, status)

	//Output:
	//test: AgentContent() -> [test:agent/example] [status:OK]

}


*/
