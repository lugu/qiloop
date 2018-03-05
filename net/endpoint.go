package net

import (
	"fmt"
	gonet "net"
	"github.com/lugu/qiloop/message"
)

type EndPoint struct {
	conn gonet.Conn
}

func DialEndPoint(addr string) (e EndPoint, err error) {
	e.conn, err = gonet.Dial("tcp", addr)
	if err != nil {
		return e, fmt.Errorf("failed to connect %s: %s", addr, err)
	}
	return e, nil
}

func AcceptedEndPoint(c gonet.Conn) EndPoint {
	return EndPoint{
		conn: c,
	}
}

func (e EndPoint) Send(m message.Message) error {
	return m.Write(e.conn)
}

func (e EndPoint) Receive() (m message.Message, err error) {
	err = m.Read(e.conn)
	return
}
