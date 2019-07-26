package bus

import (
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
)

// Channel represent an established connection between a client and a
// server.
type Channel struct {
	Cap      CapabilityMap
	EndPoint net.EndPoint
}

// SendError send a error message in response to msg.
func (c *Channel) SendError(msg *net.Message, err error) error {
	// FIXME: missing trace here.
	// o.trace(msg)
	return util.ReplyError(c.EndPoint, msg, err)
}

// SendReply send a reply message in response to msg.
func (c *Channel) SendReply(msg *net.Message, response []byte) error {
	hdr := msg.Header
	hdr.Type = net.Reply
	reply := net.NewMessage(hdr, response)
	// FIXME: missing trace here.
	// o.trace(&reply)
	return c.EndPoint.Send(reply)
}

// Send send a reply message in response to msg.
func (c *Channel) Send(msg *net.Message) error {
	// FIXME: missing trace here.
	// o.trace(&reply)
	return c.EndPoint.Send(*msg)
}

// NewContext retuns a non authenticate context.
func NewContext(e net.EndPoint) *Channel {
	return &Channel{
		EndPoint: e,
		Cap:      PreferedCap("", ""),
	}
}

// Authenticated returns true if the connection is authenticated.
func (c *Channel) Authenticated() bool {
	return c.Cap.Authenticated()
}

// SetAuthenticated marks the context as authenticated.
func (c *Channel) SetAuthenticated() {
	c.Cap.SetAuthenticated()
}
