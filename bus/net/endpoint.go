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
	// bootstap and testing.
	ReceiveAny() (*Message, error)

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
}

type endPoint struct {
	conn      gonet.Conn
	filters   []Filter
	consumers []Consumer
	closers   []Closer
	// handlerMutex: protect filters and consumers which must stay synchronized.
	handlerMutex sync.Mutex
}

// NewEndPoint creates an EndPoint which accpets messsages
func NewEndPoint(conn gonet.Conn) EndPoint {
	e := &endPoint{
		conn:      conn,
		filters:   make([]Filter, 0, 10),
		consumers: make([]Consumer, 0, 10),
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

// Close wait for a message to be received.
func (e *endPoint) closeWith(err error) error {
	e.handlerMutex.Lock()
	defer e.handlerMutex.Unlock()

	ret := e.conn.Close()
	var wait sync.WaitGroup
	wait.Add(len(e.closers))
	for id, c := range e.closers {
		if c == nil {
			wait.Done()
			continue
		}
		e.filters[id] = nil
		e.consumers[id] = nil
		e.closers[id] = nil
		go func(closer Closer) {
			closer(err)
			wait.Done()
		}(c)
	}
	wait.Wait()
	return ret
}

// Close wait for a message to be received.
func (e *endPoint) Close() error {
	return e.closeWith(io.EOF)
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
		e.closers[id] = nil
		return nil
	}
	return fmt.Errorf("invalid handler id: %d", id)
}

// AddHandler register the associated Filter and Consumer to the
// EndPoint.
func (e *endPoint) AddHandler(f Filter, c Consumer, cl Closer) int {
	e.handlerMutex.Lock()
	defer e.handlerMutex.Unlock()
	for i, filter := range e.filters {
		if filter == nil {
			e.filters[i] = f
			e.consumers[i] = c
			e.closers[i] = cl
			return i
		}
	}
	e.filters = append(e.filters, f)
	e.consumers = append(e.consumers, c)
	e.closers = append(e.closers, cl)
	return len(e.filters) - 1
}

var MessageDropped error = errors.New("message dropped")

// dispatch test all destinations for someone interrested in the
// message.
func (e *endPoint) dispatch(msg *Message) error {
	ret := MessageDropped
	e.handlerMutex.Lock()
	defer e.handlerMutex.Unlock()
	for i, f := range e.filters {
		if f == nil {
			continue
		}
		matched, keep := f(&msg.Header)
		consumer := e.consumers[i]
		if !keep {
			e.filters[i] = nil
			e.consumers[i] = nil
		}
		if matched {
			ret = nil
			go func() {
				if err := consumer(msg); err != nil {
					log.Printf("message consumer: %s", err)
				}
			}()
		}
	}
	if ret == nil {
		return nil
	}
	return fmt.Errorf("dropping message: %#v", msg.Header)
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
	var err error

	for {
		msg := new(Message)
		err := msg.Read(e.conn)
		if err == io.EOF {
			log.Printf("remote connection closed")
			break
		} else if err != nil {
			// FIXME: proper error management: recover from a
			// corrupted message by discarding the crap.
			log.Printf("error: closing connection: %s", err)
			break
		}
		queue <- msg
	}
	e.closeWith(err)
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

func NewPipe() (EndPoint, EndPoint) {
	a, b := gonet.Pipe()
	return NewEndPoint(a), NewEndPoint(b)
}
