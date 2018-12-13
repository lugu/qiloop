package session

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client"
	objproxy "github.com/lugu/qiloop/bus/client/object"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/net"
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
	Directory        services.ServiceDirectory
	cancel           chan int
	added            chan struct {
		P0 uint32
		P1 string
	}
	removed chan struct {
		P0 uint32
		P1 string
	}
}

func newObject(info services.ServiceInfo, ref object.ObjectReference) (object.Object, error) {
	endpoint, err := client.SelectEndPoint(info.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("object connection error (%s): %s",
			info.Name, err)
	}
	proxy := client.NewProxy(client.NewClient(endpoint), ref.MetaObject,
		ref.ServiceID, ref.ObjectID)
	return &objproxy.ObjectProxy{proxy}, nil
}

func newService(info services.ServiceInfo, objectID uint32) (p bus.Proxy, err error) {
	endpoint, err := client.SelectEndPoint(info.Endpoints)
	if err != nil {
		return nil, fmt.Errorf("service connection error (%s): %s", info.Name, err)
	}
	c := client.NewClient(endpoint)
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
func (s *Session) Object(ref object.ObjectReference) (o object.Object, err error) {
	info, err := s.findServiceID(ref.ServiceID)
	if err != nil {
		return o, err
	}
	return newObject(info, ref)
}

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(c bus.Client, serviceID, objectID uint32) (p bus.Proxy, err error) {
	meta, err := bus.MetaObject(c, serviceID, objectID)
	if err != nil {
		return p, fmt.Errorf("Can not reach metaObject: %s", err)
	}
	return client.NewProxy(c, meta, serviceID, objectID), nil
}

// BindSession returns a session based on a previously established
// connection.
func BindSession(c bus.Client) (*Session, error) {
	proxy, err := metaProxy(c, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get directory meta object: %s", err)
	}
	s := new(Session)
	s.Directory = &services.ServiceDirectoryProxy{objproxy.ObjectProxy{proxy}}

	s.serviceList, err = s.Directory.Services()
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %s", err)
	}
	s.cancel = make(chan int)
	s.removed, err = s.Directory.SignalServiceRemoved(s.cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe remove signal: %s", err)
	}
	s.added, err = s.Directory.SignalServiceAdded(s.cancel)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe added signal: %s", err)
	}
	return s, nil
}

// NewSession connects an address and return a new session.
func NewSession(addr string) (bus.Session, error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to contact %s: %s", addr, err)
	}
	if err = client.Authenticate(endpoint); err != nil {
		endpoint.Close()
		return nil, fmt.Errorf("authentication failed: %s", err)
	}
	c := client.NewClient(endpoint)

	sess, err := BindSession(c)
	if err != nil {
		endpoint.Close()
		return nil, err
	}
	return sess, nil
}

func (s *Session) updateServiceList() {
	var err error
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	s.serviceList, err = s.Directory.Services()
	if err != nil {
		log.Printf("error: failed to update service directory list: %s", err)
		log.Printf("error: closing session.")
		if err := s.Destroy(); err != nil {
			log.Printf("error: session destruction: %s", err)
		}
	}
}

// Destroy close the session.
func (s *Session) Destroy() error {
	// cancel both add and remove services
	s.cancel <- 1
	s.cancel <- 1
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
