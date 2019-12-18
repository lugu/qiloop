package net

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	gonet "net"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/ftrvxmtrx/fd"
	quic "github.com/lucas-clemente/quic-go"
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
	fmt.Stringer

	// Send pushes the message into the network.
	Send(m Message) error

	// ReceiveAny returns a chanel to receive a single message.
	ReceiveAny() (chan *Message, error)

	// MakeHandler registers the associated Filter and queue to
	// the incomming traffic from the EndPoint. Writing to the
	// queue must not block otherwise messages are discarded. If
	// an error occurs, closer is called, the consumer is closed.
	// closer can be nil. Do not attempt to add another handler
	// from within a Filter. Filter must not block.
	MakeHandler(f Filter, queue chan<- *Message, cl Closer) int

	// AddHandler creates a queue of 10 messages and spawn a
	// goroutine to forward those events into the consumer, then
	// calls MakeHandler. Do not attempt to add another handler
	// from within a Filter. Filter must not block.
	AddHandler(f Filter, c Consumer, cl Closer) int

	// RemoveHandler removes the associated Filter and Consumer.
	// RemoveHandler must not be called from within the Filter: use
	// the Filter returned value keep for this purpose.
	RemoveHandler(id int) error

	// Close close the underlying connection
	Close() error
}

// Handler represents a client of a incomming message stream. It
// contains a filter used to decided if the message shall be sent to
// the handler, consumer will receive the messages in order of
// arrival. If an error occurs, closer is called, then the consumer
// is closed.
type Handler struct {
	filter   Filter
	consumer chan<- *Message
	closer   Closer
}

// NewHandler returns an Handler: f is call on each incomming message,
// c is called if f returns true. cl is always called when the handler
// is effectively closed.
func NewHandler(f Filter, ch chan<- *Message, cl Closer) *Handler {
	h := &Handler{
		filter:   f,
		consumer: ch,
		closer:   cl,
	}
	return h
}

// closeWith call the closer callback and then close the consumer.
// TODO: update API to call closer only on error.
func (h *Handler) closeWith(err error) {
	if h.closer != nil {
		h.closer(err)
	}
	close(h.consumer)
}

type endPoint struct {
	stream        Stream
	handlers      []*Handler
	handlersMutex sync.Mutex
}

// EndPointFinalizer creates a new EndPoint and let you process it
// before it start handling messages. This allows you to add handler
// or/and avoid data races.
func EndPointFinalizer(stream Stream, finalizer func(EndPoint)) EndPoint {
	e := &endPoint{
		stream:   stream,
		handlers: make([]*Handler, 10),
	}
	finalizer(e)
	go e.process()
	return e
}

// NewEndPoint returns an EndPoint which already process incomming
// messages. Since no handler have been register at the time of the
// creation of the EndPoint, any message receive will be droped until
// and Handler is registered. Prefer EndPointFinalizer for a safe way
// to construct EndPoint.
func NewEndPoint(stream Stream) EndPoint {
	e := &endPoint{
		stream:   stream,
		handlers: make([]*Handler, 10),
	}
	go e.process()
	return e
}

// ConnEndPoint returns an EndPoint using a go connection.
func ConnEndPoint(conn gonet.Conn) EndPoint {
	return NewEndPoint(ConnStream(conn))
}

func dialUNIX(name string) (EndPoint, error) {
	conn, err := gonet.Dial("unix", name)
	if err != nil {
		return nil, fmt.Errorf(`failed to connect unix socket "%s": %s`,
			name, err)
	}
	return ConnEndPoint(conn), nil
}

func dialTCP(addr string) (EndPoint, error) {
	conn, err := gonet.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("connect %s: %s", addr, err)
	}
	return ConnEndPoint(conn), nil
}

// dialTLS connects regardless of the certificate.
func dialTLS(addr string) (EndPoint, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", addr, conf)
	if err != nil {
		return nil, fmt.Errorf("connect %s: %s", addr, err)
	}
	return ConnEndPoint(conn), nil
}

// dialPipe connects a unix socket and echange file descriptors to
// communiate.
func dialPipe(name string) (EndPoint, error) {

	addr := gonet.UnixAddr{
		Name: name,
		Net:  "unix",
	}
	conn, err := gonet.DialUnix("unix", nil, &addr)
	if err != nil {
		return nil, fmt.Errorf(`failed to connect unix socket "%s": %s`,
			name, err)
	}
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	err = fd.Put(conn, r)
	if err != nil {
		return nil, err
	}
	fds, err := fd.Get(conn, 1, nil)
	if err != nil {
		return nil, err
	}
	if len(fds) != 1 {
		return nil, fmt.Errorf("missing fd")
	}
	return NewEndPoint(PipeStream(fds[0], w)), nil
}

// KeyNetContext represents an entry in the sream context.
type KeyNetContext uint32

