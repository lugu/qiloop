package directory

import (
	"fmt"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
)

type ServiceDirectoryImpl struct {
	staging  map[uint32]ServiceInfo
	services map[uint32]ServiceInfo
	lastUuid uint32
	signal   ServiceDirectorySignalHelper
}

func NewServiceDirectory() ServiceDirectory {
	return &ServiceDirectoryImpl{
		staging:  make(map[uint32]ServiceInfo),
		services: make(map[uint32]ServiceInfo),
		lastUuid: 0,
	}
}

func (s *ServiceDirectoryImpl) Activate(sess *session.Session,
	serviceID, objectID uint32, signal ServiceDirectorySignalHelper) {
	s.signal = signal
}

func (s *ServiceDirectoryImpl) Service(service string) (info ServiceInfo, err error) {
	for _, info = range s.services {
		if info.Name == service {
			return info, nil
		}
	}
	return info, fmt.Errorf("Service not found: %s", service)
}

func (s *ServiceDirectoryImpl) Services() ([]ServiceInfo, error) {
	list := make([]ServiceInfo, 0, len(s.services))
	for _, info := range s.services {
		list = append(list, info)
	}
	return list, nil
}

func (s *ServiceDirectoryImpl) RegisterService(newInfo ServiceInfo) (uint32, error) {
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
	for id, info := range s.services {
		if info.Name == i.Name {
			s.services[id] = i
			return nil
		}
	}
	return fmt.Errorf("Service not found: %s", i.Name)
}

func (s *ServiceDirectoryImpl) MachineId() (string, error) {
	return util.MachineID(), nil
}

func (s *ServiceDirectoryImpl) _socketOfService(P0 uint32) (o object.ObjectReference, err error) {
	return o, fmt.Errorf("_socketOfService not yet implemented")
}
