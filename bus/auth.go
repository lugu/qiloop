package bus

import (
	"fmt"

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

// PreferedCap returns a CapabilityMap based on default settings. If
// user or token is empty, it is ignored.
func PreferedCap(user, token string) CapabilityMap {
	permissions := CapabilityMap{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(false),
		"RemoteCancelableCalls": value.Bool(false),
		"ObjectPtrUID":          value.Bool(false),
	}
	if user != "" {
		permissions[KeyUser] = value.String(user)
	}
	if token != "" {
		permissions[KeyToken] = value.String(token)
	}
	return permissions
}

func authenticateCall(endpoint net.EndPoint, permissions CapabilityMap) (CapabilityMap, error) {

	cache := NewCache(endpoint)
	cache.AddService("ServiceZero", 0, object.MetaService0)
	proxies := Services(cache)
	service0, err := proxies.ServiceServer()
	if err != nil {
		return CapabilityMap{}, err
	}

	return service0.Authenticate(permissions)
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
	return Authentication(endpoint, PreferedCap(user, token))
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
