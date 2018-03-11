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

	found := make(chan []byte)

	var handler net.Handler = func(msg *net.Message) error {
		if msg.Header.Service == serviceID && msg.Header.Object == objectID &&
			msg.Header.Action == actionID && msg.Header.ID == messageID {
			found <- msg.Payload
			return nil
		}
		return net.ErrNotForMe
	}

	// 1. starts listening for an answer.
	id := c.endpoint.AddHandler(handler)
	defer c.endpoint.RemoveHandler(id)

	// 2. send the call message.
	if err := c.endpoint.Send(msg); err != nil {
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}

	// 3. wait for a response
	buf := <-found

	return buf, nil
}

const ErrCancelled = net.ConstError("Don't wait for a reply")

// Stream returns a channel which returns the future value of a
// given signal. To stop the stream one must send a value in the
// cancel channel. Do not close the message channel.
func (c *blockingClient) Stream(serviceID, objectID, actionID uint32, cancel chan int) (chan []byte, error) {
	found := make(chan []byte)

	var handler net.Handler = func(msg *net.Message) error {
		if msg.Header.Service == serviceID && msg.Header.Object == objectID &&
			msg.Header.Action == actionID {
			found <- msg.Payload
			return nil
		}
		return net.ErrNotForMe
	}

	go func(id int) {
		<-cancel
		c.endpoint.RemoveHandler(id)
	}(c.endpoint.AddHandler(handler))

	return found, nil
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
