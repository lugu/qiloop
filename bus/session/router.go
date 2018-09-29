package session

import (
	"github.com/lugu/qiloop/bus/net"
	"log"
	gonet "net"
	"sync"
)

// Server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type Server struct {
	listen         gonet.Listener
	firewall       Firewall
	Router         Router
	endpoints      []net.EndPoint
	endpointsMutex sync.Mutex
}

func NewServer(l gonet.Listener, f Firewall, r Router) *Server {
	return &Server{l, f, r, make([]net.EndPoint, 0, 100), sync.Mutex{}}
}

func (s *Server) add(e net.EndPoint) error {
	s.endpointsMutex.Lock()
	defer s.endpointsMutex.Unlock()
	s.endpoints = append(s.endpoints, e)

	var f net.Filter = func(hdr *net.Header) (matched bool, keep bool) {
		return true, true
	}
	var c net.Consumer = func(msg *net.Message) error {
		err := s.firewall.Inspect(msg, e)
		if err != nil {
			return err
		}
		return s.Router.Dispatch(msg, e)
	}
	e.AddHandler(f, c)
	return nil
}

func (s *Server) Run() error {
	for {
		c, err := s.listen.Accept()
		if err != nil {
			return err
		}
		endpoint := net.NewEndPoint(c)
		err = s.add(endpoint)
		if err != nil {
			log.Printf("Server connection error: %s", err)
			endpoint.Close()
		}
	}
}

// CloseAll close the connecction. Return the first error if any.
func (s *Server) closeAll() error {
	var ret error = nil
	s.endpointsMutex.Lock()
	defer s.endpointsMutex.Unlock()
	for _, e := range s.endpoints {
		err := e.Close()
		if err != nil && ret == nil {
			ret = err
		}
	}
	return ret
}

func (s *Server) Stop() error {
	err := s.listen.Close()
	s.closeAll()
	return err
}

// Firewall ensures an endpoint talks only to the autorized service.
// Especially, it ensure authentication is passed.
type Firewall interface {
	Inspect(m *net.Message, from net.EndPoint) error
}

// Router dispatch the incomming messages.
type Router interface {
	Add(n Namespace) (uint32, error)
	Remove(serviceID uint32) error
	Dispatch(m *net.Message, from net.EndPoint) error
}

// Namespace represents a service
type Namespace interface {
	Add(o Object) (uint32, error)
	Remove(objectID uint32) error
	Dispatch(m *net.Message, from net.EndPoint) error
	Ref(objectID uint32) error
	Unref(objectID uint32) error
}

type Object interface {
	CallID(action uint32, payload []byte) ([]byte, error)

	// SignalSubscribe returns a channel with the values of a signal
	SubscribeID(signal uint32, cancel chan int) (chan []byte, error)
}
