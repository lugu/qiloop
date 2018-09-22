package impl

import (
	"fmt"
	. "github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/type/object"
)

type ServiceDirectoryInterface interface {
	Service(P0 string) (ServiceInfo, error)
	Services() ([]ServiceInfo, error)
	RegisterService(P0 ServiceInfo) (uint32, error)
	UnregisterService(P0 uint32) error
	ServiceReady(P0 uint32) error
	UpdateServiceInfo(P0 ServiceInfo) error
	MachineId() (string, error)
	_socketOfService(P0 uint32) (object.ObjectReference, error)
}

type ServiceDirectoryImpl struct {
	staging  map[uint32]ServiceInfo
	services map[uint32]ServiceInfo
	lastUuid uint32
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
	_, ok := s.services[id]
	if ok {
		delete(s.services, id)
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
