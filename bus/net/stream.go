package net

import (
	"context"
	"fmt"
	"io"
	gonet "net"
	"os"
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

type pipeStream struct {
	r   *os.File
	w   *os.File
	ctx context.Context
}

func (p *pipeStream) Read(d []byte) (int, error) {
	return p.r.Read(d)
}

func (p *pipeStream) Write(d []byte) (int, error) {
	return p.w.Write(d)
}

func (p *pipeStream) Close() error {
	p.r.Close()
	p.w.Close()
	return nil
}

func (p *pipeStream) String() string {
	return fmt.Sprintf("pipe://%d:%d", p.r.Fd(), p.w.Fd())
}

func (p *pipeStream) Context() context.Context {
	return p.ctx
}

// PipeStream returns a Stream based on the pipe:// protocol
func PipeStream(r, w *os.File) Stream {
	return &pipeStream{
		r:   r,
		w:   w,
		ctx: context.TODO(),
	}
}

