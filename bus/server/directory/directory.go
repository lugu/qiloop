package directory

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/server"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"sort"
)

type ServiceDirectoryImpl struct {
	staging  map[uint32]ServiceInfo
	services map[uint32]ServiceInfo
	lastUuid uint32
	signal   ServiceDirectorySignalHelper
}

func NewServiceDirectory() *ServiceDirectoryImpl {
	return &ServiceDirectoryImpl{
		staging:  make(map[uint32]ServiceInfo),
		services: make(map[uint32]ServiceInfo),
		lastUuid: 0,
	}
}

func (s *ServiceDirectoryImpl) Activate(sess bus.Session, serviceID,
	objectID uint32, signal ServiceDirectorySignalHelper) error {
	s.signal = signal
	return nil
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
		return fmt.Errorf("missing end point")
	}
	for _, e := range i.Endpoints {
		if e == "" {
			return fmt.Errorf("empty endpoint not allowed")
		}
	}
	return nil
}

func (s *ServiceDirectoryImpl) Service(service string) (info ServiceInfo, err error) {
	for _, info = range s.services {
		if info.Name == service {
			return info, nil
		}
	}
	return info, fmt.Errorf("Service not found: %s", service)
}

type services []ServiceInfo

func (a services) Len() int           { return len(a) }
func (a services) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a services) Less(i, j int) bool { return a[i].ServiceId < a[j].ServiceId }

func (s *ServiceDirectoryImpl) Services() ([]ServiceInfo, error) {

	list := make([]ServiceInfo, 0, len(s.services))
	for _, info := range s.services {
		list = append(list, info)
	}
	sort.Sort(services(list))
	return list, nil
}

func (s *ServiceDirectoryImpl) RegisterService(newInfo ServiceInfo) (uint32, error) {
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
	s.lastUuid++
	newInfo.ServiceId = s.lastUuid
	s.staging[s.lastUuid] = newInfo
	return s.lastUuid, nil
}

func (s *ServiceDirectoryImpl) UnregisterService(id uint32) error {
	i, ok := s.services[id]
	if ok {
		delete(s.services, id)
		if s.signal != nil {
			s.signal.SignalServiceRemoved(id, i.Name)
		}
		return nil
	}
	_, ok = s.staging[id]
	if ok {
		delete(s.staging, id)
		return nil
	}
	return fmt.Errorf("Service not found %d", id)
}

func (s *ServiceDirectoryImpl) ServiceReady(id uint32) error {
	i, ok := s.staging[id]
	if ok {
		delete(s.staging, id)
		s.services[id] = i
		if s.signal != nil {
			s.signal.SignalServiceAdded(id, i.Name)
		}
		return nil
	}
	return fmt.Errorf("Service id not found: %d", id)
}

func (s *ServiceDirectoryImpl) UpdateServiceInfo(i ServiceInfo) error {
	if err := checkServiceInfo(i); err != nil {
		return err
	}

	info, ok := s.services[i.ServiceId]
	if !ok {
		return fmt.Errorf("Service not found: %d (%s)", i.ServiceId, i.Name)
	}
	if info.Name != i.Name { // can not change name
		return fmt.Errorf("Invalid name: %s (expected: %s)", i.Name, info.Name)
	}

	s.services[i.ServiceId] = i
	return nil
}

func (s *ServiceDirectoryImpl) MachineId() (string, error) {
	return util.MachineID(), nil
}

func (s *ServiceDirectoryImpl) _socketOfService(P0 uint32) (o object.ObjectReference, err error) {
	return o, fmt.Errorf("_socketOfService not yet implemented")
}

func (s *ServiceDirectoryImpl) Namespace() server.Namespace {
	panic("Not yet implemented")
}
