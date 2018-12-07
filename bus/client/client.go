package client

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/value"
	"strings"
	"sync"
)

type client struct {
	endpoint       net.EndPoint
	messageID      uint32
	messageIDMutex sync.Mutex
}

func (c *client) nextMessageID() uint32 {
	c.messageIDMutex.Lock()
	defer c.messageIDMutex.Unlock()
	c.messageID += 2
	return c.messageID
}

func (c *client) newMessage(serviceID uint32, objectID uint32,
	actionID uint32, payload []byte) net.Message {

	header := net.NewHeader(net.Call, serviceID, objectID, actionID,
		c.nextMessageID())
	return net.NewMessage(header, payload)
}

func (c *client) Call(serviceID uint32, objectID uint32, actionID uint32,
	payload []byte) ([]byte, error) {

	msg := c.newMessage(serviceID, objectID, actionID, payload)
	messageID := msg.Header.ID

	reply := make(chan *net.Message)

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == serviceID && hdr.Object == objectID &&
			hdr.Action == actionID && hdr.ID == messageID {
			return true, false
		}
		return false, true
	}

	consumer := func(msg *net.Message) error {
		reply <- msg
		return nil
	}

	closer := func(err error) {
		// FIXME: remote connection closed: shall set
		// the client as unusable
		close(reply)
	}

	// 1. starts listening for an answer.
	id := c.endpoint.AddHandler(filter, consumer, closer)

	// 2. send the call message.
	if err := c.endpoint.Send(msg); err != nil {
		c.endpoint.RemoveHandler(id)
		return nil, fmt.Errorf(
			"failed to call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}

	// 3. wait for a response
	response, ok := <-reply
	if !ok {
		return nil, fmt.Errorf("Remote connection closed")
	}
	if response.Header.Type == net.Error {
		buf := bytes.NewBuffer(response.Payload)
		v, err := value.NewValue(buf)
		if err != nil {
			return nil, fmt.Errorf(
				"error response read error: %s", err)
		}
		strVal, ok := v.(value.StringValue)
		if !ok {
			return nil, fmt.Errorf("invalid error response")
		}
		return nil, fmt.Errorf(strVal.Value())
	}
	return response.Payload, nil
}

// Subscribe returns a channel which returns the future value of a
// given signal. To stop the stream one must send a value in the
// cancel channel. Do not close the message channel.
func (c *client) Subscribe(serviceID, objectID, actionID uint32,
	cancel chan int) (chan []byte, error) {

	stream := make(chan []byte)

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == serviceID && hdr.Object == objectID &&
			hdr.Action == actionID {
			return true, true
		}
		return false, true
	}
	consumer := func(msg *net.Message) error {
		stream <- msg.Payload
		return nil
	}
	closer := func(err error) {
		defer func() {
			recover()
		}()
		close(stream)
		close(cancel)
	}

	go func(id int) {
		_, ok := <-cancel
		c.endpoint.RemoveHandler(id)
		if ok {
			close(stream)
		}
	}(c.endpoint.AddHandler(filter, consumer, closer))

	return stream, nil
}

// NewClient returns a new client.
func NewClient(endpoint net.EndPoint) bus.Client {
	return &client{
		endpoint:  endpoint,
		messageID: 1,
	}
}

// SelectEndPoint connect to a remote peer using the list of
// addresses. It tries local addresses first and refuses to connect
// invalid IP addresses such as test ranges (198.18.0.x).
func SelectEndPoint(addrs []string) (endpoint net.EndPoint, err error) {
	if len(addrs) == 0 {
		return endpoint, fmt.Errorf("empty address list")
	}
	// sort the addresses based on their value
	for _, addr := range addrs {
		// do not connect the test range.
		// FIXME: unless a local interface has such IP
		// address.
		if strings.Contains(addr, "198.18.0") {
			continue
		}
		endpoint, err = net.DialEndPoint(addr)
		if err != nil {
			continue
		}
		err = Authenticate(endpoint)
		if err != nil {
			return endpoint, fmt.Errorf("authentication error: %s",
				err)
		}
		return endpoint, nil
	}
	return endpoint, err
}
