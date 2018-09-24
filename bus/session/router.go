package session

import (
	"github.com/lugu/qiloop/bus/net"
	"log"
	gonet "net"
	"sync"
)

type ClientSession struct {
	EndPoint      net.EndPoint
	Authenticated bool
}

func NewClientSession(c gonet.Conn) *ClientSession {
	return &ClientSession{net.NewEndPoint(c), false}
}

// Server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type Server struct {
	listen        gonet.Listener
	firewall      Firewall
	Router        Router
	sessions      map[*ClientSession]bool
	sessionsMutex sync.Mutex
}

func NewServer(l gonet.Listener, f Firewall, r Router) *Server {
	s := make(map[*ClientSession]bool)
	return &Server{
		listen:        l,
		firewall:      f,
		Router:        r,
		sessions:      s,
		sessionsMutex: sync.Mutex{},
	}
}

func (s *Server) handle(c gonet.Conn) error {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()
	session := NewClientSession(c)
	s.sessions[session] = true

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		return true, true
	}
	consumer := func(msg *net.Message) error {
		err := s.firewall.Inspect(msg, session)
		if err != nil {
			return err
		}
		return s.Router.Dispatch(msg, session)
	}
	session.EndPoint.AddHandler(filter, consumer)
	return nil
}

func (s *Server) Run() error {
	for {
		c, err := s.listen.Accept()
		if err != nil {
			return err
		}
		err = s.handle(c)
		if err != nil {
			log.Printf("Server connection error: %s", err)
			c.Close()
		}
	}
}

// CloseAll close the connecction. Return the first error if any.
func (s *Server) closeAll() error {
	var ret error = nil
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()
	for s, _ := range s.sessions {
		err := s.EndPoint.Close()
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
	Inspect(m *net.Message, from *ClientSession) error
}

// Router dispatch the incomming messages.
type Router interface {
	Add(n Namespace) (uint32, error)
	Remove(serviceID uint32) error
	Dispatch(m *net.Message, from *ClientSession) error
}

// Namespace represents a service
type Namespace interface {
	Add(o Object) (uint32, error)
	Remove(objectID uint32) error
	Dispatch(m *net.Message, from *ClientSession) error
	Ref(objectID uint32) error
	Unref(objectID uint32) error
}

type Object interface {
	CallID(action uint32, payload []byte) ([]byte, error)

	// SignalSubscribe returns a channel with the values of a signal
	SubscribeID(signal uint32, cancel chan int) (chan []byte, error)
}
