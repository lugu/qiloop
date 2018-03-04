package net

import (
	"fmt"
	"qiloop/message"
)

type Client interface {
	Call(service uint32, object uint32, action uint32, payload []byte) ([]byte, error)
}

type BlockingClient struct {
	directory     EndPoint
	nextMessageId uint32
}

func (c *BlockingClient) Call(service uint32, object uint32, action uint32, payload []byte) ([]byte, error) {
	id := c.nextMessageId
	c.nextMessageId += 2
	h := message.NewHeader(message.Call, service, object, action, id)
	m := message.NewMessage(h, payload)
	if err := c.directory.Send(m); err != nil {
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			service, object, action, err)
	}
	response, err := c.directory.Receive()
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply from service %d, object %d, action %d: %s",
			service, object, action, err)
	}
	if response.Header.Id != id {
		return nil, fmt.Errorf("invalid to message id (%d is expected, got %d)",
			id, response.Header.Id)
	}
	return response.Payload, nil
}

func NewClient(endpoint string) (Client, error) {
	directory, err := DialEndPoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("client failed to connect %s: %s", endpoint, err)
	}
	return &BlockingClient{directory, 1}, nil
}
