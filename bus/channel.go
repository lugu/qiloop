package bus

import (
	"bytes"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/value"
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

func errorPaylad(err error) []byte {
	var buf bytes.Buffer
	val := value.String(err.Error())
	val.Write(&buf)
	return buf.Bytes()
}

// SendError send a error message in response to msg.
func (c *channel) SendError(msg *net.Message, err error) error {
	hdr := net.NewHeader(net.Error, msg.Header.Service, msg.Header.Object,
		msg.Header.Action, msg.Header.ID)
	mError := net.NewMessage(hdr, errorPaylad(err))
	return c.Send(&mError)
}

// SendReply send a reply message in response to msg.
func (c *channel) SendReply(msg *net.Message, response []byte) error {
	hdr := msg.Header
	hdr.Type = net.Reply
	reply := net.NewMessage(hdr, response)
	return c.Send(&reply)
}

// Send send a reply message in response to msg.
func (c *channel) Send(msg *net.Message) error {
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

// tracedChannel notify a Tracer each time a message is directly sent using
// SendReply, SendError or Send.
type tracedChannel struct {
	Channel
	tracer Tracer
	id     uint32
}

func (c *tracedChannel) Send(msg *net.Message) error {
	c.tracer.Trace(msg, c.id)
	return c.Channel.Send(msg)
}

func (c *tracedChannel) SendError(msg *net.Message, err error) error {
	hdr := net.NewHeader(net.Error, msg.Header.Service, msg.Header.Object,
		msg.Header.Action, msg.Header.ID)
	mError := net.NewMessage(hdr, errorPaylad(err))
	return c.Send(&mError)
}

func (c *tracedChannel) SendReply(msg *net.Message, response []byte) error {
	hdr := msg.Header
	hdr.Type = net.Reply
	reply := net.NewMessage(hdr, response)
	return c.Send(&reply)
}
