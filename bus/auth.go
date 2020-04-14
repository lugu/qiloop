package bus

import (
	"bytes"
	"fmt"
	"time"

	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/session/token"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
)

const (
	// KeyState is the key of the authetication state in the capability map
	KeyState = "__qi_auth_state"
	// KeyUser is the key for the user name in the capability map
	KeyUser = "auth_user"
	// KeyToken is the key for the user token in the capability map
	KeyToken = "auth_token"
	// KeyNewToken is the key with the new token sent in the capability map
	KeyNewToken = "auth_newToken"

	// StateError indicates an authentication failure.
	StateError uint32 = 1
	// StateContinue indicates a new token generation procedure.
	StateContinue uint32 = 2
	// StateDone indicates an authentication success.
	StateDone uint32 = 3
)

// CapabilityMap is a data structure exchanged between the server and
// the client during the authentication procedure. It is used as a
// negociation medium between the client an the server to decides
// which features are supported and to verify the client credentials.
type CapabilityMap map[string]value.Value

// Authenticated return true if the authentication procedure has been
// completed with success.
func (c CapabilityMap) Authenticated() bool {
	statusValue, ok := c[KeyState]
	if !ok {
		return false
	}
	status, ok := statusValue.(value.UintValue)
	if !ok {
		status2, ok := statusValue.(value.IntValue)
		if !ok {
			return false
		}
		status = value.UintValue(uint32(status2.Value()))
	}
	if uint32(status) == StateDone {
		return true
	}
	return false
}

// SetAuthenticated force the done status in the capability map.
func (c CapabilityMap) SetAuthenticated() {
	c[KeyState] = value.Uint(StateDone)
}

func (c CapabilityMap) hasBool(val string, def bool) bool {
	b, ok := c[val]
	if !ok {
		return def
	}
	return b == value.Bool(true)
}

func (c CapabilityMap) MetaObjectCache() bool {
	return c.hasBool("MetaObjectCache", false)
}

func (c CapabilityMap) ObjectPtrUID() bool {
	return c.hasBool("ObjectPtrUID", false)
}

func (c CapabilityMap) RemoteCancelableCalls() bool {
	return c.hasBool("ObjectPtrUID", false)
}

// DefaultCap returns a CapabilityMap based on default settings.
func DefaultCap() CapabilityMap {
	return CapabilityMap{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(false),
		"RemoteCancelableCalls": value.Bool(false),
		"ObjectPtrUID":          value.Bool(false),
	}
}

// ClientCap returns a CapabilityMap based on the default settings.
// FIXME: get ride of the special case
func ClientCap(user, token string) CapabilityMap {
	permissions := DefaultCap()
	if user != "" {
		permissions[KeyUser] = value.String(user)
	}
	if token != "" {
		permissions[KeyToken] = value.String(token)
	}
	return permissions
}

func authenticateCall(endpoint net.EndPoint, permissions CapabilityMap) (CapabilityMap, error) {

	// 1. prepare to receive a capability message
	filter := func(hdr *net.Header) (matched bool, keep bool) {
		if hdr.Type == net.Capability {
			return true, false
		}
		return false, true
	}
	reply := make(chan *net.Message, 1)
	errors := make(chan error, 1)
	closer := func(err error) {
		if err != nil {
			errors <- err
		}
	}
	id := endpoint.MakeHandler(filter, reply, closer)
	defer endpoint.RemoveHandler(id)

	// 2. normal call to authenticate
	cache := NewCache(endpoint)
	cache.AddService("ServiceZero", 0, object.MetaService0)
	service0, err := ServiceServer(cache)
	m, err := service0.Authenticate(permissions)
	if err == nil {
		return m, nil
	}

	// 3. on error, send the capability map
	var out bytes.Buffer
	err = WriteCapabilityMap(permissions, &out)
	if err != nil {
		return m, err
	}
	hdr := net.NewHeader(net.Capability, 0, 0, 0, 2)
	err = endpoint.Send(net.NewMessage(hdr, out.Bytes()))
	if err != nil {
		return m, err
	}

	// 4. wait for the capability map
	timer := time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
		return m, fmt.Errorf("missing capability: timeout")
	case err := <-errors:
		return m, err
	case msg, ok := <-reply:
		if !ok {
			return m, fmt.Errorf("Remote connection closed")
		}
		buf := bytes.NewBuffer(msg.Payload)
		m, err := ReadCapabilityMap(buf)
		if err != nil {
			return m, err
		}
		m[KeyState] = value.Uint(StateDone)
		return m, nil
	}
}

