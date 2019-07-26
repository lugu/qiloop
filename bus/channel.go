package bus

import (
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
)

type Channel interface {
	Cap() CapabilityMap
	EndPoint() net.EndPoint
	Authenticated() bool
	Send(msg *net.Message) error
	SendError(msg *net.Message, err error) error
	SendReply(msg *net.Message, response []byte) error
	SetAuthenticated()
}

// channel represent an established connection between a client and a
// server.
type channel struct {
	capability CapabilityMap
	endpoint   net.EndPoint
}

// NewContext retuns a non authenticate context.
func NewContext(e net.EndPoint) Channel {
	return &channel{
		endpoint:   e,
		capability: PreferedCap("", ""),
	}
}

// SendError send a error message in response to msg.
func (c *channel) SendError(msg *net.Message, err error) error {
	// FIXME: missing trace here.
	// o.trace(msg)
	return util.ReplyError(c.endpoint, msg, err)
}

// SendReply send a reply message in response to msg.
func (c *channel) SendReply(msg *net.Message, response []byte) error {
	hdr := msg.Header
	hdr.Type = net.Reply
	reply := net.NewMessage(hdr, response)
	// FIXME: missing trace here.
	// o.trace(&reply)
	return c.endpoint.Send(reply)
}

// Send send a reply message in response to msg.
func (c *channel) Send(msg *net.Message) error {
	// FIXME: missing trace here.
	// o.trace(&reply)
	return c.endpoint.Send(*msg)
}

// Authenticated returns true if the connection is authenticated.
func (c *channel) Authenticated() bool {
	return c.capability.Authenticated()
}

// SetAuthenticated marks the context as authenticated.
func (c *channel) SetAuthenticated() {
	c.capability.SetAuthenticated()
}

// Cap return the capability map associated with the channel.
func (c *channel) Cap() CapabilityMap {
	return c.capability
}

// EndPoint returns the other side endpoint.
func (c *channel) EndPoint() net.EndPoint {
	return c.endpoint
}
