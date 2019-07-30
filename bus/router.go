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
	services  map[uint32]*serviceImpl
	namespace Namespace
	session   Session // nil until activation
}

// NewRouter construct a router with the service zero passed.
func NewRouter(authenticator Actor, namespace Namespace) *Router {
	return &Router{
		services: map[uint32]*serviceImpl{
			0: {
				objects: map[uint32]Actor{
					0: authenticator,
				},
				boxes: map[uint32]MailBox{
					0: NewMailBox(authenticator),
				},
			},
		},
		namespace: namespace,
		session:   nil,
	}
}

// Activate calls the Activate method on all the services. Only after
// the router can process messages.
func (r *Router) Activate(session Session) error {
	r.Lock()
	if r.session != nil {
		r.Unlock()
		return fmt.Errorf("router already activated")
	}
	r.session = session
	r.Unlock()
	for serviceID, service := range r.services {
		activation := serviceActivation(r, session, serviceID)
		err := service.Activate(activation)
		if err != nil {
			return err
		}
	}
	return nil
}

// Terminate terminates all the services.
func (r *Router) Terminate() error {
	r.Lock()
	services := r.services
	r.services = make(map[uint32]*serviceImpl)
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
func (r *Router) Add(serviceID uint32, s *serviceImpl) error {
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

// Dispatch process a message. The message will be replied if the
// router can not found the destination.
func (r *Router) Dispatch(m *net.Message, from Channel) error {
	r.RLock()
	s, ok := r.services[m.Header.Service]
	r.RUnlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return from.SendError(m, ErrServiceNotFound)
}
