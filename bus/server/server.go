package server

import (
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"log"
	"math/rand"
	gonet "net"
	"sync"
)

// ErrServiceNotFound is returned with a message refers to an unknown
// service.
var ErrServiceNotFound = errors.New("Service not found")

// ErrObjectNotFound is returned with a message refers to an unknown
// object.
var ErrObjectNotFound = errors.New("Object not found")

// ErrActionNotFound is returned with a message refers to an unknown
// action.
var ErrActionNotFound = errors.New("Action not found")

// ErrNotAuthenticated is returned with a message tries to contact a
// service without prior authentication.
var ErrNotAuthenticated = errors.New("Not authenticated")

// ActionWrapper handles messages for an action.
type ActionWrapper func(payload []byte) ([]byte, error)

// Wrapper is used to dispatch messages to ActionWrapper.
type Wrapper map[uint32]ActionWrapper

// Terminator is to be called when an object whish to disapear.
type Terminator func()

// Activation is sent during activation: it informs the object of the
// context in which the object is being used.
type Activation struct {
	ServiceID uint32
	ObjectID  uint32
	Session   bus.Session
	Terminate Terminator
}

func objectTerminator(service *ServiceImpl, objectID uint32) Terminator {
	return func() {
		service.Remove(objectID)
	}
}

func serviceTerminator(router *Router, serviceID uint32) Terminator {
	return func() {
		router.Remove(serviceID)
		router.namespace.Remove(serviceID)
	}
}

func serviceActivation(router *Router, session bus.Session, serviceID uint32) Activation {
	return Activation{
		ServiceID: serviceID,
		ObjectID:  1,
		Session:   session,
		Terminate: serviceTerminator(router, serviceID),
	}
}

func objectActivation(service *ServiceImpl, session bus.Session, serviceID, objectID uint32) Activation {
	return Activation{
		ServiceID: serviceID,
		ObjectID:  objectID,
		Session:   session,
		Terminate: objectTerminator(service, objectID),
	}
}

// ServerObject interface used by Server to manipulate services.
type ServerObject interface {
	Receive(m *net.Message, from *Context) error
	Activate(activation Activation) error
	OnTerminate()
}

// ServiceImpl implements Service. It allows a service to manage the
// object within its domain.
type ServiceImpl struct {
	sync.RWMutex
	objects   map[uint32]ServerObject
	terminate Terminator
}

// NewService returns a service with the given object associated with
// object id 1.
func NewService(o ServerObject) *ServiceImpl {
	return &ServiceImpl{
		objects: map[uint32]ServerObject{
			1: o,
		},
	}
}

// Add is used to add an object to a service domain.
func (s *ServiceImpl) Add(o ServerObject) (uint32, error) {
	var index uint32
	// assign the first object to the index 0. following objects will
	// be assigned random values.
	s.Lock()
	defer s.Unlock()
	if len(s.objects) != 0 {
		index = rand.Uint32()
	}
	s.objects[index] = o
	return index, nil
}

// Activate informs the service it will become active and shall be
// ready to handle requests.
func (s *ServiceImpl) Activate(activation Activation) error {
	var wait sync.WaitGroup
	wait.Add(len(s.objects))
	s.terminate = activation.Terminate
	ret := make(chan error, len(s.objects))
	for objectID, obj := range s.objects {
		go func(obj ServerObject, objectID uint32) {
			objActivation := objectActivation(s, activation.Session,
				activation.ServiceID, objectID)
			err := obj.Activate(objActivation)
			if err != nil {
				ret <- err
			}
			wait.Done()
		}(obj, objectID)
	}
	wait.Wait()
	close(ret)
	for err := range ret {
		if err != nil {
			return err
		}
	}
	return nil
}

// Remove removes an object from the service domain.
func (s *ServiceImpl) Remove(objectID uint32) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.objects[objectID]; ok {
		delete(s.objects, objectID)
		return nil
	}
	return fmt.Errorf("Namespace: cannot remove object %d", objectID)
}

// Dispatch forwards the message to the appropriate object.
func (s *ServiceImpl) Dispatch(m *net.Message, from *Context) error {
	s.RLock()
	o, ok := s.objects[m.Header.Object]
	s.RUnlock()
	if ok {
		return o.Receive(m, from)
	}
	return util.ReplyError(from.EndPoint, m, ErrObjectNotFound)
}

