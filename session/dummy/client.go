package dummy

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
)

type blockingClient struct {
	directory     net.EndPoint
	nextMessageID uint32
}

func (c *blockingClient) Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error) {
	id := c.nextMessageID
	c.nextMessageID += 2
	h := net.NewHeader(net.Call, serviceID, objectID, actionID, id)
	m := net.NewMessage(h, payload)
	if err := c.directory.Send(m); err != nil {
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}
	response, err := c.directory.Receive()
	if err != nil {
		return nil, fmt.Errorf("failed to receive reply from service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
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

func NewClient(endpoint string) (session.Client, error) {
	directory, err := net.DialEndPoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("client failed to connect %s: %s", endpoint, err)
	}
	return &blockingClient{directory, 1}, nil
}
