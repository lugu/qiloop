package net

import (
	"fmt"
	"io"
	"log"
	gonet "net"
	"sync"
)

// Filter returns true if given message shall be processed by a
// Consumer. Returns two values:
// - matched: true if the message should be processed by the Consumer.
// - keep: true if the handler shall be kept in the dispatcher.
type Filter func(hdr *Header) (matched bool, keep bool)

// Consumer process a message which has been selected by a filter.
type Consumer func(msg *Message) error

// EndPoint reprensents a network socket capable of sending and
// receiving messages.
type EndPoint interface {

	// Send pushes the message into the network.
	Send(m Message) error

	// ReceiveAny returns a message. Should only be used during
	// bootstap and testing.
	ReceiveAny() (*Message, error)

	// AddHandler registers the associated Filter and Consumer to the
	// EndPoint. Do not attempt to add another handler from within a
	// Filter.
	AddHandler(f Filter, c Consumer) int

	// RemoveHandler removes the associated Filter and Consumer.
	// RemoveHandler must not be called from within the Filter: use
	// the Filter returned value keep for this purpose.
	RemoveHandler(id int) error
}

type endPoint struct {
	conn      gonet.Conn
	filters   []Filter
	consumers []Consumer
	// handlerMutex: protect filters and consumers which must stay synchronized.
	handlerMutex sync.Mutex
}

// NewEndPoint creates an EndPoint which accpets messsages
func NewEndPoint(conn gonet.Conn) EndPoint {
	e := &endPoint{
		conn:      conn,
		filters:   make([]Filter, 10),
		consumers: make([]Consumer, 10),
	}
	go e.process()
	return e
}

// DialEndPoint construct an endpoint by contacting a given address.
func DialEndPoint(addr string) (EndPoint, error) {
	conn, err := gonet.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return NewEndPoint(conn), nil
}

// Send post a message to the other side of the endpoint.
func (e *endPoint) Send(m Message) error {
	return m.Write(e.conn)
}

// Close wait for a message to be received.
func (e *endPoint) Close() error {
	return e.conn.Close()
}

// RemoveHandler unregister the associated Filter and Consumer.
// WARNING: RemoveHandler must not be called from within the Filter or
// the Consumer.
func (e *endPoint) RemoveHandler(id int) error {
	e.handlerMutex.Lock()
	defer e.handlerMutex.Unlock()
	if id >= 0 && id < len(e.filters) {
		e.filters[id] = nil
		e.consumers[id] = nil
		return nil
	}
	return fmt.Errorf("invalid handler id: %d", id)
}

// AddHandler register the associated Filter and Consumer to the
// EndPoint.
func (e *endPoint) AddHandler(f Filter, c Consumer) int {
	e.handlerMutex.Lock()
	defer e.handlerMutex.Unlock()
	for i, filter := range e.filters {
		if filter == nil {
			e.filters[i] = f
			e.consumers[i] = c
			return i
		}
	}
	e.filters = append(e.filters, f)
	e.consumers = append(e.consumers, c)
	return len(e.filters) - 1
}

// dispatch test all destinations for someone interrested in the
// message.
func (e *endPoint) dispatch(msg *Message) error {
	e.handlerMutex.Lock()
	defer e.handlerMutex.Unlock()
	for i, f := range e.filters {
		if f != nil {
			matched, keep := f(&msg.Header)
			consumer := e.consumers[i]
			if !keep {
				e.filters[i] = nil
				e.consumers[i] = nil
			}
			if matched {
				go consumer(msg)
				return nil
			}
		}
	}
	return fmt.Errorf("failed to dispatch message: %#v", msg.Header)
}

// process read all messages from the end point and dispatch them one
// by one.
func (e *endPoint) process() {
	queue := make(chan *Message, 10)
	defer close(queue)

	go func() {
		for msg := range queue {
			e.dispatch(msg)
		}
	}()

	for {
		msg := new(Message)
		err := msg.Read(e.conn)
		if err == io.EOF {
			e.Close()
			break
		} else if err != nil {
			// FIXME: proper error management: recover from a
			// currupted message by discarding the crap.
			log.Printf("closing connection: %s", err)
			e.Close()
			break
		}
		queue <- msg
	}
}

// Receive wait for a message to be received.
func (e *endPoint) ReceiveAny() (*Message, error) {
	found := make(chan *Message)
	filter := func(msg *Header) (matched bool, keep bool) {
		return true, false
	}
	consumer := func(msg *Message) error {
		found <- msg
		return nil
	}
	_ = e.AddHandler(filter, consumer)
	msg := <-found
	return msg, nil
}