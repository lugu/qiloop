package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/bus/util"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"log"
	"math/rand"
	gonet "net"
	"sync"
)

var ServiceNotFound error = errors.New("Service not found")
var ObjectNotFound error = errors.New("Object not found")
var ActionNotFound error = errors.New("Action not found")
var NotAuthenticated error = errors.New("Not authenticated")

type SignalUser struct {
	signalID  uint32
	messageID uint32
	context   *Context
	clientID  uint64
}

type BasicObject struct {
	meta      object.MetaObject
	signals   []SignalUser
	Wrapper   bus.Wrapper
	serviceID uint32
	objectID  uint32
}

func fullMetaObject(meta object.MetaObject) object.MetaObject {
	for i, method := range object.ObjectMetaObject.Methods {
		meta.Methods[i] = method
	}
	for i, signal := range object.ObjectMetaObject.Signals {
		meta.Signals[i] = signal
	}
	return meta
}

func NewObject(meta object.MetaObject) *BasicObject {
	var obj BasicObject
	obj.meta = fullMetaObject(meta)
	obj.signals = make([]SignalUser, 0)
	obj.Wrapper = make(map[uint32]bus.ActionWrapper)
	obj.Wrap(uint32(0x2), obj.wrapMetaObject)
	// obj.Wrapper[uint32(0x3)] = obj.Terminate
	// obj.Wrapper[uint32(0x5)] = obj.Property
	// obj.Wrapper[uint32(0x6)] = obj.SetProperty
	// obj.Wrapper[uint32(0x7)] = obj.Properties
	// obj.Wrapper[uint32(0x8)] = obj.RegisterEventWithSignature
	return &obj
}

func (o *BasicObject) Wrap(id uint32, fn bus.ActionWrapper) {
	o.Wrapper[id] = fn
}

func (o *BasicObject) AddSignalUser(signalID, messageID uint32, from *Context) uint64 {
	clientID := rand.Uint64()
	newUser := SignalUser{
		signalID,
		messageID,
		from,
		clientID,
	}
	o.signals = append(o.signals, newUser)
	return clientID
}

