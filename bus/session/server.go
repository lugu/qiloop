package session

import (
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/object"
	"log"
	"math/rand"
	gonet "net"
	"sync"
)

var ServiceNotFound error = errors.New("Service not found")
var ObjectNotFound error = errors.New("Object not found")
var ActionNotFound error = errors.New("Action not found")

// From IDL file, generate stub interfaces like:
type ObjectImplA interface {
}

func NewObjectImplAConstructor(impl ObjectImplA) ObjectImplConstructor {
	return nil
}

type ObjectImplConstructor interface {
	Instanciate(service, object uint32, namespace ServiceImpl) ObjectImpl
}

type ObjectImpl interface {
	Destroy() error
	Ref() object.ObjectReference
}

type ActionHelper interface {
	Unregister() error
	NewHeader(typ uint8, action, id uint32) net.Header
}

func NewObject(meta object.MetaObject) ObjectWrapper {
	panic("not yet implemented")
	return nil
}

// ObjectWrapper implements the methods of type.Object
type ObjectWrapper interface {
	// UpdateSignal let an implementation update a signal without
	// having to know who listen to it.
	UpdateSignal(signalID uint32, value []byte) error
	// Terminate let an implementation destroy itself.
	Terminate()
	// Wrapper returns the wrapper for the actions of bus.Object.
	// Other actions need to be consolidated before the creation of
	// the ObjectDispather.
	Wrapper() bus.Wrapper
}

type GenericObject struct {
	signals      map[uint32][]*ClientSession
	meta         object.MetaObject
	wrapper      bus.Wrapper // TODO: implements object.Object
	actionHelper ActionHelper
}

func (o *GenericObject) MetaObject() (object.MetaObject, error) {
	return o.meta, nil
}

func (o *GenericObject) UpdateSignal(signal uint32, value []byte) (ret error) {
	if clients, ok := o.signals[signal]; ok {

		// FIXME: which id to use ? shall we record the register id
		// value?
		hdr := o.actionHelper.NewHeader(net.Event, signal, 0)
		msg := net.NewMessage(hdr, value)

		for _, client := range clients {
			err := client.EndPoint.Send(msg)
			if err != nil {
				ret = err
			}
		}
	}
	return ret
}

// FIXME: must as well call the implementation's Terminate method.
func (o *GenericObject) Terminate() {
	o.actionHelper.Unregister()
}

func (o *GenericObject) Wrapper() bus.Wrapper {
	panic("not yet implemented")
	return nil
}

type Object interface {
	Receive(m *net.Message, from *ClientSession) error
}

// ObjectDispather implements Object interface
type ObjectDispather struct {
	Wrapper bus.Wrapper
}

func (o *ObjectDispather) Receive(m *net.Message, from *ClientSession) error {
	a, ok := o.Wrapper[m.Header.Action]
	if !ok {
		return ActionNotFound
	}
	response, err := a(m.Payload)

	if err != nil {
		mError := net.NewMessage(m.Header, util.ErrorPaylad(err))
		mError.Header.Type = net.Error
		return from.EndPoint.Send(mError)
	}
	reply := net.NewMessage(m.Header, response)
	reply.Header.Type = net.Reply
	return from.EndPoint.Send(reply)
}

type ServiceImpl struct {
	objects map[uint32]Object
	mutex   sync.Mutex
}

func NewService(o Object) *ServiceImpl {
	return &ServiceImpl{
		objects: map[uint32]Object{
			0: o,
		},
	}
}

func (n *ServiceImpl) Add(o Object) (uint32, error) {
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

func (n *ServiceImpl) Remove(objectID uint32) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if _, ok := n.objects[objectID]; ok {
		delete(n.objects, objectID)
		return nil
	}
	return fmt.Errorf("Namespace: cannot remove object %d", objectID)
}

func (n *ServiceImpl) Dispatch(m *net.Message, from *ClientSession) error {
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
	services  map[uint32]*ServiceImpl
	nextIndex uint32
	mutex     sync.Mutex
}

func NewDirectoryService() *ServiceImpl {
	// FIXME: implement the service
	var o Object
	return NewService(o)
}

func NewRouter() *Router {
	return &Router{
		services:  make(map[uint32]*ServiceImpl),
		nextIndex: 0,
	}
}

func (r *Router) Add(n *ServiceImpl) (uint32, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.services[r.nextIndex] = n
	r.nextIndex++
	return r.nextIndex, nil
}

func (r *Router) Remove(serviceID uint32) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.services[serviceID]; ok {
		delete(r.services, serviceID)
		return nil
	}
	return fmt.Errorf("Router: cannot remove service %d", serviceID)
}

func (r *Router) Dispatch(m *net.Message, from *ClientSession) error {
	r.mutex.Lock()
	s, ok := r.services[m.Header.Service]
	r.mutex.Unlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return ServiceNotFound
}

// ClientSession represents the context of the request
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
	Router        *Router
	sessions      map[*ClientSession]bool
	sessionsMutex sync.Mutex
}

func NewServer(l gonet.Listener, r *Router) *Server {
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
		err = s.Router.Dispatch(msg, session)
		if err != nil {
			log.Printf("server warning: %s", err)
		}
		return nil

	}
	closer := func(err error) {
		s.sessionsMutex.Lock()
		defer s.sessionsMutex.Unlock()
		if _, ok := s.sessions[session]; ok {
			delete(s.sessions, session)
		}
	}
	session.EndPoint.AddHandler(filter, consumer, closer)
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
		go func() {
			err := s.EndPoint.Close()
			if err != nil && ret == nil {
				ret = err
			}
		}()
	}
	return ret
}

func (s *Server) Stop() error {
	err := s.listen.Close()
	s.closeAll()
	return err
}
