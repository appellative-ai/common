package messaging

type Channel struct {
	Name string
	C    chan *Message
}

func NewChannel(name string, size int) *Channel {
	c := new(Channel)
	c.Name = name
	c.C = make(chan *Message, size)
	return c
}

func NewEmissaryChannel() *Channel {
	return NewChannel(ChannelEmissary, ChannelSize)
}

func NewMasterChannel() *Channel {
	return NewChannel(ChannelMaster, ChannelSize)
}

/*
func NewPrimaryChannel(enabled bool) *Channel {
	return NewChannel(PrimaryChannel, enabled)
}
*/

func (c *Channel) String() string { return c.Name }

/*
func (c *Channel) Name() string   { return c.name }
func (c *Channel) IsClosed() bool { return c.C == nil }
*/

func (c *Channel) Close() {
	close(c.C)
}

func (c *Channel) Send(m *Message) {
	if m != nil {
		c.C <- m
	}
}

func (c *Channel) Receive() *Message {
	return <-c.C
}

//func (c *Channel) IsEnabled() bool { return c.enabled }
//func (c *Channel) Enable()         { c.enabled = true }
//func (c *Channel) Disable()        { c.enabled = false }
