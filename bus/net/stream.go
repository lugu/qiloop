package net

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	gonet "net"
	"sync"

	quic "github.com/lucas-clemente/quic-go"
	"github.com/lugu/qiloop/bus/net/cert"
)

// Stream represents a network connection. Stream abstracts
// connections to allow for various transports.
type Stream interface {
	io.Reader
	io.Writer
	io.Closer
	fmt.Stringer
	Context() context.Context
}

type connStream struct {
	gonet.Conn
	ctx context.Context
}

func (c connStream) String() string {
	return c.RemoteAddr().Network() + "://" +
		c.RemoteAddr().String()
}

func (c connStream) Context() context.Context {
	return c.ctx
}

// ConnStream construct a Stream from a connection.
func ConnStream(conn gonet.Conn) Stream {
	return connStream{
		conn,
		context.TODO(),
	}
}

type quicStream struct {
	quic.Stream
	// quic.Stream does not permit to call Close while Writing.
	// Refer to go-quic documentation.
	sync.RWMutex
}

func newQuicStream(s quic.Stream) Stream {
	return &quicStream{
		s,
		sync.RWMutex{},
	}
}

func (s *quicStream) Close() error {
	s.Lock()
	defer s.Unlock()
	return s.Stream.Close()
}

func (s *quicStream) Read(p []byte) (int, error) {
	return s.Stream.Read(p)
}

func (s *quicStream) Write(p []byte) (int, error) {
	s.RLock()
	defer s.RUnlock()
	return s.Stream.Write(p)
}

func (s *quicStream) String() string {
	return fmt.Sprintf("StreamID %d", s.StreamID())
}

type quicListener struct {
	l       quic.Listener
	streams chan Stream
	closer  chan struct{}
	errors  chan error
}

func newQuicListener(l quic.Listener, ctx context.Context) (Listener, error) {
	q := &quicListener{
		l:       l,
		streams: make(chan Stream),
		closer:  make(chan struct{}),
		errors:  make(chan error),
	}
	go q.bg(ctx)
	return q, nil
}

func (q quicListener) Accept() (Stream, error) {
	select {
	case <-q.closer:
		return nil, io.EOF
	case err := <-q.errors:
		return nil, err
	case stream := <-q.streams:
		return stream, nil
	}
}

func (q quicListener) Close() error {
	err := q.l.Close()
	close(q.closer)
	return err
}

func (q quicListener) bg(ctx context.Context) {
	for {
		sess, err := q.l.Accept(ctx)
		if err != nil {
			q.errors <- err
			return
		}
		q.handleSession(sess, ctx)
	}
}

func (q quicListener) handleSession(sess quic.Session, ctx context.Context) {
	cancel := make(chan struct{})
	go func() {
		select {
		case <-q.closer: // close the sesion on demand
		case <-cancel: // close the sesion on error
		}
		sess.Close()
	}()
	// send stream of streams into streams
	go func() {
		for {
			stream, err := sess.AcceptStream(ctx)
			if err != nil {
				log.Printf("Session error: %s <-> %s : %s",
					sess.LocalAddr().String(),
					sess.RemoteAddr().String(),
					err)
				close(cancel)
				return
			}
			select {
			case <-q.closer:
				return
			case q.streams <- newQuicStream(stream):
			}
		}
	}()
}

func listenQUIC(addr string) (Listener, error) {
	var err1, err2 error
	cer, err1 := cert.Certificate()
	if err1 != nil {
		log.Printf("Failed to read x509 certificate: %s", err1)
		cer, err2 = cert.GenerateCertificate()
		if err2 != nil {
			log.Printf("Failed to create x509 certificate: %s", err2)
			return nil, fmt.Errorf("no certificate available (%s, %s)",
				err1, err2)
		}
	}

	conf := &tls.Config{Certificates: []tls.Certificate{cer}}

	listener, err := quic.ListenAddr(addr, conf, nil)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(context.TODO(), ListenAddress, addr)
	return newQuicListener(listener, ctx)
}
