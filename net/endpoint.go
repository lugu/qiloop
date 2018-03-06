package net

import (
	"fmt"
	gonet "net"
)

// EndPoint reprensents a network socket capable of sending and
// receiving messages.
type EndPoint struct {
	conn gonet.Conn
}

// DialEndPoint construct an endpoint by contacting a given address.
func DialEndPoint(addr string) (e EndPoint, err error) {
	e.conn, err = gonet.Dial("tcp", addr)
	if err != nil {
		return e, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return e, nil
}

// AcceptedEndPoint construct an endpoint using an accepted
// connection.
func AcceptedEndPoint(c gonet.Conn) EndPoint {
	return EndPoint{
		conn: c,
	}
}

// Send post a message to the other side of the endpoint.
func (e EndPoint) Send(m Message) error {
	return m.Write(e.conn)
}

// Receive wait for a message to be received.
func (e EndPoint) Receive() (m Message, err error) {
	err = m.Read(e.conn)
	return
}
