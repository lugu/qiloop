package directory

import (
	"fmt"
	"io"
	"sort"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
)

// serviceDirectory implements ServiceDirectoryImplementor
type serviceDirectory struct {
	staging  map[uint32]ServiceInfo
	services map[uint32]ServiceInfo
	lastID   uint32
	signal   ServiceDirectorySignalHelper
}

// serviceDirectoryImpl returns an implementation of ServiceDirectory
func serviceDirectoryImpl() *serviceDirectory {
	return &serviceDirectory{
		staging:  make(map[uint32]ServiceInfo),
		services: make(map[uint32]ServiceInfo),
		lastID:   0,
	}
}

func (s *serviceDirectory) Activate(activation bus.Activation,
	helper ServiceDirectorySignalHelper) error {
	s.signal = helper
	return nil
}

func (s *serviceDirectory) OnTerminate() {
}

func checkServiceInfo(i ServiceInfo) error {
	if i.Name == "" {
		return fmt.Errorf("empty name not allowed")
	}
	if i.MachineId == "" {
		return fmt.Errorf("empty machine id not allowed")
	}
	if i.ProcessId == 0 {
		return fmt.Errorf("process id zero not allowed")
	}
	if len(i.Endpoints) == 0 {
		return fmt.Errorf("missing end point (%s)", i.Name)
	}
	for _, e := range i.Endpoints {
		if e == "" {
			return fmt.Errorf("empty endpoint not allowed")
		}
	}
	return nil
}

func (s *serviceDirectory) info(serviceID uint32) (ServiceInfo, error) {
	info, ok := s.services[serviceID]
	if !ok {
		return info, fmt.Errorf("service %d not found", serviceID)
	}
	return info, nil
}

func (s *serviceDirectory) Service(service string) (info ServiceInfo, err error) {
	for _, info = range s.services {
		if info.Name == service {
			return info, nil
		}
	}
	return info, fmt.Errorf("Service not found: %s", service)
}

type serviceList []ServiceInfo

func (a serviceList) Len() int           { return len(a) }
func (a serviceList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a serviceList) Less(i, j int) bool { return a[i].ServiceId < a[j].ServiceId }

func (s *serviceDirectory) Services() ([]ServiceInfo, error) {
	list := make([]ServiceInfo, 0, len(s.services))
	for _, info := range s.services {
		list = append(list, info)
	}
	sort.Sort(serviceList(list))
	return list, nil
}

func (s *serviceDirectory) RegisterService(newInfo ServiceInfo) (uint32, error) {
	if err := checkServiceInfo(newInfo); err != nil {
		return 0, err
	}
	for _, info := range s.staging {
		if info.Name == newInfo.Name {
			return 0, fmt.Errorf("Service name already staging: %s", info.Name)
		}
	}
	for _, info := range s.services {
		if info.Name == newInfo.Name {
			return 0, fmt.Errorf("Service name already ready: %s", info.Name)
		}
	}
	s.lastID++
	newInfo.ServiceId = s.lastID
	s.staging[s.lastID] = newInfo
	return s.lastID, nil
}

func (s *serviceDirectory) UnregisterService(id uint32) error {
	i, ok := s.services[id]
	if ok {
		delete(s.services, id)
		signal := s.signal
		if signal != nil {
			signal.SignalServiceRemoved(id, i.Name)
		}
		return nil
	}
	_, ok = s.staging[id]
	if ok {
		delete(s.staging, id)
		return nil
	}
	return fmt.Errorf("Service not found: %d", id)
}

func (s *serviceDirectory) ServiceReady(id uint32) error {
	i, ok := s.staging[id]
	if ok {
		delete(s.staging, id)
		s.services[id] = i
		signal := s.signal
		if signal != nil {
			signal.SignalServiceAdded(id, i.Name)
		}
		return nil
	}
	return fmt.Errorf("Service id not found: %d", id)
}

