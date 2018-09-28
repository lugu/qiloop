package session

import (
	"bufio"
	"fmt"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
	"sync"
)

const (
	KeyState    string = "__qi_auth_state"
	KeyUser     string = "auth_user"
	KeyToken    string = "auth_token"
	KeyNewToken string = "auth_newToken"

	StateError    uint32 = 1
	StateContinue uint32 = 2
	StateDone     uint32 = 3
)

type CapabilityMap map[string]value.Value

func SaveNewUserToken(login string, token string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	var flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(usr.HomeDir+"/.qi-auth.conf.new", flag, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open auth file: %s", err)
	}
	_, err = file.WriteString(login + "\n" + token)
	if err != nil {
		return fmt.Errorf("Failed to write auth file: %s", err)
	}
	return nil
}

func GetUserToken() (string, string) {
	usr, err := user.Current()
	if err != nil {
		return "", ""
	}
	file, err := os.Open(usr.HomeDir + "/.qi-auth.conf")
	if err != nil {
		return "", ""
	}
	defer file.Close()
	r := bufio.NewReader(file)
	user, err := r.ReadString('\n')
	if err != nil {
		return "", ""
	}
	pwd, err := r.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", ""
	}
	// FIXME: don't trim space from pwd
	return strings.TrimSpace(user), strings.TrimSpace(pwd)
}

