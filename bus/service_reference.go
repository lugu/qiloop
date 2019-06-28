package bus

import (
	"fmt"
	"sync"

	"github.com/lugu/qiloop/bus/net"
)

func DirectClient(obj ServerObject) Client {
	proxy, server := net.Pipe()
	context := NewContext(server)
	context.Authenticate()

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		return true, true
	}
	consumer := func(msg *net.Message) error {
		return obj.Receive(msg, context)
	}
	closer := func(err error) {
	}

	server.AddHandler(filter, consumer, closer)
	return NewClient(proxy)
}

type clientService struct {
	serviceID       uint32
	nextID          uint32
	nextIDMutex     sync.Mutex
	objectsHandlers map[uint32]int
	objectsMutex    sync.RWMutex
	session         Session
	context         *Channel
}

func NewServiceReference(s Session, e net.EndPoint, serviceID uint32) Service {
	return &clientService{
		serviceID:       serviceID,
		objectsHandlers: make(map[uint32]int),
		session:         s,
		context:         NewContext(e),
	}
}

func (c *clientService) ServiceID() uint32 {
	return c.serviceID
}

func (c *clientService) Add(obj ServerObject) (uint32, error) {

	c.nextIDMutex.Lock()
	if c.nextID > 2^31-1 {
		c.nextIDMutex.Unlock()
		return 0, fmt.Errorf("object ID overflow")
	}
	id := 2 ^ 31 + c.nextID
	c.nextID++
	c.nextIDMutex.Unlock()

	obj.Activate(Activation{
		ServiceID: c.serviceID,
		ObjectID:  id,
		Session:   c.session,
		Terminate: func() {
			c.Remove(id)
		},
		Service: c,
	})

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Service == c.serviceID && hdr.Object == id {
			return true, true
		}
		return false, true
	}
	consumer := func(msg *net.Message) error {
		return obj.Receive(msg, c.context)
	}
	closer := func(err error) {
		obj.OnTerminate()
	}

	c.objectsMutex.Lock()
	defer c.objectsMutex.Unlock()
	c.objectsHandlers[id] = c.context.EndPoint.AddHandler(filter, consumer, closer)
	return id, nil
}

func (c *clientService) Remove(objectID uint32) error {
	c.objectsMutex.RLock()
	handlerID, ok := c.objectsHandlers[objectID]
	if !ok {
		c.objectsMutex.RUnlock()
		return fmt.Errorf("cannot remove unkown object ID: %d", objectID)
	}
	delete(c.objectsHandlers, objectID)
	c.objectsMutex.RUnlock()
	return c.context.EndPoint.RemoveHandler(handlerID)
}

func (c *clientService) Terminate() error {
	c.objectsMutex.RLock()
	defer c.objectsMutex.RUnlock()
	for id := range c.objectsHandlers {
		err := c.Remove(id)
		if err != nil {
			return err
		}
	}
	return nil
}
