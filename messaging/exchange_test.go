package messaging

import (
	"fmt"
	"time"
)

var (
	pkgPath = "/github/appellative-ai/core/messaging"
)

type agentC struct {
	running bool
	name    string
	ch      chan *Message
	handler Handler
}

func newControlAgent(name string, handler Handler) *agentC {
	c := new(agentC)
	c.name = name
	c.ch = make(chan *Message, ChannelSize)
	if handler == nil {
		c.handler = func(m *Message) {}
	} else {
		c.handler = handler
	}
	return c
}

// Name -
func (c *agentC) Name() string { return c.name }

// String - identity
func (c *agentC) String() string { return c.Name() }

// Message - message an agent
func (c *agentC) Message(msg *Message) {
	if msg == nil {
		return
	}
	if msg.Name == StartupEvent {
		c.run()
		return
	}
	if !c.running {
		return
	}
	switch msg.Channel() {
	case ChannelControl:
		if c.ch != nil {
			c.ch <- msg
		}
	default:
	}
}

// Run - run the agent
func (c *agentC) run() {
	if c.running {
		return
	}
	c.running = true
	go controlAgentRun(c)
}

// Shutdown - shutdown the agent
func (c *agentC) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	c.Message(ShutdownMessage)
}

func (c *agentC) shutdown() {
	close(c.ch)
}

// controlAgentRun - a simple run function that only handles control messages, and dispatches via a message handler
func controlAgentRun(c *agentC) {
	if c == nil || c.handler == nil {
		return
	}
	for {
		select {
		case msg, open := <-c.ch:
			if !open {
				return
			}
			fmt.Printf("test: controlAgent.run() -> %v\n", msg)
			switch msg.Name {
			case ShutdownEvent:
				c.handler(NewMessage(ChannelControl, msg.Name))
				c.shutdown()
				return
			default:
				c.handler(msg)
			}
		default:
		}
	}
}

func emptyHandler(_ *Message) {}

func ExampleRegister() {
	testDir := NewExchange()

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())

	uri1 := "urn:test:one"
	a := testDir.Get(uri1)
	fmt.Printf("test: Get(%v) -> : [agent:%v]\n", uri1, a)

	a1 := newControlAgent(uri1, emptyHandler)
	err := testDir.Register(a1)
	fmt.Printf("test: Register(%v) -> : [err:%v]\n", uri1, err)

	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m1 := testDir.Get(uri1)
	fmt.Printf("test: Get(%v) -> : [agent:%v]\n", uri1, m1.Name())

	uri2 := "urn:test:two"
	a2 := newControlAgent(uri2, emptyHandler)
	err = testDir.Register(a2)
	fmt.Printf("test: Register(%v) -> : [err:%v]\n", uri2, err)
	fmt.Printf("test: Count() -> : %v\n", testDir.Count())
	m2 := testDir.Get(uri2)
	fmt.Printf("test: Get(%v) -> : [agent:%v]\n", uri2, m2.Name())

	fmt.Printf("test: List() -> : %v\n", testDir.List())

	//Output:
	//test: Count() -> : 0
	//test: Get(urn:test:one) -> : [agent:<nil>]
	//test: Register(urn:test:one) -> : [err:<nil>]
	//test: Count() -> : 1
	//test: Get(urn:test:one) -> : [agent:urn:test:one]
	//test: Register(urn:test:two) -> : [err:<nil>]
	//test: Count() -> : 2
	//test: Get(urn:test:two) -> : [agent:urn:test:two]
	//test: List() -> : [urn:test:one urn:test:two]

}

func ExampleRegisterError() {
	uri := "urn:agent007"
	ex := NewExchange()

	a := newControlAgent(uri, emptyHandler)
	err := ex.Register(a)
	fmt.Printf("test: Register(%v) -> [%v]\n", uri, err)

	err = ex.Register(a)
	fmt.Printf("test: Register(%v) -> [%v]\n", uri, err)

	//Output:
	//test: Register(urn:agent007) -> [<nil>]
	//test: Register(urn:agent007) -> [exchange.Register() agent already exists: [urn:agent007]]

}

func ExampleMessageError() {
	uri := "urn:test"
	ex := NewExchange()

	fmt.Printf("test: Message(%v) -> : %v\n", uri, ex.Message(nil))
	fmt.Printf("test: Message(%v) -> : %v\n", uri, ex.Message(NewMessage(ChannelControl, "")))
	fmt.Printf("test: Message(%v) -> : %v\n", uri, ex.Message(NewMessage(ChannelControl, "")))

	//Output:
	//test: Message(urn:test) -> : false
	//test: Message(urn:test) -> : false
	//test: Message(urn:test) -> : false

}

