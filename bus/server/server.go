package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
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

type signalUser struct {
	signalID  uint32
	messageID uint32
	context   *Context
	clientID  uint64
}

// ActionWrapper handles messages for an action.
type ActionWrapper func(payload []byte) ([]byte, error)

// Wrapper is used to dispatch messages to ActionWrapper.
type Wrapper map[uint32]ActionWrapper

// BasicObject implements the Object interface. It handles the generic
// method and signal. Services implementation embedded a BasicObject
// and fill it with the extra actions they wish to handle using the
// Wrap method. See type/object.Object for a list of the default
// methods.
type BasicObject struct {
	meta      object.MetaObject
	signals   []signalUser
	Wrapper   Wrapper
	serviceID uint32
	objectID  uint32
}

// NewObject construct a BasicObject from a MetaObject.
func NewObject(meta object.MetaObject) *BasicObject {
	var obj BasicObject
	obj.meta = object.FullMetaObject(meta)
	obj.signals = make([]signalUser, 0)
	obj.Wrapper = make(map[uint32]ActionWrapper)
	obj.Wrap(uint32(0x2), obj.wrapMetaObject)
	// obj.Wrapper[uint32(0x3)] = obj.Terminate
	// obj.Wrapper[uint32(0x5)] = obj.Property
	// obj.Wrapper[uint32(0x6)] = obj.SetProperty
	// obj.Wrapper[uint32(0x7)] = obj.Properties
	// obj.Wrapper[uint32(0x8)] = obj.RegisterEventWithSignature
	return &obj
}

// Wrap let a BasicObject owner extend it with custom actions.
func (o *BasicObject) Wrap(id uint32, fn ActionWrapper) {
	o.Wrapper[id] = fn
}

func (o *BasicObject) addSignalUser(signalID, messageID uint32, from *Context) uint64 {
	clientID := rand.Uint64()
	newUser := signalUser{
		signalID,
		messageID,
		from,
		clientID,
	}
	o.signals = append(o.signals, newUser)
	return clientID
}

func (o *BasicObject) removeSignalUser(id uint64) error {
	return nil
}

