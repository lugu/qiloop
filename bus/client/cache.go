package client

import (
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/object"
)

// Cache implments bus.Session interface
type Cache struct {
	Names    map[string]uint32
	Services map[uint32]object.MetaObject
	Endpoint net.EndPoint
}

func (s *Cache) Proxy(name string, objectID uint32) (bus.Proxy, error) {
	serviceID, ok := s.Names[name]
	if !ok {
		return nil, fmt.Errorf("service not cached: %s", name)
	}
	meta := s.Services[serviceID]

	client := NewClient(s.Endpoint)
	return NewProxy(client, meta, serviceID, objectID), nil
}

func (s *Cache) Object(ref object.ObjectReference) (o object.Object, err error) {
	return o, fmt.Errorf("Not yet implemented")
}

func (s *Cache) Destroy() error {
	return s.Endpoint.Close()
}

func (s *Cache) AddService(name string, serviceID uint32,
	meta object.MetaObject) {
	s.Names[name] = serviceID
	s.Services[serviceID] = meta
}

func NewCache(e net.EndPoint) Cache {
	return Cache{
		Names:    make(map[string]uint32),
		Services: make(map[uint32]object.MetaObject),
		Endpoint: e,
	}
}

func (s *Cache) Session() bus.Session {
	return s
}