func AuthenticateUser(endpoint net.EndPoint, user, token string) (CapabilityMap, error) {
	permissions := CapabilityMap{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	if user != "" {
		permissions[KeyUser] = value.String(user)
	}
	if token != "" {
		permissions[KeyToken] = value.String(token)
	}
	return AuthenticateCall(endpoint, permissions)

}

func AuthenticateCall(endpoint net.EndPoint, permissions CapabilityMap) (CapabilityMap, error) {
	const serviceID = 0
	const objectID = 0

	client0 := newClient(endpoint)
	proxy0 := NewProxy(client0, object.MetaService0, serviceID, objectID)
	server0 := services.ServerProxy{proxy0}
	return server0.Authenticate(permissions)
}

func AuthenticateContinue(endpoint net.EndPoint, user string, resp CapabilityMap) error {
	newTokenValue, ok := resp[KeyNewToken]
	if !ok {
		return fmt.Errorf("missing authentication new token")
	}
	newToken, ok := newTokenValue.(value.StringValue)
	if !ok {
		return fmt.Errorf("new token format error")
	}
	resp, err := AuthenticateUser(endpoint, user, string(newToken))
	if err != nil {
		return fmt.Errorf("authentication failed: %s", err)
	}
	statusValue, ok := resp[KeyState]
	if !ok {
		return fmt.Errorf("missing authentication state")
	}
	status, ok := statusValue.(value.IntValue)
	if !ok {
		return fmt.Errorf("authentication state format error")
	}
	switch uint32(status) {
	case StateDone:
		return SaveNewUserToken(user, string(newToken))
	case StateContinue:
		return fmt.Errorf("Authentication dropped")
	case StateError:
		return fmt.Errorf("Authentication failed")
	default:
		return fmt.Errorf("invalid state type: %d", status)
	}
}

func Authenticate(endpoint net.EndPoint) error {

	user, token := GetUserToken()
	resp, err := AuthenticateUser(endpoint, user, token)
	if err != nil {
		return fmt.Errorf("authentication failed: %s", err)
	}
	statusValue, ok := resp[KeyState]
	if !ok {
		return fmt.Errorf("missing authentication state")
	}
	status, ok := statusValue.(value.IntValue)
	if !ok {
		return fmt.Errorf("authentication status error")
	}
	switch uint32(status) {
	case StateDone:
		return nil
	case StateContinue:
		return AuthenticateContinue(endpoint, user, resp)
	case StateError:
		return fmt.Errorf("Authentication failed")
	default:
		return fmt.Errorf("invalid state type: %d", status)
	}
}

// Session implements the Session interface. It is an
// implementation of Session. It does not update the list of services
// and returns clients.

type Session struct {
	serviceList      []services.ServiceInfo
	serviceListMutex sync.Mutex
	Directory        services.ServiceDirectory
	cancel           chan int
	added            chan struct {
		P0 uint32
		P1 string
	}
	removed chan struct {
		P0 uint32
		P1 string
	}
}

func newEndPoint(info services.ServiceInfo) (endpoint net.EndPoint, err error) {
	if len(info.Endpoints) == 0 {
		return endpoint, fmt.Errorf("missing address for service %s", info.Name)
	}
	endpoint, err = net.DialEndPoint(info.Endpoints[0])
	if err != nil {
		return endpoint, fmt.Errorf("%s: connection failed (%s) : %s", info.Name, info.Endpoints[0], err)
	}
	if err = Authenticate(endpoint); err != nil {
		return endpoint, fmt.Errorf("authentication error (%s): %s", info.Name, err)
	}
	return endpoint, nil
}

func newObject(info services.ServiceInfo, ref object.ObjectReference) (object.Object, error) {
	endpoint, err := newEndPoint(info)
	if err != nil {
		return nil, fmt.Errorf("object connection error (%s): %s", info.Name, err)
	}
	proxy := newProxy(endpoint, ref.MetaObject, ref.ServiceID, ref.ObjectID)
	return &services.ObjectProxy{proxy}, nil
}

func newService(info services.ServiceInfo, objectID uint32) (p bus.Proxy, err error) {
	endpoint, err := newEndPoint(info)
	if err != nil {
		return nil, fmt.Errorf("service connection error (%s): %s", info.Name, err)
	}
	proxy, err := metaProxy(endpoint, info.ServiceId, objectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service meta object (%s): %s", info.Name, err)
	}
	return proxy, nil
}

// FIXME: objectID does not seems needed
func (s *Session) Proxy(name string, objectID uint32) (p bus.Proxy, err error) {
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()

	for _, service := range s.serviceList {
		if service.Name == name {
			return newService(service, objectID)
		}
	}
	return p, fmt.Errorf("service not found: %s", name)
}

func (s *Session) Object(ref object.ObjectReference) (o object.Object, err error) {
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	for _, service := range s.serviceList {
		if service.ServiceId == ref.ServiceID {
			return newObject(service, ref)
		}
	}
	return o, fmt.Errorf("Not yet implemented")
}

func newProxy(e net.EndPoint, meta object.MetaObject, serviceID, objectID uint32) bus.Proxy {
	return NewProxy(newClient(e), meta, serviceID, objectID)
}

// metaProxy is to create proxies to the directory and server
// services needed for a session.
func metaProxy(e net.EndPoint, serviceID, objectID uint32) (p bus.Proxy, err error) {
	client := newClient(e)
	meta, err := bus.MetaObject(client, serviceID, objectID)
	if err != nil {
		return p, fmt.Errorf("Can not reach metaObject: %s", err)
	}
	return NewProxy(client, meta, serviceID, objectID), nil
}

func NewSession(addr string) (bus.Session, error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to contact %s: %s", addr, err)
	}
	if err = Authenticate(endpoint); err != nil {
		endpoint.Close()
		return nil, fmt.Errorf("authenitcation failed: %s", err)
	}

	proxy, err := metaProxy(endpoint, 1, 1)
	if err != nil {
		endpoint.Close()
		return nil, fmt.Errorf("failed to get directory meta object: %s", err)
	}
	s := new(Session)
	s.Directory = &services.ServiceDirectoryProxy{proxy}

	s.serviceList, err = s.Directory.Services()
	if err != nil {
		endpoint.Close()
		return nil, fmt.Errorf("failed to list services: %s", err)
	}
	s.cancel = make(chan int)
	s.removed, err = s.Directory.SignalServiceRemoved(s.cancel)
	if err != nil {
		endpoint.Close()
		return nil, fmt.Errorf("failed to subscribe remove signal: %s", err)
	}
	s.added, err = s.Directory.SignalServiceAdded(s.cancel)
	if err != nil {
		endpoint.Close()
		return nil, fmt.Errorf("failed to subscribe added signal: %s", err)
	}
	return s, nil
}

func (s *Session) updateServiceList() {
	var err error
	s.serviceListMutex.Lock()
	defer s.serviceListMutex.Unlock()
	s.serviceList, err = s.Directory.Services()
	if err != nil {
		log.Printf("error: failed to update service directory list: %s", err)
		log.Printf("error: closing session.")
		if err := s.Destroy(); err != nil {
			log.Printf("error: session destruction: %s", err)
		}
	}
}

func (s *Session) Destroy() error {
	// cancel both add and remove services
	s.cancel <- 1
	s.cancel <- 1
	return s.Directory.Disconnect()
}

func (s *Session) updateLoop() {
	for {
		select {
		case _, ok := <-s.removed:
			if !ok {
				return
			}
			s.updateServiceList()
		case _, ok := <-s.added:
			if !ok {
				return
			}
			s.updateServiceList()
		}
	}
}