func ExampleExist() {
	name1 := "common:core:agent/test/exist"
	c := NewChannel("test", ChannelSize)
	a := newTestAgent(name1, c, nil)
	ex := NewExchange()

	fmt.Printf("test: Exists(%v) -> : %v\n", name1, ex.Exist(name1))
	ex.Register(a)
	fmt.Printf("test: Exists(%v) -> : %v\n", name1, ex.Exist(name1))
	name1 = "bad:agent"
	fmt.Printf("test: Exists(%v) -> : %v\n", name1, ex.Exist(name1))

	//Output:
	//test: Exists(common:core:agent/test/exist) -> : false
	//test: Exists(common:core:agent/test/exist) -> : true
	//test: Exists(bad:agent) -> : false

}

func ExampleMessage() {
	uri1 := "urn:agent-1"
	uri2 := "urn:agent-2"
	uri3 := "urn:agent-3"
	c := NewChannel("test", ChannelSize) //make(chan *Message, 16)
	ex := NewExchange()

	a1 := newTestAgent(uri1, c, nil)
	ex.Register(a1)
	a2 := newTestAgent(uri2, c, nil)
	ex.Register(a2)
	a3 := newTestAgent(uri3, c, nil)
	ex.Register(a3)

	ex.Message(NewAddressableMessage(ChannelControl, StartupEvent, uri1, pkgPath))
	ex.Message(NewAddressableMessage(ChannelControl, StartupEvent, uri2, pkgPath))
	ex.Message(NewAddressableMessage(ChannelControl, StartupEvent, uri3, pkgPath))

	time.Sleep(time.Second * 1)
	resp1 := <-c.C
	resp2 := <-c.C
	resp3 := <-c.C
	fmt.Printf("test: <- c -> : %v %v %v\n", resp1.To(), resp2.To(), resp3.To())
	close(c.C)

	//Output:
	//test: <- c -> : [urn:agent-1] [urn:agent-2] [urn:agent-3]

}

func ExampleMessage_To() {
	name1 := "*:*:agent/test-1"
	name2 := "*:*:agent/test-2"
	name3 := "*:*:agent/test-3"
	c := NewChannel("test", ChannelSize) //make(chan *Message, 16)
	ex := NewExchange()

	a1 := newTestAgent(name1, c, nil)
	ex.Register(a1)
	a2 := newTestAgent(name2, c, nil)
	ex.Register(a2)
	a3 := newTestAgent(name3, c, nil)
	ex.Register(a3)

	m := NewMessage(ChannelControl, "event/test-multiple-to")
	m.AddTo(name1)
	sent := ex.Message(m)
	fmt.Printf("test: exchange.Message() -> [count:%v] [sent:%v]\n", len(m.To()), sent)

	m = NewMessage(ChannelControl, "event/test-multiple-to")
	m.AddTo(name1, name2, name3)
	sent = ex.Message(m)
	fmt.Printf("test: exchange.Message() -> [count:%v] [sent:%v]\n", len(m.To()), sent)

	//Output:
	//test: exchange.Message() -> [count:1] [sent:true]
	//test: exchange.Message() -> [count:3] [sent:true]

}

func ExampleListCount() {
	uri1 := "urn:agent-1"
	uri2 := "urn:agent-2"
	ex := NewExchange()

	a1 := newTestAgent(uri1, nil, nil)
	ex.Register(a1)
	a2 := newTestAgent(uri2, nil, nil)
	ex.Register(a2)

	fmt.Printf("test: Count() -> : %v\n", ex.Count())
	fmt.Printf("test: List() -> : %v\n", ex.List())

	//Output:
	//test: Count() -> : 2
	//test: List() -> : [urn:agent-1 urn:agent-2]

}

