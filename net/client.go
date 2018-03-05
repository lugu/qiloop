package net

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/message"
	"github.com/lugu/qiloop/value"
)

// Client represents a client connection to a service.
type Client interface {
	Call(service uint32, object uint32, action uint32, payload []byte) ([]byte, error)
}

type blockingClient struct {
	directory     EndPoint
	nextMessageID uint32
}

func (c *blockingClient) Call(service uint32, object uint32, action uint32, payload []byte) ([]byte, error) {
	id := c.nextMessageID
	c.nextMessageID += 2
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
	if response.Header.ID != id {
		return nil, fmt.Errorf("invalid to message id (%d is expected, got %d)",
			id, response.Header.ID)
	}
	if response.Header.Type == message.Error {
		message, err := value.NewValue(bytes.NewBuffer(response.Payload))
		if err != nil {
			return nil, fmt.Errorf("Error: failed to parse error message: %s", string(response.Payload))
		}
		return nil, fmt.Errorf("Error: %s", message)
	}
	return response.Payload, nil
}

// NewClient returns a Client connected to the specified endpoint.
func NewClient(endpoint string) (Client, error) {
	directory, err := DialEndPoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("client failed to connect %s: %s", endpoint, err)
	}
	return &blockingClient{directory, 1}, nil
}
