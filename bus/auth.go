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

func authenticateUser(endpoint net.EndPoint, user, token string) (CapabilityMap, error) {
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
	return authenticateCall(endpoint, permissions)
}

func authenticateCall(endpoint net.EndPoint, permissions CapabilityMap) (CapabilityMap, error) {
	const serviceID = 0
	const objectID = 0

	client0 := NewClient(endpoint)
	proxy0 := NewProxy(client0, object.MetaService0, serviceID, objectID)
	server0 := proxyServer0{
		proxy0,
	}
	return server0.Authenticate(permissions)
}

func authenticateContinue(endpoint net.EndPoint, user string, resp CapabilityMap) error {
	newTokenValue, ok := resp[KeyNewToken]
	if !ok {
		return fmt.Errorf("missing authentication new token")
	}
	newToken, ok := newTokenValue.(value.StringValue)
	if !ok {
		return fmt.Errorf("new token format error")
	}
	resp2, err := authenticateUser(endpoint, user, string(newToken))
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
		return token.WriteUserToken(user, string(newToken))
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
	resp, err := authenticateUser(endpoint, user, token)
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
		return authenticateContinue(endpoint, user, resp)
	case StateError:
		return fmt.Errorf("Authentication failed")
	default:
		return fmt.Errorf("invalid state type: %d", status)
	}
}

// Authenticate runs the authentication procedure of a given
// connection. If the clients fails to authenticate itself, it returns
// an error. In response the connection shall be closed.
func Authenticate(endpoint net.EndPoint) error {
	user, token := token.GetUserToken()
	return AuthenticateUser(endpoint, user, token)
}