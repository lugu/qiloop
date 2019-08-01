package bus

import (
	"fmt"
	"sync"

	"github.com/lugu/qiloop/bus/net"
)

// Router dispatch the incomming messages. A Router shall be Activated
// before calling NewService.
type Router struct {
	sync.RWMutex
	services  map[uint32]ServiceReceiver
	namespace Namespace
	session   Session
}

// NewRouter construct a router with the service zero passed.
func NewRouter(authenticator Actor, namespace Namespace, session Session) *Router {
	r := &Router{
		services:  make(map[uint32]ServiceReceiver),
		namespace: namespace,
		session:   session,
	}
	activation := Activation{
		ServiceID: 0,
		ObjectID:  0,
	}
	service, _ := NewService(authenticator, activation)
	r.Add(0, service)
	return r
}

// Terminate terminates all the services.
func (r *Router) Terminate() error {
	r.Lock()
	services := r.services
	r.services = make(map[uint32]ServiceReceiver)
	r.Unlock()

	var ret error
	for serviceID, service := range services {
		err := service.Terminate()
		if err != nil && ret == nil {
			ret = fmt.Errorf("service %d terminate: %s",
				serviceID, err)
		}
	}
	return ret
}

// Add Add a service ot a router. This does not call the Activate()
// method on the service.
func (r *Router) Add(serviceID uint32, s ServiceReceiver) error {
	r.Lock()
	_, ok := r.services[serviceID]
	if ok {
		r.Unlock()
		return fmt.Errorf("service id already used: %d", serviceID)
	}
	r.services[serviceID] = s
	r.Unlock()
	return nil
}

// Remove removes a service from a router.
func (r *Router) Remove(serviceID uint32) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.services[serviceID]; ok {
		delete(r.services, serviceID)
		return nil
	}
	return fmt.Errorf("Router: cannot remove service %d", serviceID)
}

// Receive process a message. The message will be replied if the
// router can not found the destination.
func (r *Router) Receive(m *net.Message, from Channel) error {
	r.RLock()
	s, ok := r.services[m.Header.Service]
	r.RUnlock()
	if ok {
		return s.Receive(m, from)
	}
	return from.SendError(m, ErrServiceNotFound)
}
