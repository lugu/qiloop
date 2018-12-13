package server

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client"
	objproxy "github.com/lugu/qiloop/bus/client/object"
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"sync"
)

type privateNamespace struct {
	sync.Mutex
	reserved  map[string]uint32
	activated map[string]uint32
	next      uint32
}

// PrivateNamespace implements bus.Namespace without relying on a service
// directory. Used for testing purpose.
func PrivateNamespace() bus.Namespace {
	return &privateNamespace{
		reserved:  make(map[string]uint32),
		activated: make(map[string]uint32),
		next:      1,
	}
}

func (ns *privateNamespace) Reserve(name string) (uint32, error) {
	if name == "" {
		return 0, fmt.Errorf("empty string")
	}
	ns.Lock()
	defer ns.Unlock()
	_, ok := ns.reserved[name]
	if ok {
		return 0, fmt.Errorf("service %s already used", name)
	}
	ns.reserved[name] = ns.next
	ns.next++
	return ns.reserved[name], nil
}

func (ns *privateNamespace) Remove(serviceID uint32) error {
	ns.Lock()
	defer ns.Unlock()
	for name, id := range ns.activated {
		if id == serviceID {
			delete(ns.activated, name)
			break
		}
	}
	for name, id := range ns.reserved {
		if id == serviceID {
			delete(ns.reserved, name)
			return nil
		}
	}
	return fmt.Errorf("service %d not in use", serviceID)
}

func (ns *privateNamespace) Enable(serviceID uint32) error {
	ns.Lock()
	defer ns.Unlock()
	for name, id := range ns.reserved {
		if id == serviceID {
			ns.activated[name] = serviceID
			return nil
		}
	}
	return fmt.Errorf("service %d not reserved", serviceID)
}

func (ns *privateNamespace) Resolve(name string) (uint32, error) {
	ns.Lock()
	defer ns.Unlock()
	serviceID, ok := ns.activated[name]
	if !ok {
		_, ok := ns.reserved[name]
		if ok {
			return 0, fmt.Errorf("service %s not enabled", name)
		}
		return 0, fmt.Errorf("unknown service %s", name)
	}
	return serviceID, nil
}

func (ns *privateNamespace) Session(s bus.Server) bus.Session {
	return &localSession{
		namespace: ns,
		server:    s,
	}
}

type localSession struct {
	server    bus.Server
	namespace bus.Namespace
}

func (s *localSession) Client(serviceID uint32) (bus.Client, error) {
	return s.server.Client(), nil
}

func (s *localSession) Proxy(name string, objectID uint32) (bus.Proxy, error) {
	serviceID, err := s.namespace.Resolve(name)
	if err != nil {
		return nil, err
	}
	clt, err := s.Client(serviceID)
	if err != nil {
		return nil, err
	}
	meta, err := bus.MetaObject(clt, serviceID, objectID)
	if err != nil {
		return nil, fmt.Errorf("metaObject (service %d, object %d): %s",
			serviceID, objectID, err)
	}
	return client.NewProxy(clt, meta, serviceID, objectID), nil
}

func (s *localSession) Object(ref object.ObjectReference) (object.Object,
	error) {
	clt, err := s.Client(ref.ServiceID)
	if err != nil {
		return nil, err
	}
	proxy := client.NewProxy(clt, ref.MetaObject, ref.ServiceID,
		ref.ObjectID)
	return &objproxy.ObjectProxy{proxy}, nil
}
func (s *localSession) Destroy() error {
	return nil
}

// Namespace returns a Namespace which connects to a remote service
// directory.
func Namespace(session bus.Session, endpoints []string) (bus.Namespace, error) {
	services := services.Services(session)
	directory, err := services.ServiceDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to connect ServiceDirectory: %s",
			err)
	}
	return &remoteNamespace{
		session:   session,
		Directory: directory,
		EndPoints: endpoints,
		SessionID: "",
		ProcessID: util.ProcessID(),
		MachineID: util.MachineID(),
	}, nil
}

type remoteNamespace struct {
	session   bus.Session
	Directory services.ServiceDirectory
	EndPoints []string
	MachineID string
	ProcessID uint32
	SessionID string
}

func (ns *remoteNamespace) serviceInfo(name string) services.ServiceInfo {
	return services.ServiceInfo{
		Name:      name,
		ServiceId: 0,
		MachineId: ns.MachineID,
		ProcessId: ns.ProcessID,
		Endpoints: ns.EndPoints,
		SessionId: ns.SessionID,
	}
}

func (ns *remoteNamespace) Reserve(name string) (uint32, error) {
	return ns.Directory.RegisterService(ns.serviceInfo(name))
}

func (ns *remoteNamespace) Remove(serviceID uint32) error {
	return ns.Directory.UnregisterService(serviceID)
}

func (ns *remoteNamespace) Enable(serviceID uint32) error {
	return ns.Directory.ServiceReady(serviceID)
}

func (ns *remoteNamespace) Resolve(name string) (uint32, error) {
	info, err := ns.Directory.Service(name)
	if err != nil {
		return 0, err
	}
	return info.ServiceId, nil
}

func (ns *remoteNamespace) Session(s bus.Server) bus.Session {
	return ns.session
}
