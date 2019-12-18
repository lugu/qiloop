package bus

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/value"
)

type client struct {
	endpoint       net.EndPoint
	messageID      uint32
	messageIDMutex sync.Mutex
	state          map[string]int
	stateMutex     sync.Mutex
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

func (c *client) cancelMessage(hdr net.Header) net.Message {
	msg := net.NewMessage(hdr, []byte{})
	msg.Header.Type = net.Cancel
	return msg
}

func (c *client) Call(cancel <-chan struct{}, serviceID, objectID, actionID uint32,
	payload []byte) ([]byte, error) {

	// Do nothing if cancel is already closed.
	if cancel != nil {
		select {
		case <-cancel:
			return nil, ErrCancelled
		default:
		}
	}

	msg := c.newMessage(serviceID, objectID, actionID, payload)
	messageID := msg.Header.ID

	reply := make(chan *net.Message, 1)
	errors := make(chan error, 1)

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == serviceID && hdr.Object == objectID &&
			hdr.Action == actionID && hdr.ID == messageID {
			return true, false
		}
		return false, true
	}

	closer := func(err error) {
	    if err != nil {
		errors <- err
	    }
	}

	// 1. starts listening for an answer.
	id := c.endpoint.MakeHandler(filter, reply, closer)

	// 2. send the call message.
	if err := c.endpoint.Send(msg); err != nil {
		c.endpoint.RemoveHandler(id)
		return nil, fmt.Errorf(
			"call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}

	// 3. wait for a response
	var (
		ok       bool
		response *net.Message
	)

	if cancel == nil {
		cancel = make(chan struct{})
	}

	select {
	case err := <-errors:
		return nil, err
	case response, ok = <-reply:
		if !ok {
			return nil, fmt.Errorf("Remote connection closed")
		}
	case <-cancel:
		msg := c.cancelMessage(msg.Header)
		if err := c.endpoint.Send(msg); err != nil {
			return nil, fmt.Errorf(
				"cancel failed: service %d, object %d, action %d: %s",
				serviceID, objectID, actionID, err)
		}
		return nil, ErrCancelled
	}
	// 4. analyse response
	switch response.Header.Type {
	case net.Reply:
		return response.Payload, nil
	case net.Error:
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
	case net.Cancelled:
		return nil, ErrCancelled
	default:
		return nil, fmt.Errorf("Unexpected message type: %d",
			response.Header.Type)
	}
}

// Subscribe returns a channel which returns the future value of a
// given signal. Use the cancel method to stop the stream.
// Do not close the events channel.
func (c *client) Subscribe(serviceID, objectID, actionID uint32) (
	cancel func(), events chan []byte, err error) {

	abort := make(chan struct{})
	events = make(chan []byte)
	cancel = func() {
		close(abort)
	}

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == serviceID && hdr.Object == objectID &&
			hdr.Action == actionID {
			if hdr.Type == net.Error {
				// unsubscribe on error
				return true, false
			}
			return true, true
		}
		return false, true
	}
	queue := make(chan *net.Message, 10)

	go func(id int) {
		for {
			select {
			case msg, ok := <-queue:
				if !ok {
					close(events)
					return
				} else if msg.Header.Type == net.Event {
					events <- msg.Payload
				}
			case <-abort:
				c.endpoint.RemoveHandler(id)
				close(events)
				return
			}
		}
	}(c.endpoint.MakeHandler(filter, queue, nil))

	return cancel, events, nil
}

// OnDisconnect calls cb when the network connection is closed.
func (c *client) OnDisconnect(closer func(error)) error {
	if closer == nil {
		return nil
	}
	filter := func(hdr *net.Header) (bool, bool) { return false, true }
	consumer := make(chan *net.Message)
	c.endpoint.MakeHandler(filter, consumer, closer)
	return nil
}

func (c *client) State(signal string, add int) int {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	previous, ok := c.state[signal]
	if !ok && add != 0 {
		c.state[signal] = add
		return add
	}
	next := previous + add
	if next == 0 {
		delete(c.state, signal)
		return 0
	}
	c.state[signal] = next
	return next
}

// NewClient returns a new client.
func NewClient(endpoint net.EndPoint) Client {
	return &client{
		endpoint:  endpoint,
		messageID: 1,
		state:     map[string]int{},
	}
}

// SelectEndPoint connect to a remote peer using the list of
// addresses. It tries local addresses first and refuses to connect
// invalid IP addresses such as test ranges (198.18.0.x).
// User and token are user during Authentication. If user and token
// are empty, the file .qiloop-auth.conf is read.
func SelectEndPoint(addrs []string, user, token string) (addr string, endpoint net.EndPoint, err error) {
	if len(addrs) == 0 {
		return "", nil, fmt.Errorf("empty address list")
	}
	// sort the addresses based on their value
	for _, addr := range addrs {
		// do not connect the test range.
		// TODO: unless a local interface has such IP
		// address.
		if strings.Contains(addr, "198.18.0") {
			continue
		}
		endpoint, err = net.DialEndPoint(addr)
		if err != nil {
			continue
		}
		if user == "" && token == "" {
			err = Authenticate(endpoint)
		} else {
			err = AuthenticateUser(endpoint, user, token)
		}
		if err != nil {
			return "", nil, fmt.Errorf("authentication error: %s",
				err)
		}
		return addr, endpoint, nil
	}
	return "", nil, err
}
