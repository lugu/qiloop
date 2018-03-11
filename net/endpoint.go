package net

import (
	"fmt"
	"io"
	"log"
	gonet "net"
	"sync"
)

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const ErrNotForMe = ConstError("This message is not for me")

// Handler consum messages. If a message is not destinated to the
// Handle, it must return ErrNotForMe otherwise it must return nil.
// Messages are passed through handlers until one handler returns nil.
type Handler func(msg *Message) error

// EndPoint reprensents a network socket capable of sending and
// receiving messages.
type EndPoint interface {
	Send(m Message) error
	ReceiveAny() (*Message, error)

	// TODO: split Handler into two parts: a filter for the message
	// header which must not block and another part which does the
	// processing which will be run it its own goroutine.
	AddHandler(h Handler) int
	RemoveHandler(id int) error
}

type endPoint struct {
	conn         gonet.Conn
	destinations []Handler
	destMutex    sync.Mutex
}

// NewEndPoint creates an EndPoint which accpets messsages
func NewEndPoint(conn gonet.Conn) EndPoint {
	e := &endPoint{
		conn:         conn,
		destinations: make([]Handler, 10),
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

func (e *endPoint) RemoveHandler(id int) error {
	e.destMutex.Lock()
	defer e.destMutex.Unlock()
	if id >= 0 && id < len(e.destinations) {
		e.destinations[id] = nil
		return nil
	}
	return fmt.Errorf("invalid Handler id: %d", id)
}

func (e *endPoint) AddHandler(h Handler) int {
	e.destMutex.Lock()
	defer e.destMutex.Unlock()
	for i, handler := range e.destinations {
		if handler == nil {
			e.destinations[i] = h
			return i
		}
	}
	e.destinations = append(e.destinations, h)
	return len(e.destinations) - 1
}

// dispatch test all destinations for someone interrested in the
// message.
//
// BUG: do not holds the destMutex lock during processing since one
// of the Handler might want to unregister itself.
func (e *endPoint) dispatch(msg *Message) error {
	e.destMutex.Lock()
	defer e.destMutex.Unlock()
	for _, dest := range e.destinations {
		if dest != nil && dest(msg) == nil {
			return nil
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
	var handler Handler = func(msg *Message) error {
		found <- msg
		return nil
	}
	id := e.AddHandler(handler)
	defer e.RemoveHandler(id)
	msg := <-found
	return msg, nil
}
