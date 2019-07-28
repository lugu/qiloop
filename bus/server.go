package bus

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	gonet "net"
	"sync"

	"github.com/lugu/qiloop/bus/net"
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

// ErrTerminate is returned with an object lifetime ends while
// clients subscribes to its signals.
var ErrTerminate = errors.New("Object terminated")

// Activation is sent during activation: it informs the object of the
// context in which the object is being used.
type Activation struct {
	ServiceID uint32
	ObjectID  uint32
	Session   Session
	Terminate func()
	Service   Service
}

func objectTerminator(service Service, objectID uint32) func() {
	return func() {
		service.Remove(objectID)
	}
}

func serviceTerminator(router *Router, serviceID uint32) func() {
	return func() {
		router.Remove(serviceID)
		router.namespace.Remove(serviceID)
	}
}

func serviceActivation(router *Router, session Session, serviceID uint32) Activation {
	return Activation{
		ServiceID: serviceID,
		ObjectID:  1,
		Session:   session,
		Terminate: serviceTerminator(router, serviceID),
		Service:   nil,
	}
}

func objectActivation(service *serviceImpl, session Session, serviceID, objectID uint32) Activation {
	return Activation{
		ServiceID: serviceID,
		ObjectID:  objectID,
		Session:   session,
		Terminate: objectTerminator(service, objectID),
		Service:   service,
	}
}

// Receiver handles incomming messages.
type Receiver interface {
	Receive(m *net.Message, from Channel) error
}

// Actor interface used by Server to manipulate services.
type Actor interface {
	Receiver
	Activate(activation Activation) error
	OnTerminate()
}

type pendingObject struct{}

func (p pendingObject) Receive(m *net.Message, from Channel) error {
	return ErrObjectNotFound
}

func (p pendingObject) Activate(activation Activation) error {
	panic("can not activate pending object")
}

func (p pendingObject) OnTerminate() {
}

// serviceImpl implements Service. It allows a service to manage the
// object within its domain.
type serviceImpl struct {
	sync.RWMutex
	objects   map[uint32]Actor
	terminate func()
	session   Session
	serviceID uint32
}

// newService returns a service with the given object associated with
// object id 1.
func newService(o Actor) *serviceImpl {
	return &serviceImpl{
		objects: map[uint32]Actor{
			1: o,
		},
	}
}

func (s *serviceImpl) ServiceID() uint32 {
	return s.serviceID
}

// Add is used to add an object to a service domain.
func (s *serviceImpl) Add(obj Actor) (index uint32, err error) {
	// assign the first object to the index 0. following objects will
	// be assigned random values.
	s.Lock()
	if _, ok := s.objects[1]; ok {
		index = (rand.Uint32() << 1) >> 1
		if _, ok = s.objects[index]; ok {
			s.Unlock()
			return s.Add(obj)
		}
	}
	if s.session == nil { // service not yet activated
		s.objects[index] = obj
		s.Unlock()
		return
	}
	s.objects[index] = pendingObject{}
	s.Unlock()

	a := objectActivation(s, s.session, s.serviceID, index)
	err = obj.Activate(a)

	s.Lock()
	if err != nil {
		s.objects[index] = nil
	} else {
		s.objects[index] = obj
	}
	s.Unlock()
	return
}

