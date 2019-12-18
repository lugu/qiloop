package bus

import (
	"errors"
	"log"
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

	s := &server{
		listen:        listener,
		namespace:     namespace,
		contexts:      make(map[Channel]bool),
		contextsMutex: sync.Mutex{},
		closeChan:     make(chan int, 1),
		waitChan:      make(chan error, 1),
	}
	s.Router = NewRouter(service0, namespace, s.Session())
	_, err := s.NewService("ServiceDirectory", service1)
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

	s := &server{
		listen:        listener,
		namespace:     namespace,
		contexts:      make(map[Channel]bool),
		contextsMutex: sync.Mutex{},
		closeChan:     make(chan int, 1),
		waitChan:      make(chan error, 1),
	}
	s.Router = NewRouter(service0, namespace, s.Session())
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

// ServiceReceiver represents the service used by the Server.
type ServiceReceiver interface {
	Service
	Receiver
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

	// 1. reserve the name
	serviceID, err := s.namespace.Reserve(name)
	if err != nil {
		return nil, err
	}

	// 2. activate the service
	activation := serviceActivation(s.Router, session, serviceID)
	service, err := NewService(object, activation)
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
	consumer := make(chan *net.Message, 10)
	go func() {
		for msg := range consumer {
			err := firewall(msg, context)
			if err != nil {
				log.Printf("missing authentication from %s: %#v",
					context.EndPoint().String(), msg.Header)
				context.SendError(msg, err)
				stream.Close()
				return
			}
			err = s.Router.Receive(msg, context)
			if err != nil {
				log.Printf("error %v: %s", msg.Header, err)
			}
		}
	}()
	closer := func(err error) {
		s.contextsMutex.Lock()
		defer s.contextsMutex.Unlock()
		if _, ok := s.contexts[context]; ok {
			delete(s.contexts, context)
		}
	}
	finalize := func(e net.EndPoint) {
		context.endpoint = e
		e.MakeHandler(filter, consumer, closer)
		s.contextsMutex.Lock()
		s.contexts[context] = true
		s.contextsMutex.Unlock()
	}
	net.EndPointFinalizer(stream, finalize)
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