func (o *BasicObject) handleRegisterEvent(from *Context, msg *net.Message) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	if objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	signalID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read signal uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	_, err = basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	messageID := msg.Header.ID
	clientID := o.addSignalUser(signalID, messageID, from)
	var out bytes.Buffer
	err = basic.WriteUint64(clientID, &out)
	if err != nil {
		err = fmt.Errorf("cannot write client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	return o.reply(from, msg, out.Bytes())
}

func (o *BasicObject) handleUnregisterEvent(from *Context, msg *net.Message) error {
	buf := bytes.NewBuffer(msg.Payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read object uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	if objectID != o.objectID {
		err := fmt.Errorf("wrong object id, expecting %d, got %d",
			msg.Header.Object, objectID)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	_, err = basic.ReadUint32(buf)
	if err != nil {
		err = fmt.Errorf("cannot read action uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	clientID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	err = o.removeSignalUser(clientID)
	if err != nil {
		return util.ReplyError(from.EndPoint, msg, err)
	}
	var out bytes.Buffer
	return o.reply(from, msg, out.Bytes())
}

func (o *BasicObject) wrapMetaObject(payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer(payload)
	objectID, err := basic.ReadUint32(buf)
	if err != nil {
		return nil, err
	}
	if objectID != o.objectID {
		return nil, fmt.Errorf("invalid object id: %d instead of %d",
			objectID, o.objectID)
	}
	var out bytes.Buffer
	err = object.WriteMetaObject(o.meta, &out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// UpdateSignal informs the registered clients of the new state.
func (o *BasicObject) UpdateSignal(signal uint32, value []byte) error {
	var ret error
	for _, client := range o.signals {
		if client.signalID == signal {
			hdr := o.newHeader(net.Event, signal, client.messageID)
			msg := net.NewMessage(hdr, value)
			// FIXME: catch writing to close connection
			err := client.context.EndPoint.Send(msg)
			if err != nil {
				ret = err
			}
		}
	}
	return ret
}

func (o *BasicObject) reply(from *Context, m *net.Message, response []byte) error {
	hdr := o.newHeader(net.Reply, m.Header.Action, m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return from.EndPoint.Send(reply)
}

func (o *BasicObject) handleDefault(from *Context, msg *net.Message) error {
	a, ok := o.Wrapper[msg.Header.Action]
	if !ok {
		return util.ReplyError(from.EndPoint, msg, ErrActionNotFound)
	}
	response, err := a(msg.Payload)
	if err != nil {
		return util.ReplyError(from.EndPoint, msg, err)
	}
	return o.reply(from, msg, response)
}

// Receive processes the incoming message and responds to the client.
// The returned error is not destinated to the client which have
// already be replied.
func (o *BasicObject) Receive(m *net.Message, from *Context) error {
	switch m.Header.Action {
	case 0x0:
		return o.handleRegisterEvent(from, m)
	case 0x1:
		return o.handleUnregisterEvent(from, m)
	default:
		return o.handleDefault(from, m)
	}
}

// OnTerminate is called when the object is terminated.
func (o *BasicObject) OnTerminate() {
}

func (o *BasicObject) newHeader(typ uint8, action, id uint32) net.Header {
	return net.NewHeader(typ, o.serviceID, o.objectID, action, id)
}

// Activate informs the object when it becomes online. After this
// method returns, the object will start receiving incomming messages.
func (o *BasicObject) Activate(sess bus.Session, serviceID,
	objectID uint32) error {
	o.serviceID = serviceID
	o.objectID = objectID
	return nil
}

// Object interface used by Server to manipulate services.
type Object interface {
	Receive(m *net.Message, from *Context) error
	Activate(sess bus.Session, serviceID, objectID uint32) error
	// - call Server.directory.Remove()
	// - call Server.Router.Remove()
	// - call Terminate on all objects
	// - signal WaitTerminate condition
	OnTerminate()
}

// ServiceImpl implements Service. It allows a service to manage the
// object within its domain.
type ServiceImpl struct {
	sync.RWMutex
	objects map[uint32]Object
}

// NewService returns a service with the given object associated with
// object id 1.
func NewService(o Object) *ServiceImpl {
	return &ServiceImpl{
		objects: map[uint32]Object{
			1: o,
		},
	}
}

// Add is used to add an object to a service domain.
func (s *ServiceImpl) Add(o Object) (uint32, error) {
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

// Activate is called when a Service becomes online. After the
// Activate method is call, the service will start receiving incomming
// messages.
func (s *ServiceImpl) Activate(sess bus.Session, serviceID uint32) error {
	var wait sync.WaitGroup
	wait.Add(len(s.objects))
	ret := make(chan error, len(s.objects))
	for objectID, obj := range s.objects {
		go func(obj Object, objectID uint32) {
			err := obj.Activate(sess, serviceID, objectID)
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
	return nil
}

// Router dispatch the incomming messages.
type Router struct {
	sync.RWMutex
	services map[uint32]*ServiceImpl
}

// NewRouter construct a router with the service zero passed.
func NewRouter(authenticator Object) *Router {
	return &Router{
		services: map[uint32]*ServiceImpl{
			0: {
				objects: map[uint32]Object{
					0: authenticator,
				},
			},
		},
	}
}

// Activate calls the Activate method on all the services. Only after
// the router can process messages.
func (r *Router) Activate(sess bus.Session) error {
	r.RLock()
	defer r.RUnlock()
	for serviceID, service := range r.services {
		err := service.Activate(sess, serviceID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Terminate terminates all the services.
func (r *Router) Terminate() error {
	r.RLock()
	defer r.RUnlock()
	var ret error
	for serviceID, service := range r.services {
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
	session       bus.Session
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
	s := &Server{
		listen:        l,
		addrs:         addrs,
		session:       sess,
		namespace:     namespace,
		Router:        NewRouter(ServiceAuthenticate(auth)),
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

	s := &Server{
		listen:        listener,
		namespace:     namespace,
		Router:        NewRouter(service0),
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
//
// FIXME: Server registers the service to the namespace, it shall as
// well unregister it on Server.Terminate() *and* when requested by
// the service itself. TODO: need to add a Terminate() callback into
// the signal helper which is part of the stub interface... Shall the
// stub do the registration as well ? It has knowledge of the life
// cycle of the service.
//
// - call Server.namespace.Remove(serviceID)
// - call Server.Router.Remove()
// - call Terminate on all objects
// - signal WaitTerminate condition
//
// TODO:
// 1. fix the proxy generation issue
// 2. write a test auto with a service connecting a remote server
// 3. refactor this name registration into the stub
func (s *Server) NewService(name string, object Object) (Service, error) {

	service := NewService(object)

	// 1. reserve the name
	serviceID, err := s.namespace.Reserve(name)
	if err != nil {
		return nil, err
	}
	// 2. initialize the service
	err = service.Activate(s.Session(), serviceID)
	if err != nil {
		return nil, err
	}
	// 3. bring it online
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
	err := s.Router.Activate(s.session)
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