func (o *BasicObject) RemoveSignalUser(id uint64) error {
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
	clientID, err := basic.ReadUint64(buf)
	if err != nil {
		err = fmt.Errorf("cannot read client uid: %s", err)
		return util.ReplyError(from.EndPoint, msg, err)
	}
	messageID := msg.Header.ID
	clientID = o.AddSignalUser(signalID, messageID, from)
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
	err = o.RemoveSignalUser(clientID)
	if err != nil {
		return util.ReplyError(from.EndPoint, msg, err)
	}
	var out bytes.Buffer
	return o.reply(from, msg, out.Bytes())
}
func (o *BasicObject) MetaObject() (object.MetaObject, error) {
	return o.meta, nil
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
	meta, err := o.MetaObject()
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	err = object.WriteMetaObject(meta, &out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (o *BasicObject) UpdateSignal(signal uint32, value []byte) error {
	var ret error = nil
	for _, client := range o.signals {
		if client.signalID == signal {
			hdr := o.NewHeader(net.Event, signal, client.messageID)
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
	hdr := o.NewHeader(net.Reply, m.Header.Action, m.Header.ID)
	reply := net.NewMessage(hdr, response)
	return from.EndPoint.Send(reply)
}

func (o *BasicObject) handleDefault(from *Context, msg *net.Message) error {
	a, ok := o.Wrapper[msg.Header.Action]
	if !ok {
		return util.ReplyError(from.EndPoint, msg, ActionNotFound)
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

func (o *BasicObject) Terminate() {
	// TODO:
	// - call Server.directory.Remove()
	// - call Server.Router.Remove()
	// - call Terminate on all objects
	// - signal WaitTerminate condition
	panic("not yet implemented")
}

func (o *BasicObject) NewHeader(typ uint8, action, id uint32) net.Header {
	return net.NewHeader(typ, o.serviceID, o.objectID, action, id)
}

func (o *BasicObject) Activate(sess bus.Session, serviceID,
	objectID uint32) error {
	o.serviceID = serviceID
	o.objectID = objectID
	return nil
}

type Object interface {
	Receive(m *net.Message, from *Context) error
	Activate(sess bus.Session, serviceID, objectID uint32) error
}

type ServiceImpl struct {
	sync.RWMutex
	objects map[uint32]Object
}

func NewService(o Object) *ServiceImpl {
	return &ServiceImpl{
		objects: map[uint32]Object{
			1: o,
		},
	}
}

func (n *ServiceImpl) Add(o Object) (uint32, error) {
	var index uint32 = 0
	// assign the first object to the index 0. following objects will
	// be assigned random values.
	n.Lock()
	defer n.Unlock()
	if len(n.objects) != 0 {
		index = rand.Uint32()
	}
	n.objects[index] = o
	return index, nil
}

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

func (n *ServiceImpl) Remove(objectID uint32) error {
	n.Lock()
	defer n.Unlock()
	if _, ok := n.objects[objectID]; ok {
		delete(n.objects, objectID)
		return nil
	}
	return fmt.Errorf("Namespace: cannot remove object %d", objectID)
}

func (n *ServiceImpl) Dispatch(m *net.Message, from *Context) error {
	n.RLock()
	o, ok := n.objects[m.Header.Object]
	n.RUnlock()
	if ok {
		return o.Receive(m, from)
	}
	return util.ReplyError(from.EndPoint, m, ObjectNotFound)
}

func (n *ServiceImpl) Terminate() error {
	return fmt.Errorf("terminate not yet implemented")
}

// Router dispatch the incomming messages.
type Router struct {
	sync.RWMutex
	services map[uint32]*ServiceImpl
}

func NewRouter(authenticator Object) *Router {
	return &Router{
		services: map[uint32]*ServiceImpl{
			0: &ServiceImpl{
				objects: map[uint32]Object{
					0: authenticator,
				},
			},
		},
	}
}

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

func (r *Router) Terminate() error {
	r.RLock()
	defer r.RUnlock()
	var ret error = nil
	for serviceID, service := range r.services {
		err := service.Terminate()
		if err != nil && ret == nil {
			ret = fmt.Errorf("service %d terminate: %s",
				serviceID, err)
		}
	}
	return ret
}

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

func (r *Router) Remove(serviceID uint32) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.services[serviceID]; ok {
		delete(r.services, serviceID)
		return nil
	}
	return fmt.Errorf("Router: cannot remove service %d", serviceID)
}

func (r *Router) Dispatch(m *net.Message, from *Context) error {
	r.RLock()
	s, ok := r.services[m.Header.Service]
	r.RUnlock()
	if ok {
		return s.Dispatch(m, from)
	}
	return util.ReplyError(from.EndPoint, m, ServiceNotFound)
}

// Context represents the context of the request
type Context struct {
	EndPoint      net.EndPoint
	Authenticated bool
}

func NewContext(e net.EndPoint) *Context {
	return &Context{e, false}
}

// Firewall ensures an endpoint talks only to autorized services.
// Especially, it ensure authentication is passed.
func Firewall(m *net.Message, from *Context) error {
	if from.Authenticated == false && m.Header.Service != 0 {
		return NotAuthenticated
	}
	return nil
}

// Server listen from incomming connections, set-up the end points and
// forward the EndPoint to the dispatcher.
type Server struct {
	listen        gonet.Listener
	addrs         []string
	namespace     Namespace
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
// FIXME: update NewServer signature to add an authenticator
func NewServer(session *session.Session, addr string) (*Server, error) {
	l, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listen:        l,
		addrs:         []string{addr},
		session:       session,
		Router:        NewRouter(ServiceAuthenticate(Yes{})),
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
	namespace Namespace) (*Server, error) {

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

type Service interface {
	Terminate() error
}

func (s *Server) NewService(name string, object Object) (Service, error) {

	service := NewService(object)

	// 1. reserve the name
	serviceID, err := s.namespace.Reserve(name)
	if err != nil {
		return nil, err
	}
	// 2. initialize the service
	err = service.Activate(s.namespace.Session(s), serviceID)
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
	session := s.namespace.Session(s)
	err := s.Router.Activate(session)
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
	var ret error = nil
	s.contextsMutex.Lock()
	defer s.contextsMutex.Unlock()
	for context, _ := range s.contexts {
		err := context.EndPoint.Close()
		if err != nil && ret == nil {
			ret = err
		}
	}
	return ret
}

func (s *Server) WaitTerminate() chan error {
	return s.waitChan
}

func (s *Server) Stop() error {
	close(s.closeChan)
	err := s.listen.Close()
	s.stoppedWith(err)
	return err
}

func (s *Server) NewClient() bus.Client {
	ctl, srv := gonet.Pipe()
	s.handle(srv, true)
	return client.NewClient(net.NewEndPoint(ctl))
}

func (s *Server) Session() bus.Session {
	return s.namespace.Session(s)
}
