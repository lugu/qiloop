package main

import (
	"bytes"
	"fmt"
	server "github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
)

func manualProxy(e net.EndPoint, service, object uint32) session.Proxy {
	return session.NewProxy(&blockingClient{e, 3}, service, object)
}

func authenticate(e net.EndPoint) error {
	server0 := server.Server{manualProxy(e, 0, 0)}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	if _, err := server0.Authenticate(permissions); err != nil {
		fmt.Errorf("authentication failed: %s", err)
	}
	return nil
}

type blockingClient struct {
	conn          net.EndPoint
	nextMessageID uint32
}

func (c *blockingClient) Call(service uint32, object uint32, action uint32, payload []byte) ([]byte, error) {
	id := c.nextMessageID
	c.nextMessageID += 2
	h := net.NewHeader(net.Call, service, object, action, id)
	m := net.NewMessage(h, payload)
	if err := c.conn.Send(m); err != nil {
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			service, object, action, err)
	}
	response, err := c.conn.Receive()
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply from service %d, object %d, action %d: %s",
			service, object, action, err)
	}
	if response.Header.ID != id {
		return nil, fmt.Errorf("invalid to message id (%d is expected, got %d)",
			id, response.Header.ID)
	}
	if response.Header.Type == net.Error {
		message, err := value.NewValue(bytes.NewBuffer(response.Payload))
		if err != nil {
			return nil, fmt.Errorf("Error: failed to parse error message: %s", string(response.Payload))
		}
		return nil, fmt.Errorf("Error: %s", message)
	}
	return response.Payload, nil
}