const (
	// DialAddress is the addressed dialed.
	DialAddress KeyNetContext = iota
	// ListenAddress is the addressed listen.
	ListenAddress
)

// dialQUIC connects regardless of the certificate.
// TOOD: does not multiplex sessions
func dialQUIC(addr string) (EndPoint, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"qi-messaging"},
	}
	session, err := quic.DialAddr(addr, conf, nil)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(context.TODO(), DialAddress, addr)
	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		return nil, err
	}
	return NewEndPoint(newQuicStream(stream)), nil
}

// DialEndPoint construct an endpoint by contacting a given address.
func DialEndPoint(addr string) (EndPoint, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("dial: invalid address: %s", err)
	}
	switch u.Scheme {
	case "tcp":
		return dialTCP(u.Host)
	case "tcps":
		return dialTLS(u.Host)
	case "quic":
		return dialQUIC(u.Host)
	case "unix":
		return dialUNIX(strings.TrimPrefix(addr, "unix://"))
	case "pipe":
		return dialPipe(strings.TrimPrefix(addr, "pipe://"))
	default:
		return nil, fmt.Errorf("unknown URL scheme: %s", addr)
	}
}

// Send post a message to the other side of the endpoint.
func (e *endPoint) Send(m Message) error {
	return m.Write(e.stream)
}

// closeWith close all handler
func (e *endPoint) closeWith(err error) error {

	ret := e.stream.Close()

	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()

	for id, handler := range e.handlers {
		if handler != nil {
			go handler.closeWith(err)
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
		e.handlers[id].closeWith(nil)
		e.handlers[id] = nil
		return nil
	}
	return fmt.Errorf("invalid handler id: %d", id)
}

// AddDirectHandler register the associated Filter and Consumer to the
// EndPoint.
func (e *endPoint) MakeHandler(f Filter, queue chan<- *Message, cl Closer) int {
	newHandler := NewHandler(f, queue, cl)
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

// AddHandler register the associated Filter and Consumer to the
// EndPoint. A goroutine associated with a queue of 10 messages is
// created.
func (e *endPoint) AddHandler(f Filter, c Consumer, cl Closer) int {
	ch := make(chan *Message, 10)
	go func() {
		for msg := range ch {
			err := c(msg)
			if err != nil {
				log.Printf("consumer: %s", err)
			}
		}
	}()
	return e.MakeHandler(f, ch, cl)
}

// ErrNoMatch is returned when the message did not match any handler
var ErrNoMatch = errors.New("message dropped: no handler match")

// ErrConsumerBlocked is returned when the message cannot be sent to the
// consumer queue (the queue is full).
var ErrConsumerBlocked = errors.New("message dropped: consumer blocked")

// ErrNoHandler is returned when there is no handler registered.
var ErrNoHandler = errors.New("message dropped: no handler registered")

// dispatch requests each handler if it match the message, if so it
// sends it to the handler queue. If no handler is registered, it
// returns ErrNoHandler and if no handler match the message it returns
// ErrNoMatch.
func (e *endPoint) dispatch(msg *Message) error {
	e.handlersMutex.Lock()
	defer e.handlersMutex.Unlock()
	if len(e.handlers) == 0 {
		return ErrNoHandler
	}
	ret := ErrNoMatch
	for i, h := range e.handlers {
		if h == nil {
			continue
		}
		matched, keep := h.filter(&msg.Header)
		if matched {
			select {
			case h.consumer <- msg:
				if ret == ErrNoMatch {
					ret = nil
				}
			default:
				ret = ErrConsumerBlocked
			}
		}
		if !keep {
			h.closeWith(nil)
			e.handlers[i] = nil
		}
	}
	return ret
}

// process read all messages from the end point and dispatch them one
// by one.
func (e *endPoint) process() {
	var err error

	for {
		msg := new(Message)
		err = msg.Read(e.stream)
		if err != nil {
			e.closeWith(err)
			return
		}
		err = e.dispatch(msg)
		if err != nil {
			if msg.Header.Type == Error {
				log.Printf("%s: %v, %s", err, msg.Header,
					readError(msg))
			} else {
				log.Printf("%s: %v", err, msg.Header)
			}
		}
	}
}

// ReceiveAny returns a chanel to receive one message. If the
// connection close, the chanel is closed.
func (e *endPoint) ReceiveAny() (chan *Message, error) {
	filter := func(hdr *Header) (matched bool, keep bool) {
		return true, false
	}
	consumer := make(chan *Message, 1)
	e.MakeHandler(filter, consumer, nil)
	return consumer, nil
}

func (e *endPoint) String() string {
	return e.stream.String()
}

// Pipe returns a set of EndPoint connected to each other.
func Pipe() (EndPoint, EndPoint) {
	a, b := gonet.Pipe()
	return ConnEndPoint(a), ConnEndPoint(b)
}