// Terminate calls OnTerminate on all its objects.
func (s *ServiceImpl) Terminate() error {
	s.RLock()
	defer s.RUnlock()

	for _, obj := range s.objects {
		obj.OnTerminate()
	}
	if s.terminate != nil {
		s.terminate()
	}
	return nil
}

// Router dispatch the incomming messages. A Router shall be Activated
// before calling NewService.
type Router struct {
	sync.RWMutex
	services  map[uint32]*ServiceImpl
	namespace bus.Namespace
	session   bus.Session // nil until activation
}

// NewRouter construct a router with the service zero passed.
func NewRouter(authenticator ServerObject, namespace bus.Namespace) *Router {
	return &Router{
		services: map[uint32]*ServiceImpl{
			0: {
				objects: map[uint32]ServerObject{
					0: authenticator,
				},
			},
		},
		namespace: namespace,
		session:   nil,
	}
}

// Activate calls the Activate method on all the services. Only after
// the router can process messages.
func (r *Router) Activate(session bus.Session) error {
	r.Lock()
	if r.session != nil {
		r.Unlock()
		return fmt.Errorf("router already activated")
	}
	r.session = session
	r.Unlock()
	for serviceID, service := range r.services {
		activation := serviceActivation(r, session, serviceID)
		err := service.Activate(activation)
		if err != nil {
			return err
		}
	}
	return nil
}

// Terminate terminates all the services.
func (r *Router) Terminate() error {
	r.Lock()
	services := r.services
	r.services = make(map[uint32]*ServiceImpl)
	r.Unlock()

	var ret error
	for serviceID, service := range services {
		err := service.Terminate()
		if err != nil && ret == nil {
			ret = fmt.Errorf("service %d terminate: %s",
				serviceID, err)
		}
	}
	return ret
}

// NewService brings a new service online. It requires the router to
// be activated (as part of a session). The steps involves are:
// 1. request the name to the namespace (service directory)
// 2. activate the service
// 3. add the service to the router dispatcher
// 4. advertize the service to the namespace (service directory)
func (r *Router) NewService(name string, object ServerObject) (Service, error) {

	r.RLock()
	session := r.session
	r.RUnlock()

	// if the router is not yet activated, this is an error
	if session == nil {
		return nil, fmt.Errorf("cannot create service prior to activation")
	}

	service := NewService(object)
	// 1. reserve the name
	serviceID, err := r.namespace.Reserve(name)
	if err != nil {
		return nil, err
	}

	// 2. activate the service
	err = service.Activate(serviceActivation(r, session, serviceID))
	if err != nil {
		return nil, err
	}
	// 3. make it available
	err = r.Add(serviceID, service)
	if err != nil {
		return nil, err
	}
	// 4. advertize it
	err = r.namespace.Enable(serviceID)
	if err != nil {
		r.Remove(serviceID)
		return nil, err
	}
	return service, nil
}

// Add Add a service ot a router. This does not call the Activate()
// method on the service.
func (r *Router) Add(serviceID uint32, s *ServiceImpl) error {
	r.Lock()
	_, ok := r.services[serviceID]
	if ok {
		r.Unlock()
		return fmt.Errorf("service id already used: %d", serviceID)
	}
	r.services[serviceID] = s
	r.Unlock()
	return nil
}

// Remove removes a service from a router.
func (r *Router) Remove(serviceID uint32) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.services[serviceID]; ok {
		delete(r.services, serviceID)
		return nil
	}
	return fmt.Errorf("Router: cannot remove service %d", serviceID)
}

// Dispatch process a message. The message will be replied if the
// router can not found the destination.
func (r *Router) Dispatch(m *net.Message, from *Context) error {
	r.RLock()
	s, ok := r.services[m.Header.Service]
	r.RUnlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return util.ReplyError(from.EndPoint, m, ErrServiceNotFound)
}

// Context represents the context of the request
type Context struct {
	EndPoint      net.EndPoint
	Authenticated bool
}

// NewContext retuns a non authenticate context.
func NewContext(e net.EndPoint) *Context {
	return &Context{e, false}
}

// Firewall ensures an endpoint talks only to autorized services.
// Especially, it ensure authentication is passed.
func Firewall(m *net.Message, from *Context) error {
	if from.Authenticated == false && m.Header.Service != 0 {
		return ErrNotAuthenticated
	}
	return nil
}

// Server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type Server struct {
	listen        gonet.Listener
	addrs         []string
	namespace     bus.Namespace
	Router        *Router
	contexts      map[*Context]bool
	contextsMutex sync.Mutex
	closeChan     chan int
	waitChan      chan error
}

