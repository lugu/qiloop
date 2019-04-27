package net

import (
	"context"
	"io"
	gonet "net"
)

// Stream represents a network connection. Stream abstracts
// connections to allow for various transports.
type Stream interface {
	io.Reader
	io.Writer
	io.Closer
	Context() context.Context
	String() string
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

func ConnStream(conn gonet.Conn) Stream {
	return connStream{
		conn,
		context.TODO(),
	}
}
