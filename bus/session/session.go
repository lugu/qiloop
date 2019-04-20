package session

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/type/object"
	"log"
	"sync"
)

// Session implements the Session interface. It is an
// implementation of Session. It does not update the list of services
// and returns clients.
type Session struct {
	serviceList      []services.ServiceInfo
	serviceListMutex sync.Mutex
	Directory        services.ServiceDirectoryProxy
	cancel           func()
	added            chan services.ServiceAdded
	removed          chan services.ServiceRemoved
}

func newObject(info services.ServiceInfo, ref object.ObjectReference) (bus.ObjectProxy, error) {
	endpoint, err := bus.SelectEndPoint(info.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("object connection error (%s): %s",
			info.Name, err)
	}
	proxy := bus.NewProxy(bus.NewClient(endpoint), ref.MetaObject,
		ref.ServiceID, ref.ObjectID)
	return bus.MakeObject(proxy), nil
}

func newService(info services.ServiceInfo, objectID uint32) (p bus.Proxy, err error) {
	endpoint, err := bus.SelectEndPoint(info.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("service connection error (%s): %s", info.Name, err)
	}
	c := bus.NewClient(endpoint)
	proxy, err := metaProxy(c, info.ServiceId, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service meta object (%s): %s", info.Name, err)
	}
	return proxy, nil
}

func (s *Session) findServiceName(name string) (i services.ServiceInfo, err error) {
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	for _, service := range s.serviceList {
		if service.Name == name {
			return service, nil
		}
	}
	return i, fmt.Errorf("Service not found: %s", name)
}

func (s *Session) findServiceID(uid uint32) (i services.ServiceInfo, err error) {
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	for _, service := range s.serviceList {
		if service.ServiceId == uid {
			return service, nil
		}
	}
	return i, fmt.Errorf("Service ID not found: %d", uid)
}

// Proxy resolve the service name and returns a proxy to it.
func (s *Session) Proxy(name string, objectID uint32) (p bus.Proxy, err error) {
	info, err := s.findServiceName(name)
	if err != nil {
		return p, err
	}
	return newService(info, objectID)
}

// Object returns a reference to ref.
func (s *Session) Object(ref object.ObjectReference) (o bus.Proxy, err error) {
	info, err := s.findServiceID(ref.ServiceID)
	if err != nil {
		return o, err
	}
	return newObject(info, ref)
}

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(c bus.Client, serviceID, objectID uint32) (p bus.Proxy, err error) {
	meta, err := bus.GetMetaObject(c, serviceID, objectID)
	if err != nil {
		return p, fmt.Errorf("Can not reach metaObject: %s", err)
	}
	return bus.NewProxy(c, meta, serviceID, objectID), nil
}

// NewSession connects an address and return a new session.
func NewSession(addr string) (bus.Session, error) {

	s := new(Session)
	// Manually create a serviceList with just the ServiceInfo
	// needed to contact ServiceDirectory.
	s.serviceList = []services.ServiceInfo{
		services.ServiceInfo{
			Name:      "ServiceDirectory",
			ServiceId: 1,
			Endpoints: []string{
				addr,
			},
		},
	}
	var err error
	s.Directory, err = services.Services(s).ServiceDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to contact server: %s", err)
	}

	s.serviceList, err = s.Directory.Services()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %s", err)
	}
	var cancelRemoved, cancelAdded func()
	cancelRemoved, s.removed, err = s.Directory.SubscribeServiceRemoved()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe remove signal: %s", err)
	}
	cancelAdded, s.added, err = s.Directory.SubscribeServiceAdded()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe added signal: %s", err)
	}
	s.cancel = func() {
		cancelRemoved()
		cancelAdded()
	}
	return s, nil
}

func (s *Session) updateServiceList() {
	services, err := s.Directory.Services()
	if err != nil {
		log.Printf("error: failed to update service directory list: %s", err)
		log.Printf("error: closing session.")
		if err := s.Destroy(); err != nil {
			log.Printf("error: session destruction: %s", err)
		}
	}
	s.serviceListMutex.Lock()
	s.serviceList = services
	s.serviceListMutex.Unlock()
}

// Destroy close the session.
func (s *Session) Destroy() error {
	s.cancel()
	return s.Directory.Disconnect()
}

func (s *Session) updateLoop() {
	for {
		select {
		case _, ok := <-s.removed:
			if !ok {
				return
			}
			s.updateServiceList()
		case _, ok := <-s.added:
			if !ok {
				return
			}
			s.updateServiceList()
		}
	}
}
