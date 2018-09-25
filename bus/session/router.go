package session

import (
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"log"
	"math/rand"
	gonet "net"
	"sync"
)

var ObjectNotFound error = errors.New("Object not found")
var ServiceNotFound error = errors.New("Service not found")

type Object interface {
	Receive(m *net.Message, from *ClientSession) error
}

type Namespace struct {
	objects map[uint32]Object
	mutex   sync.Mutex
}

func (n *Namespace) Add(o Object) (uint32, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	var index uint32 = 0
	// assign the first object to the index 0. following objects will
	// be assigned random values.
	if len(n.objects) != 0 {
		index = rand.Uint32()
	}
	n.objects[index] = o
	return index, nil
}

func (n Namespace) Remove(objectID uint32) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if _, ok := n.objects[objectID]; ok {
		delete(n.objects, objectID)
		return nil
	}
	return fmt.Errorf("Namespace: cannot remove object %d", objectID)
}

func (n Namespace) Dispatch(m *net.Message, from *ClientSession) error {
	n.mutex.Lock()
	o, ok := n.objects[m.Header.Object]
	n.mutex.Unlock()
	if ok {
		return o.Receive(m, from)
	}
	return ObjectNotFound
}

// Router dispatch the incomming messages.
type Router struct {
	services  map[uint32]Namespace
	nextIndex uint32
	mutex     sync.Mutex
}

func (r *Router) Add(n Namespace) (uint32, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.services[r.nextIndex] = n
	r.nextIndex++
	return r.nextIndex, nil
}

func (r Router) Remove(serviceID uint32) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.services[serviceID]; ok {
		delete(r.services, serviceID)
		return nil
	}
	return fmt.Errorf("Router: cannot remove service %d", serviceID)
}

func (r Router) Dispatch(m *net.Message, from *ClientSession) error {
	r.mutex.Lock()
	s, ok := r.services[m.Header.Service]
	r.mutex.Unlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return ServiceNotFound
}

type ClientSession struct {
	EndPoint      net.EndPoint
	Authenticated bool
}

func NewClientSession(c gonet.Conn) *ClientSession {
	return &ClientSession{net.NewEndPoint(c), false}
}

// Firewall ensures an endpoint talks only to autorized services.
// Especially, it ensure authentication is passed.
func Firewall(m *net.Message, from *ClientSession) error {
	if from.Authenticated == false && m.Header.Service != 0 {
		return errors.New("Client not yet authenticated")
	}
	return nil
}

// Server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type Server struct {
	listen        gonet.Listener
	Router        Router
	sessions      map[*ClientSession]bool
	sessionsMutex sync.Mutex
}

func NewServer(l gonet.Listener, r Router) *Server {
	s := make(map[*ClientSession]bool)
	return &Server{
		listen:        l,
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
		err := Firewall(msg, session)
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
