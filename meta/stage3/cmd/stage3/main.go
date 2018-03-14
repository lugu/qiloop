package main

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/object"
	"github.com/lugu/qiloop/value"
	"io"
	"log"
	"os"
	"strings"
)

type dummyClient struct {
	conn          net.EndPoint
	nextMessageID uint32
}

func (c dummyClient) Call(serviceID uint32, objectID uint32, actionID uint32, payload []byte) ([]byte, error) {
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

type dummyProxy struct {
	client          dummyClient
	serviceID       uint32
	objectID        uint32
	defaultActionID uint32
}

// Call ignores the action and call a pre-defined actionID.
func (p dummyProxy) CallID(actionID uint32, payload []byte) ([]byte, error) {
	return p.client.Call(p.serviceID, p.objectID, actionID, payload)
}

func (p dummyProxy) Call(action string, payload []byte) ([]byte, error) {
	return p.client.Call(p.serviceID, p.objectID, p.defaultActionID, payload)
}

// ServiceID returns the service identifier.
func (p dummyProxy) ServiceID() uint32 {
	return p.serviceID
}

// ObjectID returns the object identifier within the service.
func (p dummyProxy) ObjectID() uint32 {
	return p.objectID
}

// SignalStreamID is not implemented. Does nothing in stage3.
func (p dummyProxy) SignalStreamID(signal uint32, cancel chan int) (chan []byte, error) {
	return nil, fmt.Errorf("SignalStreamID not available during stage 3")
}

// SignalStream is not implemented. Does nothing in stage3.
func (p dummyProxy) SignalStream(signal string, cancel chan int) (chan []byte, error) {
	return nil, fmt.Errorf("SignalStream not available during stage 3")
}

// MethodUid is not implemented. Does nothing in stage3.
func (p dummyProxy) MethodUid(name string) (uint32, error) {
	return 0, fmt.Errorf("MethodUid not available during stage 3")
}

// SignalUid is not implemented. Does nothing in stage3.
func (p dummyProxy) SignalUid(name string) (uint32, error) {
	return 0, fmt.Errorf("SignalUid not available during stage 3")
}

type dummySession struct {
	endpoint         net.EndPoint
	defaultSerivceID uint32
	defaultObjectID  uint32
	defaultActionID  uint32
}

// Proxy ignores the service name and use a pre-defined serviceID and
// objectID.
func (s dummySession) Proxy(name string, objectID uint32) (bus.Proxy, error) {
	return dummyProxy{
		client: dummyClient{
			conn:          s.endpoint,
			nextMessageID: 3,
		},
		serviceID:       s.defaultSerivceID,
		objectID:        s.defaultObjectID,
		defaultActionID: s.defaultActionID,
	}, nil
}

func (d dummySession) Object(ref object.ObjectReference) (o object.Object, err error) {
	return o, fmt.Errorf("Not yet implemented")
}

func NewSession(conn net.EndPoint, serviceID, objectID, actionID uint32) bus.Session {

	sess0 := dummySession{conn, 0, 0, 8}
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
	return dummySession{
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

func main() {
	var outputImplementation io.Writer
	var outputInterface io.Writer

	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		outputInterface = file
		defer file.Close()
	} else {
		outputInterface = os.Stdout
	}

	if len(os.Args) > 2 {
		filename := os.Args[2]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		outputImplementation = file
		defer file.Close()
	} else {
		outputImplementation = os.Stdout
	}

	addr := ":9559"
	// directoryServiceID := 1
	// directoryObjectID := 1
	dir, err := NewServiceDirectory(addr, 1, 1, 101)
	if err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}

	serviceInfoList, err := dir.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	objects := make([]object.MetaObject, 0)
	objects = append(objects, object.MetaService0)
	objects = append(objects, object.ObjectMetaObject)

	for i, s := range serviceInfoList {

		// FIXME: change me to 100
		if i > 1 {
			continue
		}

		addr := strings.TrimPrefix(s.Endpoints[0], "tcp://")
		obj, err := NewObject(addr, s.ServiceId, 1, 2)
		if err != nil {
			log.Printf("failed to create servinceof %s: %s", s.Name, err)
			continue
		}
		meta, err := obj.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
			continue
		}
		meta.Description = s.Name
		objects = append(objects, meta)
	}
	proxy.GenerateInterfaces(objects, "services", outputInterface)
	proxy.GenerateProxys(objects, "services", outputImplementation)
}