// Activate informs the service it will become active and shall be
// ready to handle requests. activation.Service is nil.
func (s *serviceImpl) Activate(activation Activation) error {
	var wait sync.WaitGroup
	s.terminate = activation.Terminate
	s.session = activation.Session
	s.serviceID = activation.ServiceID
	s.Lock()
	ret := make(chan error, len(s.objects))
	wait.Add(len(s.objects))
	for objectID, obj := range s.objects {
		go func(obj Actor, objectID uint32) {
			objActivation := objectActivation(s, activation.Session,
				activation.ServiceID, objectID)
			err := obj.Activate(objActivation)
			if err != nil {
				ret <- err
			}
			wait.Done()
		}(obj, objectID)
	}
	s.Unlock()
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
func (s *serviceImpl) Remove(objectID uint32) error {
	s.Lock()
	if obj, ok := s.objects[objectID]; ok {
		delete(s.objects, objectID)
		s.Unlock()
		obj.OnTerminate()
		return nil
	}
	s.Unlock()
	return fmt.Errorf("cannot remove object %d", objectID)
}

// Dispatch forwards the message to the appropriate object.
func (s *serviceImpl) Dispatch(m *net.Message, from Channel) error {
	s.RLock()
	o, ok := s.objects[m.Header.Object]
	s.RUnlock()
	if ok {
		return o.Receive(m, from)
	}
	return from.SendError(m, ErrObjectNotFound)
}

// Terminate calls OnTerminate on all its objects.
func (s *serviceImpl) Terminate() error {
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
	services  map[uint32]*serviceImpl
	namespace Namespace
	session   Session // nil until activation
}

// NewRouter construct a router with the service zero passed.
func NewRouter(authenticator Actor, namespace Namespace) *Router {
	return &Router{
		services: map[uint32]*serviceImpl{
			0: {
				objects: map[uint32]Actor{
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
func (r *Router) Activate(session Session) error {
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
	r.services = make(map[uint32]*serviceImpl)
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

// Add Add a service ot a router. This does not call the Activate()
// method on the service.
func (r *Router) Add(serviceID uint32, s *serviceImpl) error {
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
func (r *Router) Dispatch(m *net.Message, from Channel) error {
	r.RLock()
	s, ok := r.services[m.Header.Service]
	r.RUnlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return from.SendError(m, ErrServiceNotFound)
}

// firewall ensures an endpoint talks only to autorized services.
// Especially, it ensure authentication is passed.
func firewall(m *net.Message, from Channel) error {
	if from.Authenticated() == false && m.Header.Service != 0 {
		return ErrNotAuthenticated
	}
	return nil
}

// server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type server struct {
	listen        net.Listener
	addrs         []string
	namespace     Namespace
	Router        *Router
	contexts      map[Channel]bool
	contextsMutex sync.Mutex
	closeChan     chan int
	waitChan      chan error
}

// NewServer creates a new server which respond to incomming
// connection requests.
func NewServer(listener net.Listener, auth Authenticator,
	namespace Namespace, service1 Actor) (Server, error) {

	service0 := ServiceAuthenticate(auth)

	router := NewRouter(service0, namespace)

	s := &server{
		listen:        listener,
		namespace:     namespace,
		Router:        router,
		contexts:      make(map[Channel]bool),
		contextsMutex: sync.Mutex{},
		closeChan:     make(chan int, 1),
		waitChan:      make(chan error, 1),
	}
	err := s.activate()
	if err != nil {
		return nil, err
	}

	_, err = s.NewService("ServiceDirectory", service1)
	if err != nil {
		s.Terminate()
		return nil, err
	}

	go s.run()
	return s, nil
}

// StandAloneServer starts a new server
func StandAloneServer(listener net.Listener, auth Authenticator,
	namespace Namespace) (Server, error) {

	service0 := ServiceAuthenticate(auth)

	router := NewRouter(service0, namespace)

	s := &server{
		listen:        listener,
		namespace:     namespace,
		Router:        router,
		contexts:      make(map[Channel]bool),
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
	ServiceID() uint32
	Add(o Actor) (uint32, error)
	Remove(objectID uint32) error
	Terminate() error
}

// NewService returns a new service. The service is activated as part
// of the creation.
// This brings the new service online. The steps involves are:
// 1. request the name to the namespace (service directory)
// 2. activate the service
// 3. add the service to the router dispatcher
// 4. advertize the service to the namespace (service directory)
func (s *server) NewService(name string, object Actor) (Service, error) {

	s.Router.RLock()
	session := s.Router.session
	s.Router.RUnlock()

	// if the router is not yet activated, this is an error
	if session == nil {
		return nil, fmt.Errorf("cannot create service prior to activation")
	}

	// 1. reserve the name
	serviceID, err := s.namespace.Reserve(name)
	if err != nil {
		return nil, err
	}
	service := newService(object)

	// 2. activate the service
	err = service.Activate(serviceActivation(s.Router, session, serviceID))
	if err != nil {
		return nil, err
	}
	// 3. make it available
	err = s.Router.Add(serviceID, service)
	if err != nil {
		return nil, err
	}
	// 4. advertize it
	err = s.namespace.Enable(serviceID)
	if err != nil {
		s.Router.Remove(serviceID)
		return nil, err
	}
	return service, nil
}

func (s *server) handle(stream net.Stream, authenticated bool) {

	context := &channel{
		capability: PreferedCap("", ""),
	}
	if authenticated {
		context.SetAuthenticated()
	}
	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Type == net.Reply || hdr.Type == net.Error ||
			hdr.Type == net.Event || hdr.Type == net.Cancelled {
			return false, true
		}
		return true, true
	}
	consumer := func(msg *net.Message) error {
		err := firewall(msg, context)
		if err != nil {
			log.Printf("missing authentication from %s: %#v",
				context.EndPoint().String(), msg.Header)
			return context.SendError(msg, err)
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
		context.endpoint = e
		e.AddHandler(filter, consumer, closer)
		s.contextsMutex.Lock()
		s.contexts[context] = true
		s.contextsMutex.Unlock()
	}
	net.EndPointFinalizer(stream, finalize)
}

func (s *server) activate() error {
	err := s.Router.Activate(s.Session())
	if err != nil {
		return err
	}
	return nil
}

func (s *server) run() {
	for {
		stream, err := s.listen.Accept()
		if err != nil {
			select {
			case <-s.closeChan:
			default:
				s.listen.Close()
				s.stoppedWith(err)
			}
			break
		}
		s.handle(stream, false)
	}
}

func (s *server) stoppedWith(err error) {
	// 1. informs all services
	s.Router.Terminate()
	// 2. close all connections
	s.closeAll()
	// 3. inform server's user
	s.waitChan <- err
	close(s.waitChan)
}

// CloseAll close the connecction. Return the first error if any.
func (s *server) closeAll() error {
	var ret error
	s.contextsMutex.Lock()
	defer s.contextsMutex.Unlock()
	for context := range s.contexts {
		err := context.EndPoint().Close()
		if err != nil && ret == nil {
			ret = err
		}
	}
	return ret
}

// WaitTerminate blocks until the server has terminated.
func (s *server) WaitTerminate() chan error {
	return s.waitChan
}

// Terminate stops a server.
func (s *server) Terminate() error {
	close(s.closeChan)
	err := s.listen.Close()
	s.stoppedWith(err)
	return err
}

// Client returns a local client able to contact services without
// creating a new connection.
func (s *server) Client() Client {
	ctl, srv := gonet.Pipe()
	s.handle(net.ConnStream(srv), true)
	return NewClient(net.ConnEndPoint(ctl))
}

// Session returns a local session able to contact local services
// without creating a new connection to the server.
func (s *server) Session() Session {
	return s.namespace.Session(s)
}
