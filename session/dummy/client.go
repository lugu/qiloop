package dummy

import (
	"fmt"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/session"
	"sync"
)

type blockingClient struct {
	endpoint       net.EndPoint
	messageID      uint32
	messageIDMutex sync.Mutex
}

func (c *blockingClient) nextMessageID() uint32 {
	c.messageIDMutex.Lock()
	defer c.messageIDMutex.Unlock()
	c.messageID += 2
	return c.messageID
}

func (c *blockingClient) newMessage(serviceID uint32, objectID uint32, actionID uint32, payload []byte) net.Message {
	header := net.NewHeader(net.Call, serviceID, objectID, actionID, c.nextMessageID())
	return net.NewMessage(header, payload)
}

func (c *blockingClient) Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error) {

	msg := c.newMessage(serviceID, objectID, actionID, payload)
	messageID := msg.Header.ID

	reply := make(chan []byte)

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == serviceID && hdr.Object == objectID &&
			hdr.Action == actionID && hdr.ID == messageID {
			return true, false
		}
		return false, true
	}

	consumer := func(msg *net.Message) error {
		reply <- msg.Payload
		return nil
	}

	// 1. starts listening for an answer.
	id := c.endpoint.AddHandler(filter, consumer)

	// 2. send the call message.
	if err := c.endpoint.Send(msg); err != nil {
		c.endpoint.RemoveHandler(id)
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}

	// 3. wait for a response
	buf := <-reply

	return buf, nil
}

// Stream returns a channel which returns the future value of a
// given signal. To stop the stream one must send a value in the
// cancel channel. Do not close the message channel.
func (c *blockingClient) Stream(serviceID, objectID, actionID uint32, cancel chan int) (chan []byte, error) {
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

	go func(id int) {
		<-cancel
		c.endpoint.RemoveHandler(id)
		close(stream)
	}(c.endpoint.AddHandler(filter, consumer))

	return stream, nil
}

func NewClient(addr string) (session.Client, error) {
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return nil, fmt.Errorf("client failed to connect %s: %s", endpoint, err)
	}
	return &blockingClient{
		endpoint:  endpoint,
		messageID: 1,
	}, nil
}

func newClient(endpoint net.EndPoint) session.Client {
	return &blockingClient{
		endpoint:  endpoint,
		messageID: 3,
	}
}
