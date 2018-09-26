package basic

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"log"
)

type basicClient struct {
	conn          net.EndPoint
	nextMessageID uint32
}

func (c basicClient) Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error) {
	id := c.nextMessageID
	c.nextMessageID += 2
	h := net.NewHeader(net.Call, serviceID, objectID, actionID, id)
	m := net.NewMessage(h, payload)

	received := make(chan *net.Message)
	go func() {
		msg, err := c.conn.ReceiveAny()
		if err != nil {
			fmt.Errorf("failed to receive net. %s", err)
		}
		received <- msg
	}()

	if err := c.conn.Send(m); err != nil {
		return nil, fmt.Errorf("failed to call service %d, object %d, action %d: %s",
			serviceID, objectID, actionID, err)
	}

	response := <-received

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

type basicProxy struct {
	client          basicClient
	serviceID       uint32
	objectID        uint32
	defaultActionID uint32
}

// Call ignores the action and call a pre-defined actionID.
func (p basicProxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.serviceID, p.objectID, actionID, payload)
}

func (p basicProxy) Call(action string, payload []byte) ([]byte, error) {
	return p.client.Call(p.serviceID, p.objectID, p.defaultActionID, payload)
}

// ServiceID returns the service identifier.
func (p basicProxy) ServiceID() uint32 {
	return p.serviceID
}

// ObjectID returns the object identifier within the service.
func (p basicProxy) ObjectID() uint32 {
	return p.objectID
}

// SubscribeID is not implemented.
func (p basicProxy) SubscribeID(signal uint32, cancel chan int) (chan []byte, error) {
	return nil, fmt.Errorf("SubscribeID not available during stage 3")
}

// Subscribe is not implemented.
func (p basicProxy) SubscribeSignal(signal string, cancel chan int) (chan []byte, error) {
	return nil, fmt.Errorf("Subscribe not available during stage 3")
}

// MethodUid is not implemented.
func (p basicProxy) MethodUid(name string) (uint32, error) {
	return 0, fmt.Errorf("MethodUid not available during stage 3")
}

// SignalUid is not implemented.
func (p basicProxy) SignalUid(name string) (uint32, error) {
	return 0, fmt.Errorf("SignalUid not available during stage 3")
}

// Disconnect is not implemented.
func (p basicProxy) Disconnect() error {
	return fmt.Errorf("SignalUid not available during stage 3")
}

type basicSession struct {
	endpoint         net.EndPoint
	defaultSerivceID uint32
	defaultObjectID  uint32
	defaultActionID  uint32
}

// Proxy ignores the service name and use a pre-defined serviceID and
// objectID.
func (s basicSession) Proxy(name string, objectID uint32) (bus.Proxy, error) {
	return basicProxy{
		client: basicClient{
			conn:          s.endpoint,
			nextMessageID: 3,
		},
		serviceID:       s.defaultSerivceID,
		objectID:        s.defaultObjectID,
		defaultActionID: s.defaultActionID,
	}, nil
}

func (d basicSession) Object(ref object.ObjectReference) (o object.Object, err error) {
	return o, fmt.Errorf("Not yet implemented")
}

func (d basicSession) Destroy() error {
	return d.endpoint.Close()
}

func NewSession(conn net.EndPoint, serviceID, objectID, actionID uint32) bus.Session {

	sess0 := basicSession{conn, 0, 0, 8}
	service0, err := stage2.NewServer(sess0, 0)
	if err != nil {
		log.Fatalf("failed to create proxy: %s", err)
	}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	_, err = service0.Authenticate(permissions)
	if err != nil {
		log.Fatalf("failed to authenticate: %s", err)
	}
	return basicSession{
		conn,
		serviceID,
		objectID,
		actionID,
	}
}

func NewObject(addr string, serviceID, objectID, actionID uint32) (d stage2.Object, err error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return d, fmt.Errorf("failed to connect: %s", err)
	}
	sess := NewSession(endpoint, serviceID, objectID, actionID)
	if err != nil {
		return d, fmt.Errorf("failed to create session: %s", err)
	}

	return stage2.NewObject(sess, 1)
}

func NewServiceDirectory(addr string, serviceID, objectID, actionID uint32) (d stage2.ServiceDirectory, err error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return d, fmt.Errorf("failed to connect: %s", err)
	}
	sess := NewSession(endpoint, serviceID, objectID, actionID)
	if err != nil {
		return d, fmt.Errorf("failed to create session: %s", err)
	}

	return stage2.NewServiceDirectory(sess, 1)
}
