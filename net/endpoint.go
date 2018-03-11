package net

import (
	"fmt"
	"log"
	gonet "net"
)

// Handler consum messages. If a message is not destinated to the
// Handle, it must return an error. It must return nil if the message
// was processed.
type Handler func(msg *Message) error

// EndPoint reprensents a network socket capable of sending and
// receiving messages.
type EndPoint struct {
	conn         gonet.Conn
	destinations []Handler
	msg          chan *Message
}

// NewEndPoint creates an EndPoint which accpets messsages
// TODO: add a Start() Pause(), Terminate() methods
func NewEndPoint(conn gonet.Conn) EndPoint {
	e := EndPoint{
		conn:         conn,
		destinations: make([]Handler, 10),
	}
	// FIXME: not yet switched
	// e.process()
	return e
}

// DialEndPoint construct an endpoint by contacting a given address.
func DialEndPoint(addr string) (e EndPoint, err error) {
	conn, err := gonet.Dial("tcp", addr)
	if err != nil {
		return e, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return NewEndPoint(conn), nil
}

// AcceptedEndPoint construct an endpoint using an accepted
// connection.
func AcceptedEndPoint(c gonet.Conn) EndPoint {
	return NewEndPoint(c)
}

// Send post a message to the other side of the endpoint.
func (e EndPoint) Send(m Message) error {
	return m.Write(e.conn)
}

// Close wait for a message to be received.
func (e EndPoint) Close() error {
	return e.conn.Close()
}

func (e EndPoint) removeHandler(id int) error {
	if id >= 0 && id < len(e.destinations) {
		e.destinations[id] = nil
		return nil
	} else {
		return fmt.Errorf("invalid Handler id: %d", id)
	}
}

func (e EndPoint) addHandler(h Handler) int {
	for i, handler := range e.destinations {
		if handler == nil {
			e.destinations[i] = h
			return i
		}
	}
	e.destinations = append(e.destinations, h)
	return len(e.destinations) - 1
}

func (e EndPoint) dispatch(msg *Message) error {
	for _, dest := range e.destinations {
		if dest != nil && dest(msg) != nil {
			return nil
		}
	}
	return fmt.Errorf("failed to dispatch message: %#v", msg.Header)
}

// process read all messages from the end point and dispatch them one
// by one.
func (e EndPoint) process() {
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
		if err != nil {
			// TODO: stop when the connection closed
			log.Printf("message read error: %s", err)
			continue
		}
		queue <- msg
	}
}

func (e EndPoint) consume(chan *Message) {
}

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const ErrNotForMe = ConstError("This message is not for me")
const ErrCancelled = ConstError("Don't wait for a reply")

// ReceiveOne wait for a message to be received.
func (e EndPoint) ReceiveOne(serviceID, objectID, actionID, messageID uint32, cancel chan int) (m *Message, err error) {
	found := make(chan *Message)
	var handler Handler = func(msg *Message) error {
		if msg.Header.Service == serviceID && msg.Header.Object == objectID &&
			msg.Header.Action == actionID && msg.Header.ID == messageID {
			found <- msg
			return nil
		}
		return ErrNotForMe
	}
	id := e.addHandler(handler)
	select {
	case msg := <-found:
		e.removeHandler(id)
		return msg, nil
	case <-cancel:
		e.removeHandler(id)
		return m, ErrCancelled
	}
}

// Receive wait for a message to be received.
func (e EndPoint) Receive() (m Message, err error) {
	err = m.Read(e.conn)
	return
}
