package bus

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/lugu/qiloop/bus/net"
)

func objectTerminator(service Service, objectID uint32) func() {
	return func() {
		service.Remove(objectID)
	}
}

func serviceTerminator(router *Router, serviceID uint32) func() {
	return func() {
		router.Remove(serviceID)
		router.namespace.Remove(serviceID)
	}
}

func serviceActivation(router *Router, session Session, serviceID uint32) Activation {
	return Activation{
		ServiceID: serviceID,
		ObjectID:  1,
		Session:   session,
		Terminate: serviceTerminator(router, serviceID),
		Service:   nil,
	}
}

func objectActivation(service *serviceImpl, session Session, serviceID, objectID uint32) Activation {
	return Activation{
		ServiceID: serviceID,
		ObjectID:  objectID,
		Session:   session,
		Terminate: objectTerminator(service, objectID),
		Service:   service,
	}
}

type pendingObject struct{}

func (p pendingObject) Receive(m *net.Message, from Channel) error {
	return ErrObjectNotFound
}

func (p pendingObject) Activate(activation Activation) error {
	panic("can not activate pending object")
}

func (p pendingObject) OnTerminate() {
}

// serviceImpl implements Service and Receiver. It allows a service to
// manage the object within its domain.
type serviceImpl struct {
	sync.RWMutex
	objects   map[uint32]Actor
	terminate func()
	session   Session
	serviceID uint32
}

// NewService returns a service basedon the given object.
func NewService(o Actor, activation Activation) (ServiceReceiver, error) {
	s := &serviceImpl{
		objects: map[uint32]Actor{
			activation.ObjectID: o,
		},
	}
	return s, s.activate(activation)
}

func (s *serviceImpl) ServiceID() uint32 {
	return s.serviceID
}

// Add is used to add an object to a service domain.
func (s *serviceImpl) Add(obj Actor) (index uint32, err error) {
	// assign the first object to the index 0. following objects will
	// be assigned random values.
	s.Lock()
	if _, ok := s.objects[1]; ok {
		index = (rand.Uint32() << 1) >> 1
		if _, ok = s.objects[index]; ok {
			s.Unlock()
			return s.Add(obj)
		}
	}
	if s.session == nil { // service not yet activated
		s.objects[index] = obj
		s.Unlock()
		return
	}
	s.objects[index] = pendingObject{}
	s.Unlock()

	a := objectActivation(s, s.session, s.serviceID, index)
	err = obj.Activate(a)

	s.Lock()
	if err != nil {
		s.objects[index] = nil
	} else {
		s.objects[index] = obj
	}
	s.Unlock()
	return
}

// activate informs the service it will become active and shall be
// ready to handle requests. activation.Service is nil.
func (s *serviceImpl) activate(activation Activation) error {
	var wait sync.WaitGroup
	s.terminate = activation.Terminate
	s.session = activation.Session
	s.serviceID = activation.ServiceID
	s.Lock()
	ret := make(chan error, len(s.objects))
	wait.Add(len(s.objects))
	for objectID, obj := range s.objects {
		go func(obj Actor, objectID uint32) {
			objActivation := objectActivation(s, activation.Session,
				activation.ServiceID, objectID)
			err := obj.Activate(objActivation)
			if err != nil {
				ret <- err
			}
			wait.Done()
		}(obj, objectID)
	}
	s.Unlock()
	wait.Wait()
	close(ret)
	for err := range ret {
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove removes an object from the service domain.
func (s *serviceImpl) Remove(objectID uint32) error {
	s.Lock()
	if obj, ok := s.objects[objectID]; ok {
		delete(s.objects, objectID)
		s.Unlock()
		obj.OnTerminate()
		return nil
	}
	s.Unlock()
	return fmt.Errorf("cannot remove object %d", objectID)
}

// Receive forwards the message to the appropriate object.
func (s *serviceImpl) Receive(m *net.Message, from Channel) error {
	s.RLock()
	o, ok := s.objects[m.Header.Object]
	s.RUnlock()
	if ok {
		return o.Receive(m, from)
	}
	return from.SendError(m, ErrObjectNotFound)
}

// Terminate calls OnTerminate on all its objects.
func (s *serviceImpl) Terminate() error {
	s.RLock()
	defer s.RUnlock()

	for _, obj := range s.objects {
		obj.OnTerminate()
	}
	if s.terminate != nil {
		s.terminate()
	}
	return nil
}