// NewServer starts a new server which listen to addr and serve
// incomming requests. It uses the given session to register and
// activate the new services.
func NewServer(sess bus.Session, addr string, auth Authenticator) (*Server, error) {
	l, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}

	addrs := []string{addr}
	namespace, err := Namespace(sess, addrs)
	if err != nil {
		return nil, err
	}
	router := NewRouter(ServiceAuthenticate(auth), namespace)
	s := &Server{
		listen:        l,
		addrs:         addrs,
		namespace:     namespace,
		Router:        router,
		contexts:      make(map[*Context]bool),
		contextsMutex: sync.Mutex{},
		closeChan:     make(chan int, 1),
		waitChan:      make(chan error, 1),
	}
	err = s.activate()
	if err != nil {
		return nil, err
	}
	go s.run()
	return s, nil
}

// StandAloneServer starts a new server
func StandAloneServer(listener gonet.Listener, auth Authenticator,
	namespace bus.Namespace) (*Server, error) {

	service0 := ServiceAuthenticate(auth)

	router := NewRouter(service0, namespace)

	s := &Server{
		listen:        listener,
		namespace:     namespace,
		Router:        router,
		contexts:      make(map[*Context]bool),
		contextsMutex: sync.Mutex{},
		closeChan:     make(chan int, 1),
		waitChan:      make(chan error, 1),
	}
	err := s.activate()
	if err != nil {
		return nil, err
	}
	go s.run()
	return s, nil
}

// Service represents a running service.
type Service interface {
	Terminate() error
}

// NewService returns a new service. The service is activated as part
// of the creation.
func (s *Server) NewService(name string, object ServerObject) (Service, error) {
	return s.Router.NewService(name, object)
}

func (s *Server) handle(c gonet.Conn, authenticated bool) {

	context := &Context{
		Authenticated: authenticated,
	}
	filter := func(hdr *net.Header) (matched bool, keep bool) {
		return true, true
	}
	consumer := func(msg *net.Message) error {
		err := Firewall(msg, context)
		if err != nil {
			log.Printf("missing authentication from %s: %#v",
				context.EndPoint.String(), msg.Header)
			return util.ReplyError(context.EndPoint, msg, err)
		}
		return s.Router.Dispatch(msg, context)
	}
	closer := func(err error) {
		s.contextsMutex.Lock()
		defer s.contextsMutex.Unlock()
		if _, ok := s.contexts[context]; ok {
			delete(s.contexts, context)
		}
	}
	finalize := func(e net.EndPoint) {
		context.EndPoint = e
		e.AddHandler(filter, consumer, closer)
		s.contextsMutex.Lock()
		s.contexts[context] = true
		s.contextsMutex.Unlock()
	}
	net.EndPointFinalizer(c, finalize)
}

func (s *Server) activate() error {
	err := s.Router.Activate(s.Session())
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) run() {
	for {
		c, err := s.listen.Accept()
		if err != nil {
			select {
			case <-s.closeChan:
			default:
				s.listen.Close()
				s.stoppedWith(err)
			}
			break
		}
		s.handle(c, false)
	}
}

func (s *Server) stoppedWith(err error) {
	// 1. informs all services
	s.Router.Terminate()
	// 2. close all connections
	s.closeAll()
	// 3. inform server's user
	s.waitChan <- err
	close(s.waitChan)
}

// CloseAll close the connecction. Return the first error if any.
func (s *Server) closeAll() error {
	var ret error
	s.contextsMutex.Lock()
	defer s.contextsMutex.Unlock()
	for context := range s.contexts {
		err := context.EndPoint.Close()
		if err != nil && ret == nil {
			ret = err
		}
	}
	return ret
}

// WaitTerminate blocks until the server has terminated.
func (s *Server) WaitTerminate() chan error {
	return s.waitChan
}

// Terminate stops a server.
func (s *Server) Terminate() error {
	close(s.closeChan)
	err := s.listen.Close()
	s.stoppedWith(err)
	return err
}

// Client returns a local client able to contact services without
// creating a new connection.
func (s *Server) Client() bus.Client {
	ctl, srv := gonet.Pipe()
	s.handle(srv, true)
	return client.NewClient(net.NewEndPoint(ctl))
}

// Session returns a local session able to contact local services
// without creating a new connection to the server.
func (s *Server) Session() bus.Session {
	return s.namespace.Session(s)
}
