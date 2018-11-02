package session

import (
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
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

// TODO: implements type.Object methods and wrap them
type BasicObject struct {
	meta      object.MetaObject
	signals   map[uint32][]*Context
	Wrapper   bus.Wrapper
	serviceID uint32
	objectID  uint32
}

func NewObject(meta object.MetaObject) *BasicObject {
	return &BasicObject{
		meta:    meta,
		signals: make(map[uint32][]*Context),
		Wrapper: make(map[uint32]bus.ActionWrapper),
	}
}

func (o *BasicObject) MetaObject() (object.MetaObject, error) {
	return o.meta, nil
}

func (o *BasicObject) UpdateSignal(signal uint32, value []byte) (ret error) {
	if clients, ok := o.signals[signal]; ok {

		// FIXME: which id to use ? shall we record the register id
		// value?
		hdr := o.NewHeader(net.Event, signal, 0)
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

func (o *BasicObject) Receive(m *net.Message, from *Context) error {
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

// FIXME: must as well call the implementation's Terminate method.
func (o *BasicObject) Terminate() {
	panic("not yet implemented")
}

func (o *BasicObject) NewHeader(typ uint8, action, id uint32) net.Header {
	panic("not yet implemented")
}

// TODO: leave a door to call the specialized activate method
func (o *BasicObject) Activate(sess Session, serviceID, objectID uint32) {
	o.serviceID = serviceID
	o.objectID = objectID
	panic("not yet implemented")
}

type Object interface {
	Receive(m *net.Message, from *Context) error
	Activate(sess Session, serviceID, objectID uint32)
}

// ObjectDispatcher implements Object interface
type ObjectDispatcher struct {
	Wrapper bus.Wrapper
}

func (o *ObjectDispatcher) Activate(sess Session, serviceID, objectID uint32) {
}
func (o *ObjectDispatcher) Receive(m *net.Message, from *Context) error {
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

func (n *ServiceImpl) Dispatch(m *net.Message, from *Context) error {
	n.mutex.Lock()
	o, ok := n.objects[m.Header.Object]
	n.mutex.Unlock()
	if ok {
		return o.Receive(m, from)
	}
	return ObjectNotFound
}

func (n *ServiceImpl) Terminate() error {
	panic("not yet implemented")
}
func (n *ServiceImpl) WaitTerminate() chan int {
	panic("not yet implemented")
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

func (r *Router) Register(uid uint32, n *ServiceImpl) error {
	panic("not yet implemented")
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

func (r *Router) Dispatch(m *net.Message, from *Context) error {
	r.mutex.Lock()
	s, ok := r.services[m.Header.Service]
	r.mutex.Unlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return ServiceNotFound
}

// Context represents the context of the request
type Context struct {
	EndPoint      net.EndPoint
	Authenticated bool
}

func NewClientSession(c gonet.Conn) *Context {
	return &Context{net.NewEndPoint(c), false}
}

// Firewall ensures an endpoint talks only to autorized services.
// Especially, it ensure authentication is passed.
func Firewall(m *net.Message, from *Context) error {
	if from.Authenticated == false && m.Header.Service != 0 {
		return errors.New("Client not yet authenticated")
	}
	return nil
}

// Server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type Server struct {
	listen        gonet.Listener
	addrs         []string
	session       Session
	Router        *Router
	contexts      map[*Context]bool
	contextsMutex sync.Mutex
}

func NewServer(session Session, addr string) (*Server, error) {
	l, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		listen:        l,
		addrs:         []string{addr},
		session:       session,
		Router:        NewRouter(),
		contexts:      make(map[*Context]bool),
		contextsMutex: sync.Mutex{},
	}, nil
}

type Service interface {
	Terminate() error
	WaitTerminate() chan int
}

func (s *Server) NewService(name string, object Object) (Service, error) {

	info := services.ServiceInfo{
		Name:      name,
		ServiceId: 0,
		MachineId: util.MachineID(),
		ProcessId: util.ProcessID(),
		Endpoints: s.addrs,
		SessionId: "", // FIXME
	}

	uid, err := s.session.Directory.RegisterService(info)
	if err != nil {
		return nil, err
	}
	service := NewService(object)
	err = s.Router.Register(uid, service)
	if err != nil {
		return nil, err
	}
	// TODO
	// object.Activate(s.session, uid, 1)
	return service, nil
}

func NewServer2(l gonet.Listener, r *Router) *Server {
	s := make(map[*Context]bool)
	return &Server{
		listen:        l,
		Router:        r,
		contexts:      s,
		contextsMutex: sync.Mutex{},
	}
}

func (s *Server) handle(c gonet.Conn) error {
	s.contextsMutex.Lock()
	defer s.contextsMutex.Unlock()
	session := NewClientSession(c)
	s.contexts[session] = true

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
		s.contextsMutex.Lock()
		defer s.contextsMutex.Unlock()
		if _, ok := s.contexts[session]; ok {
			delete(s.contexts, session)
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
	s.contextsMutex.Lock()
	defer s.contextsMutex.Unlock()
	for s, _ := range s.contexts {
		go func(c *Context) {
			err := c.EndPoint.Close()
			if err != nil && ret == nil {
				ret = err
			}
		}(s)
	}
	return ret
}

func (s *Server) Stop() error {
	err := s.listen.Close()
	s.closeAll()
	return err
}