// authenticateContinue update the token with the server provided
// token. The prefered CapabilityMap is updated with an entry
// KeyNewToken containing the new token.
func authenticateContinue(endpoint net.EndPoint, prefered, resp CapabilityMap) error {
	newTokenValue, ok := resp[KeyNewToken]
	if !ok {
		return fmt.Errorf("missing authentication new token")
	}
	newToken, ok := newTokenValue.(value.StringValue)
	if !ok {
		return fmt.Errorf("new token format error")
	}
	oldToken := prefered[KeyToken]
	prefered[KeyToken] = value.String(string(newToken))
	resp2, err := authenticateCall(endpoint, prefered)
	if err != nil {
		return fmt.Errorf("new token authentication failed: %s", err)
	}
	statusValue, ok := resp2[KeyState]
	if !ok {
		return fmt.Errorf("missing authentication state")
	}
	status, ok := statusValue.(value.UintValue)
	if !ok {
		status2, ok := statusValue.(value.IntValue)
		if !ok {
			return fmt.Errorf("authentication status error (%#v)",
				statusValue)
		}
		status = value.UintValue(uint32(status2.Value()))
	}
	switch uint32(status) {
	case StateDone:
		prefered[KeyNewToken] = prefered[KeyToken]
		prefered[KeyToken] = oldToken
		return nil
	case StateContinue:
		return fmt.Errorf("new token authentication dropped")
	case StateError:
		return fmt.Errorf("new token authentication failed")
	default:
		return fmt.Errorf("invalid state type: %d", status)
	}
}

// AuthenticateUser trigger an authentication procedure using the user
// and token. If user or token is empty, it is not included in the
// request.
func AuthenticateUser(endpoint net.EndPoint, user, token string) error {
	return Authentication(endpoint, ClientCap(user, token))
}

// Authenticate runs the authentication procedure of a given
// connection. If the clients fails to authenticate itself, it returns
// an error. In response the connection shall be closed.
func Authenticate(endpoint net.EndPoint) error {
	user, token := token.GetUserToken()
	return AuthenticateUser(endpoint, user, token)
}

// Authentication runs the authentication procedure. The prefered
// CapabilityMap is updated with the negociated capabilities.
// If a new token was issued, it is store in the KeyNewToken entry of
// the prefered capapbility map.
func Authentication(endpoint net.EndPoint, prefered CapabilityMap) error {

	resp, err := authenticateCall(endpoint, prefered)
	if err != nil {
		return fmt.Errorf("authentication failed: %s", err)
	}
	statusValue, ok := resp[KeyState]
	if !ok {
		return fmt.Errorf("missing authentication state")
	}
	status, ok := statusValue.(value.UintValue)
	if !ok {
		status2, ok := statusValue.(value.IntValue)
		if !ok {
			return fmt.Errorf("authentication status error (%#v)",
				statusValue)
		}
		status = value.UintValue(uint32(status2.Value()))
	}
	switch uint32(status) {
	case StateDone:
		return nil
	case StateContinue:
		return authenticateContinue(endpoint, prefered, resp)
	case StateError:
		return fmt.Errorf("Authentication failed")
	default:
		return fmt.Errorf("invalid state type: %d", status)
	}
}
