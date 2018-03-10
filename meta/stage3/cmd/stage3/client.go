package main

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/object"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
)

type directoryClient struct {
	conn          net.EndPoint
	nextMessageID uint32
}

func (c directoryClient) Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error) {
	id := c.nextMessageID
	c.nextMessageID += 2
	h := net.NewHeader(net.Call, serviceID, objectID, actionID, id)
	m := net.NewMessage(h, payload)
	if err := c.conn.Send(m); err != nil {
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}
	response, err := c.conn.Receive()
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

type directoryProxy struct {
	client          directoryClient
	serviceID       uint32
	objectID        uint32
	defaultActionID uint32
}

// Call ignores the action and call a pre-defined actionID.
func (p directoryProxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.serviceID, p.objectID, actionID, payload)
}

func (p directoryProxy) Call(action string, payload []byte) ([]byte, error) {
	return p.client.Call(p.serviceID, p.objectID, p.defaultActionID, payload)
}

// ServiceID returns the service identifier.
func (p directoryProxy) ServiceID() uint32 {
	return p.serviceID
}

// ObjectID returns the object identifier within the service.
func (p directoryProxy) ObjectID() uint32 {
	return p.objectID
}

type directorySession struct {
	endpoint         net.EndPoint
	defaultSerivceID uint32
	defaultObjectID  uint32
	defaultActionID  uint32
}

// Proxy ignores the service name and use a pre-defined serviceID and
// objectID.
func (s directorySession) Proxy(name string, objectID uint32) (session.Proxy, error) {
	return directoryProxy{
		client: directoryClient{
			conn:          s.endpoint,
			nextMessageID: 3,
		},
		serviceID:       s.defaultSerivceID,
		objectID:        s.defaultObjectID,
		defaultActionID: s.defaultActionID,
	}, nil
}

func (d directorySession) Object(ref object.ObjectReference) (o object.Object, err error) {
	return o, fmt.Errorf("Not yet implemented")
}
