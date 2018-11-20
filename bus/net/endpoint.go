package net

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	gonet "net"
	"net/url"
	"strings"
	"sync"
)

// Filter returns true if given message shall be processed by a
// Consumer. Returns two values:
// - matched: true if the message should be processed by the Consumer.
// - keep: true if the handler shall be kept in the dispatcher.
type Filter func(hdr *Header) (matched bool, keep bool)

// Consumer process a message which has been selected by a filter.
type Consumer func(msg *Message) error

// Closer informs the handler about a disconnection
type Closer func(err error)

// EndPoint reprensents a network socket capable of sending and
// receiving messages.
type EndPoint interface {

	// Send pushes the message into the network.
	Send(m Message) error

	// ReceiveAny returns a message. Should only be used during
	// bootstap and testing. Deprecated.
	ReceiveAny() (*Message, error)

	// ReceiveOne returns a chanel to receive a message.
	ReceiveOne() (chan *Message, error)

	// AddHandler registers the associated Filter and Consumer to the
	// EndPoint. Do not attempt to add another handler from within a
	// Filter.
	AddHandler(f Filter, c Consumer, cl Closer) int

	// RemoveHandler removes the associated Filter and Consumer.
	// RemoveHandler must not be called from within the Filter: use
	// the Filter returned value keep for this purpose.
	RemoveHandler(id int) error

	// Close close the underlying connection
	Close() error

	String() string
}

type Handler struct {
	filter   Filter
	consumer Consumer
	closer   Closer
	queue    chan *Message
	err      error
}

func NewHandler(f Filter, c Consumer, cl Closer) *Handler {
	h := &Handler{
		filter:   f,
		consumer: c,
		closer:   cl,
		queue:    make(chan *Message, 10),
	}
	go h.run()
	return h
}

func (h *Handler) run() {
	for {
		msg, ok := <-h.queue
		if !ok {
			h.closer(h.err)
			return
		}
		err := h.consumer(msg)
		if err != nil {
			log.Printf("failed to consume message: %s", err)
		}
	}
}

type endPoint struct {
	conn          gonet.Conn
	handlers      []*Handler
	handlersMutex sync.Mutex
}

// EndPointFinalizer creates a new EndPoint and let you process it
// before it start handling messages. This allows you to add handler
// or/and avoid data races.
func EndPointFinalizer(conn gonet.Conn, finalizer func(EndPoint)) EndPoint {
	e := &endPoint{
		conn:     conn,
		handlers: make([]*Handler, 10),
	}
	finalizer(e)
	go e.process()
	return e
}

func NewEndPoint(conn gonet.Conn) EndPoint {
	e := &endPoint{
		conn:     conn,
		handlers: make([]*Handler, 10),
	}
	go e.process()
	return e
}

func dialUNIX(name string) (EndPoint, error) {
	conn, err := gonet.Dial("unix", name)
	if err != nil {
		return nil, fmt.Errorf(`failed to connect unix socket "%s": %s`,
			name, err)
	}
	return NewEndPoint(conn), nil
}

func dialTCP(addr string) (EndPoint, error) {
	conn, err := gonet.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return NewEndPoint(conn), nil
}

// dialTCPS connects regardless of the certificate.
func dialTCPS(addr string) (EndPoint, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return NewEndPoint(conn), nil
}

// DialEndPoint construct an endpoint by contacting a given address.
func DialEndPoint(addr string) (EndPoint, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return dialTCP(addr)
	} else {
		switch u.Scheme {
		case "tcp":
			return dialTCP(u.Host)
		case "tcps":
			return dialTCPS(u.Host)
		case "unix":
			return dialUNIX(strings.TrimPrefix(addr, "unix://"))
		default:
			return nil, fmt.Errorf("unknown URL scheme: %s", addr)
		}
	}
}

// Send post a message to the other side of the endpoint.
func (e *endPoint) Send(m Message) error {
	return m.Write(e.conn)
}

// closeWith close all handler
func (e *endPoint) closeWith(err error) error {

	ret := e.conn.Close()

	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()

	for id, handler := range e.handlers {
		if handler != nil {
			handler.err = err
			close(handler.queue)
			e.handlers[id] = nil
		}
	}
	return ret
}

// Close wait for a message to be received.
func (e *endPoint) Close() error {
	return e.closeWith(nil)
}

// RemoveHandler unregister the associated Filter and Consumer.
// WARNING: RemoveHandler must not be called from within the Filter or
// the Consumer.
func (e *endPoint) RemoveHandler(id int) error {
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	if id >= 0 && id < len(e.handlers) && e.handlers[id] != nil {
		close(e.handlers[id].queue)
		e.handlers[id] = nil
		return nil
	}
	return fmt.Errorf("invalid handler id: %d", id)
}

// AddHandler register the associated Filter and Consumer to the
// EndPoint.
func (e *endPoint) AddHandler(f Filter, c Consumer, cl Closer) int {
	newHandler := NewHandler(f, c, cl)
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	for i, handler := range e.handlers {
		if handler == nil {
			e.handlers[i] = newHandler
			return i
		}
	}
	e.handlers = append(e.handlers, newHandler)
	return len(e.handlers) - 1
}

var MessageDropped error = errors.New("message dropped")

// dispatch test all destinations for someone interrested in the
// message.
func (e *endPoint) dispatch(msg *Message) error {
	dispatched := false
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	for i, h := range e.handlers {
		if h == nil {
			continue
		}
		matched, keep := h.filter(&msg.Header)
		if matched {
			h.queue <- msg
			dispatched = true
		}
		if !keep {
			close(h.queue)
			e.handlers[i] = nil
		}
	}
	if dispatched {
		return nil
	}
	return fmt.Errorf("dropping message: %#v (len handlers: %d)",
		msg.Header, len(e.handlers))
}

// process read all messages from the end point and dispatch them one
// by one.
func (e *endPoint) process() {
	var err error

	for {
		msg := new(Message)
		err = msg.Read(e.conn)
		if err != nil {
			e.closeWith(err)
			return
		}
		err = e.dispatch(msg)
		if err != nil {
			log.Printf("dispatch err: %s\n", err)
		}
	}
}

// ReceiveOne returns a chanel to receive one message. If the
// connection close, the chanel is closed.
func (e *endPoint) ReceiveOne() (chan *Message, error) {
	found := make(chan *Message, 1)
	filter := func(hdr *Header) (matched bool, keep bool) {
		return true, false
	}
	consumer := func(msg *Message) error {
		found <- msg
		return nil
	}
	closer := func(err error) {
		close(found)
	}
	_ = e.AddHandler(filter, consumer, closer)
	return found, nil
}

// Receive wait for a message to be received.
func (e *endPoint) ReceiveAny() (*Message, error) {
	found := make(chan *Message)
	filter := func(hdr *Header) (matched bool, keep bool) {
		return true, false
	}
	consumer := func(msg *Message) error {
		found <- msg
		return nil
	}
	closer := func(err error) {
		close(found)
	}
	_ = e.AddHandler(filter, consumer, closer)
	msg, ok := <-found
	if !ok {
		return nil, io.EOF
	}
	return msg, nil
}

func (e *endPoint) String() string {
	return e.conn.RemoteAddr().Network() + "://" +
		e.conn.RemoteAddr().String()
}

func NewPipe() (EndPoint, EndPoint) {
	a, b := gonet.Pipe()
	return NewEndPoint(a), NewEndPoint(b)
}
