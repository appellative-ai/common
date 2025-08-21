package messaging

import (
	"fmt"
	"strings"
)

const (
	ChannelSize        = 16
	FragmentIdentifier = "#"
)

// NewAgentFunc - agent constructor
type NewAgentFunc func() Agent

// Agent - agent
type Agent interface {
	Name() string
	Message(m *Message)
}

func Name(agent Agent) string {
	if agent == nil {
		return ""
	}
	return getName(agent.Name())
}

func getName(uri string) string {
	i := strings.Index(uri, FragmentIdentifier)
	if i == -1 {
		return uri
	}
	return uri[:i]
}

type agentT struct {
	name    string
	handler Handler
}

func NewAgent(name string, handler Handler) Agent {
	a := new(agentT)
	a.name = name
	if handler == nil {
		handler = func(m *Message) {
			fmt.Printf("NewAgent %v handler is nil\n", name)
		}
	}
	a.handler = handler
	return a
}

func (a *agentT) Name() string       { return a.name }
func (a *agentT) String() string     { return a.Name() }
func (a *agentT) Message(m *Message) { a.handler(m) }

/*
func Shutdown(agent Agent) {
	if agent != nil {
		agent.Message(ShutdownMessage)
	}
}

func Startup(agent Agent) {
	if agent != nil {
		agent.Message(StartupMessage)
	}
}


*/
