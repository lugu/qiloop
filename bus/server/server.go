package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/client/services"
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
	obj.Wrapper[uint32(0x2)] = obj.wrapMetaObject
	// obj.Wrapper[uint32(0x3)] = obj.Terminate
	// obj.Wrapper[uint32(0x5)] = obj.Property
	// obj.Wrapper[uint32(0x6)] = obj.SetProperty
	// obj.Wrapper[uint32(0x7)] = obj.Properties
	// obj.Wrapper[uint32(0x8)] = obj.RegisterEventWithSignature
	return &obj
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
	panic("not yet implemented")
}

func (o *BasicObject) NewHeader(typ uint8, action, id uint32) net.Header {
	return net.NewHeader(typ, o.serviceID, o.objectID, action, id)
}

func (o *BasicObject) Activate(sess *session.Session, serviceID, objectID uint32) {
	o.serviceID = serviceID
	o.objectID = objectID
}

type Object interface {
	Receive(m *net.Message, from *Context) error
	Activate(sess *session.Session, serviceID, objectID uint32)
}

type Dispatcher interface {
	Wrap(id uint32, wrap bus.ActionWrapper)
	UpdateSignal(signal uint32, value []byte) error
}

// ObjectDispatcher implements both Object and Dispatcher
type ObjectDispatcher struct {
	wrapper bus.Wrapper
}

func (o *ObjectDispatcher) UpdateSignal(signal uint32, value []byte) error {
	panic("not available")
}

func (o *ObjectDispatcher) Wrap(id uint32, fn bus.ActionWrapper) {
	if o.wrapper == nil {
		o.wrapper = make(map[uint32]bus.ActionWrapper)
	}
	o.wrapper[id] = fn
}

func (o *ObjectDispatcher) Activate(sess *session.Session, serviceID, objectID uint32) {
}
func (o *ObjectDispatcher) Receive(m *net.Message, from *Context) error {
	if o.wrapper == nil {
		return util.ReplyError(from.EndPoint, m, ActionNotFound)
	}
	a, ok := o.wrapper[m.Header.Action]
	if !ok {
		return util.ReplyError(from.EndPoint, m, ActionNotFound)
	}
	response, err := a(m.Payload)

	if err != nil {
		return util.ReplyError(from.EndPoint, m, err)
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
			1: o,
		},
	}
}

func (n *ServiceImpl) Add(o Object) (uint32, error) {
	var index uint32 = 0
	// assign the first object to the index 0. following objects will
	// be assigned random values.
	if len(n.objects) != 0 {
		index = rand.Uint32()
	}
	n.mutex.Lock()
	n.objects[index] = o
	n.mutex.Unlock()
	return index, nil
}

func (s *ServiceImpl) Activate(sess *session.Session, serviceID uint32) {
	for objectID, obj := range s.objects {
		obj.Activate(sess, serviceID, objectID)
	}
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
	return util.ReplyError(from.EndPoint, m, ObjectNotFound)
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

func NewRouter(authenticator Object) *Router {
	return &Router{
		services: map[uint32]*ServiceImpl{
			0: &ServiceImpl{
				objects: map[uint32]Object{
					0: authenticator,
				},
			},
		},
		nextIndex: 1,
	}
}

func (r *Router) Activate(sess *session.Session) {
	for serviceID, service := range r.services {
		service.Activate(sess, serviceID)
	}
}

func (r *Router) Register(uid uint32, s *ServiceImpl) error {
	r.mutex.Lock()
	_, ok := r.services[uid]
	if ok {
		r.mutex.Unlock()
		return fmt.Errorf("service id already used: %d", uid)
	}
	r.services[uid] = s
	r.mutex.Unlock()
	return nil
}

func (r *Router) Add(s *ServiceImpl) (uint32, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for {
		_, ok := r.services[r.nextIndex]
		if !ok {
			break
		}
		r.nextIndex++
	}
	uid := r.nextIndex
	r.nextIndex++
	r.services[uid] = s
	return uid, nil
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
	return util.ReplyError(from.EndPoint, m, ServiceNotFound)
}

// Context represents the context of the request
type Context struct {
	EndPoint      net.EndPoint
	Authenticated bool
}

func NewContext(c gonet.Conn) *Context {
	return &Context{net.NewEndPoint(c), false}
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
	session       session.Session
	Router        *Router
	contexts      map[*Context]bool
	contextsMutex sync.Mutex
}

func NewServer(session session.Session, addr string) (*Server, error) {
	l, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		listen:  l,
		addrs:   []string{addr},
		session: session,
		// FIXME: update NewServer signature to add an
		// authenticator
		Router:        NewRouter(NewServiceAuthenticate(make(map[string]string))),
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
		SessionId: "", // TODO
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
	service.Activate(&s.session, uid)
	return service, nil
}

func StandAloneServer(l gonet.Listener, r *Router) *Server {
	s := make(map[*Context]bool)
	r.Activate(nil)
	return &Server{
		listen:        l,
		Router:        r,
		contexts:      s,
		contextsMutex: sync.Mutex{},
	}
}

func (s *Server) handle(c gonet.Conn) error {
	context := NewContext(c)
	filter := func(hdr *net.Header) (matched bool, keep bool) {
		return true, true
	}
	consumer := func(msg *net.Message) error {
		err := Firewall(msg, context)
		if err != nil {
			log.Printf("missing authentication from %s:%s: %v",
				c.RemoteAddr().Network(), c.RemoteAddr().String(),
				msg.Header)
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

	s.contextsMutex.Lock()
	s.contexts[context] = true
	s.contextsMutex.Unlock()
	context.EndPoint.AddHandler(filter, consumer, closer)
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
