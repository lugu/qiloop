package bus

import (
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/type/object"
)

// Cache implements Session interface without connecting to a
// service directory.
type Cache struct {
	Names    map[string]uint32
	Services map[uint32]object.MetaObject
	Endpoint net.EndPoint
}

// Proxy returns a proxy object to the desired service.
func (s *Cache) Proxy(name string, objectID uint32) (Proxy, error) {
	serviceID, ok := s.Names[name]
	if !ok {
		return nil, fmt.Errorf("service not cached: %s", name)
	}
	meta := s.Services[serviceID]

	client := NewClient(s.Endpoint)
	return NewProxy(client, meta, serviceID, objectID), nil
}

// Object creates an object from a reference.
func (s *Cache) Object(ref object.ObjectReference) (o Proxy, err error) {
	return o, fmt.Errorf("Not yet implemented")
}

// Destroy closes the connection.
func (s *Cache) Destroy() error {
	return s.Endpoint.Close()
}

// AddService manually associates the service name with a service id
// and a meta objects.
func (s *Cache) AddService(name string, serviceID uint32,
	meta object.MetaObject) {

	s.Names[name] = serviceID
	s.Services[serviceID] = meta
}

// Lookup query the given service id for its meta object and add it to
// the cache.
func (s *Cache) Lookup(name string, serviceID uint32) error {
	objectID := uint32(1)
	meta, err := GetMetaObject(NewClient(s.Endpoint),
		serviceID, objectID)
	if err != nil {
		return fmt.Errorf("Can not reach metaObject: %s", err)
	}
	s.AddService(name, serviceID, meta)
	return nil
}

// NewCachedSession authenticate to the address and returns a Cache
// object which can act as a Session.
func NewCachedSession(addr string) (*Cache, error) {
	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}
	err = Authenticate(endpoint)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %s", err)
	}

	return NewCache(endpoint), nil
}

// NewCache returns a new cache.
func NewCache(e net.EndPoint) *Cache {
	return &Cache{
		Names:    make(map[string]uint32),
		Services: make(map[uint32]object.MetaObject),
		Endpoint: e,
	}
}