func (s *serviceDirectory) UpdateServiceInfo(i ServiceInfo) error {
	if err := checkServiceInfo(i); err != nil {
		return err
	}

	info, ok := s.services[i.ServiceId]
	if !ok {
		return fmt.Errorf("Service not found: %d (%s)", i.ServiceId, i.Name)
	}
	if info.Name != i.Name {
		return fmt.Errorf("Invalid name: %s (expected: %s)", i.Name,
			info.Name)
	}

	s.services[i.ServiceId] = i
	return nil
}

// MachineId returns a machine identifier.
func (s *serviceDirectory) MachineId() (string, error) {
	return util.MachineID(), nil
}

func (s *serviceDirectory) _socketOfService(P0 uint32) (
	o object.ObjectReference, err error) {
	return o, fmt.Errorf("_socketOfService not yet implemented")
}

func (s *serviceDirectory) Namespace(addr string) bus.Namespace {
	return &directoryNamespace{
		directory: s,
		addrs: []string{
			addr,
		},
	}
}

type directoryNamespace struct {
	directory *serviceDirectory
	addrs     []string
}

func (ns *directoryNamespace) Reserve(name string) (uint32, error) {
	info := ServiceInfo{
		Name:      name,
		ServiceId: 0,
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: ns.addrs,
		SessionId: "",
	}
	return ns.directory.RegisterService(info)
}

func (ns *directoryNamespace) Remove(serviceID uint32) error {
	return ns.directory.UnregisterService(serviceID)
}

func (ns *directoryNamespace) Enable(serviceID uint32) error {
	return ns.directory.ServiceReady(serviceID)
}
func (ns *directoryNamespace) Resolve(name string) (uint32, error) {
	info, err := ns.directory.Service(name)
	if err != nil {
		return 0, err
	}
	return info.ServiceId, nil
}
func (ns *directoryNamespace) Session(server bus.Server) bus.Session {
	return &directorySession{
		server:    server,
		namespace: ns,
	}
}

type directorySession struct {
	server    bus.Server
	namespace *directoryNamespace
}

func (s *directorySession) Proxy(name string, objectID uint32) (bus.Proxy, error) {

	info, err := s.namespace.directory.Service(name)
	if err != nil {
		return nil, fmt.Errorf("service not found: %s", name)
	}

	clt, err := s.client(info)
	if err != nil {
		return nil, err
	}
	meta, err := bus.GetMetaObject(clt, info.ServiceId, objectID)
	if err != nil {
		return nil, fmt.Errorf("call metaObject (service %d, object %d): %s",
			info.ServiceId, objectID, err)
	}
	return bus.NewProxy(clt, meta, info.ServiceId, objectID), nil

}
func (s *directorySession) client(info ServiceInfo) (bus.Client, error) {

	if info.MachineId == util.MachineID() &&
		info.ProcessId == util.ProcessID() {
		// if on the same process and listening one of the
		// same port, then by-pass connection
		for _, addr1 := range info.Endpoints {
			for _, addr2 := range s.namespace.addrs {
				if addr1 == addr2 {
					return s.server.Client(), nil
				}
			}
		}
	}
	_, endpoint, err := bus.SelectEndPoint(info.Endpoints, "", "")
	if err != nil {
		return nil, fmt.Errorf("object connection error (%s): %s",
			info.Name, err)
	}
	return bus.NewClient(endpoint), nil
}

func (s *directorySession) Object(ref object.ObjectReference) (bus.Proxy,
	error) {

	info, err := s.namespace.directory.info(ref.ServiceID)
	if err != nil {
		return nil, err
	}

	clt, err := s.client(info)
	if err != nil {
		return nil, err
	}
	return bus.NewProxy(clt, ref.MetaObject,
		ref.ServiceID, ref.ObjectID), nil
}

func (s *directorySession) Terminate() error {
	return nil
}

// ReadServiceInfo unmarshal a ServiceInfo struct.
func ReadServiceInfo(r io.Reader) (s ServiceInfo, err error) {
	return readServiceInfo(r)
}

// WriteServiceInfo marshal a ServiceInfo struct.
func WriteServiceInfo(s ServiceInfo, w io.Writer) (err error) {
	return writeServiceInfo(s, w)
}