func _ExampleExchangeOnShutdown() {
	uri1 := "urn:agent-1"
	uri2 := "urn:agent-2"
	ex := NewExchange()

	a1 := newTestAgent(uri1, nil, nil)
	ex.Register(a1)
	a2 := newTestAgent(uri2, nil, nil)
	ex.Register(a2)

	fmt.Printf("test: Get(%v) -> : %v\n", uri1, ex.Get(uri1))
	fmt.Printf("test: Get(%v) -> : %v\n", uri2, ex.Get(uri2))

	a1.running = true
	a1.Shutdown()

	a2.running = true
	a2.Shutdown()

	fmt.Printf("test: Get-Shutdown(%v) -> : %v\n", uri1, ex.Get(uri1))
	fmt.Printf("test: Get-Shutdown(%v) -> : %v\n", uri2, ex.Get(uri2))

	ex2 := NewExchange()
	ex.Register(a1)
	ex2.Register(a1)

	fmt.Printf("test: Get-Ex1(%v) -> : %v\n", uri1, ex.Get(uri1))
	fmt.Printf("test: Get-Ex2(%v) -> : %v\n", uri1, ex.Get(uri1))

	a1.running = true
	a1.Shutdown()
	fmt.Printf("test: Get-Ex1-Shutdown(%v) -> : %v\n", uri1, ex.Get(uri1))
	fmt.Printf("test: Get-Ex2-Shutdown(%v) -> : %v\n", uri1, ex.Get(uri1))

	//Output:
	//test: Get(urn:agent-1) -> : urn:agent-1
	//test: Get(urn:agent-2) -> : urn:agent-2
	//test: Get-Shutdown(urn:agent-1) -> : <nil>
	//test: Get-Shutdown(urn:agent-2) -> : <nil>
	//test: Get-Ex1(urn:agent-1) -> : urn:agent-1
	//test: Get-Ex2(urn:agent-1) -> : urn:agent-1
	//test: Get-Ex1-Shutdown(urn:agent-1) -> : <nil>
	//test: Get-Ex2-Shutdown(urn:agent-1) -> : <nil>

}

/*
const (
	PingEvent        = "core:event/ping"
	ReconfigureEvent = "core:event/reconfigure"
)

func newAgentCtrlHandler(msg *Message) {
	fmt.Printf(fmt.Sprintf("test: NewControlAgent_CtrlHandler() -> %v\n", msg.Name()))
}

func ExampleNewControlAgent() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("test: NewControlAgent() -> [recovered:%v]\n", r)
		}
	}()
	//ex := NewExchange() //any(NewDirectory()).(*directory)
	//c := make(chan Message, 16)
	uri := "github.com/advanced-go/example-domain/activity"

	a, error0 := NewControlAgent(uri, newAgentCtrlHandler)
	if error0 != nil {
		fmt.Printf("test: NewControlAgent() -> [err:%v]\n", error0)
	}
	//status = a.Register(agentDir)
	//if !status.OK() {
	//	fmt.Printf("test: Register() -> [status:%v]\n", status)
	//}
	// 1 -10 Nanoseconds works for a direct send to a channel, sending via an controller2 needs a longer sleep time
	//d := time.Nanosecond * 10
	// Needed time.Nanoseconds * 50 for directory send with mutex
	// Needed time.Nanoseconds * 1 for directory send via sync.Map
	d := time.Nanosecond * 1
	//Startup(a)
	a.Message(NewMessage(ChannelControl, StartupEvent))
	//c <- Message{To: "", From: "", Event: aspect.StartupEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.Message(NewMessage(ChannelControl, PauseEvent))
	//c <- Message{To: "", From: "", Event: aspect.PauseEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.Message(NewMessage(ChannelControl, ResumeEvent))
	//c <- Message{To: "", From: "", Event: aspect.ResumeEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.Message(NewMessage(ChannelControl, PingEvent))
	//c <- Message{To: "", From: "", Event: aspect.PingEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.Message(NewMessage(ChannelControl, ReconfigureEvent))
	//c <- Message{To: "", From: "", Event: aspect.ReconfigureEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(d)
	a.Message(ShutdownMessage) //.SendCtrl(Message{To: uri, From: "", Event: aspect.ShutdownEvent})
	//c <- Message{To: "", From: "", Event: aspect.ShutdownEvent, RelatesTo: "", Status: nil, Content: nil, ReplyTo: nil}
	time.Sleep(time.Millisecond * 100)

	a.Message(ShutdownMessage)
	// will panic
	//c <- Message{}

	//Output:
	//test: NewControlAgent_CtrlHandler() -> core:event/pause
	//test: NewControlAgent_CtrlHandler() -> core:event/resume
	//test: NewControlAgent_CtrlHandler() -> core:event/ping
	//test: NewControlAgent_CtrlHandler() -> core:event/reconfigure
	//test: NewControlAgent_CtrlHandler() -> core:event/shutdown
	//test: NewControlAgent() -> [recovered:send on closed channel]

}


*/
