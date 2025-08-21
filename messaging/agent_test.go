package messaging

import (
	"fmt"
	"time"
)

type testAgent struct {
	running bool
	//agentId    string
	name       string
	ctrl       *Channel
	data       *Channel
	handler    Handler
	shutdownFn func()
}

func NewTestAgent(name string) Agent {
	return newTestAgent(name, nil, nil)
}

func newTestAgent(name string, ctrl, data *Channel) *testAgent {
	t := new(testAgent)
	t.name = name
	if ctrl == nil {
		t.ctrl = NewChannel(ChannelData, ChannelSize)
	} else {
		t.ctrl = ctrl
	}
	if data == nil {
		t.data = NewChannel(ChannelControl, ChannelSize)
	} else {
		t.data = data
	}
	return t
}

// func (t *testAgent) IsFinalized() bool { return t.data.IsFinalized() && t.ctrl.IsFinalized() }
// func (t *testAgent) NUri() string    { return t.agentId }
func (t *testAgent) String() string { return t.Name() }
func (t *testAgent) Name() string   { return t.name }
func (t *testAgent) Message(msg *Message) {
	if msg == nil {
		return
	}
	switch msg.Channel() {
	case ChannelControl:
		if t.ctrl != nil {
			t.ctrl.C <- msg
		}
	case ChannelData:
		if t.data != nil {
			t.data.C <- msg
		}
	default:
	}
}
func (t *testAgent) Notify(statusCode int) { fmt.Printf("%v", statusCode) }
func (t *testAgent) Run() {
	if t.running {
		return
	}
	t.running = true
	go testAgentRun(t)
}

// Shutdown - shutdown the agent
func (t *testAgent) Shutdown() {
	if !t.running {
		return
	}
	t.running = false
	if t.shutdownFn != nil {
		t.shutdownFn()
	}
	t.Message(ShutdownMessage)
}

func testAgentRun(t *testAgent) {
	for {
		select {
		case msg, open := <-t.ctrl.C:
			if !open {
				return
			}
			fmt.Printf("test: AgentRun() -> %v\n", msg)
			if msg.Name == ShutdownEvent {
				return
			}
		default:
		}
		select {
		case msg, open := <-t.data.C:
			if !open {
				return
			}
			fmt.Printf("test: AgentRun() -> %v\n", msg)
		default:
		}
	}
}

/*
func printAgentRun(uri string, ctrl, data <-chan *Message, state any) {
	fmt.Printf("test: AgentRun() -> [uri:%v] [ctrl:%v] [data:%v] [state:%v]\n", uri, ctrl != nil, data != nil, state != nil)
}

func ExampleNewAgent_Error() {
	a, err := newAgent("", nil, nil, nil, nil)
	fmt.Printf("test: newAgent() -> [agent:%v] [%v]\n", a, err)

	a, err = newAgent("urn:agent7", nil, nil, nil, nil)
	fmt.Printf("test: newAgent() -> [agent:%v] [%v]\n", a, err)

	//Output:
	//test: newAgent() -> [agent:<nil>] [error: agent URI is empty]
	//test: newAgent() -> [agent:<nil>] [error: agent AgentFunc is nil]

}

func ExampleNewAgent() {
	uri := "urn:agent007"
	uri1 := "urn:agent009"

	a := newTestAgent(uri)
	a.Run()
	time.Sleep(time.Second)

	a, _ = NewAgentWithChannels(uri1, nil, nil, printAgentRun, "data")
	a.Run()
	time.Sleep(time.Second)

	//Output:
	//test: AgentRun() -> [uri:urn:agent007] [ctrl:true] [data:true] [state:false]
	//test: AgentRun() -> [uri:urn:agent009] [ctrl:true] [data:false] [state:true]

}

*/

func ExampleAgentRun() {
	uri := "urn:agent007"
	a := newTestAgent(uri, nil, nil)
	a.Run()
	a.Message(NewAddressableMessage(ChannelControl, StartupEvent, uri, "ExampleAgentRun()"))
	//a.Message(NewDataMessage(uri, "ExampleAgentRun()", DataEvent))
	time.Sleep(time.Second)
	a.Shutdown()
	time.Sleep(time.Second)

	//Output:
	//test: AgentRun() -> [chan:ctrl] [from:ExampleAgentRun()] [to:[urn:agent007]] [common:core:event/startup]
	//test: AgentRun() -> [chan:ctrl] [from:] [to:[]] [common:core:event/shutdown]

}

func ExampleName() {
	uri := "resiliency:agent/core/test"
	n := getName(uri)
	fmt.Printf("test: name(\"%v\") -> [%v]\n", uri, n)

	uri = "resiliency:agent/core/test#west1.zone1.host"
	n = getName(uri)
	fmt.Printf("test: name(\"%v\") -> [%v]\n", uri, n)

	//Output:
	//test: name("resiliency:agent/core/test") -> [resiliency:agent/core/test]
	//test: name("resiliency:agent/core/test#west1.zone1.host") -> [resiliency:agent/core/test]

}
