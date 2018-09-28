package fuzz

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"log"
)

var serverURL string = "tcps://10.0.165.120:9503"

var user string = "tablet"
var token string = "0ee88b39-831d-4f36-a813-3cf92a84810d"

const serviceID = 0
const objectID = 0
const actionID = 8

func ParseCapabilityMap(data []byte) (session.CapabilityMap, error) {
	buf := bytes.NewBuffer(data)
	ret, err := func() (m map[string]value.Value, err error) {
		size, err := basic.ReadUint32(buf)
		if err != nil {
			return m, fmt.Errorf("failed to read map size: %s", err)
		}
		m = make(map[string]value.Value, size)
		for i := 0; i < int(size); i++ {
			k, err := basic.ReadString(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map key: %s", err)
			}
			v, err := value.NewValue(buf)
			if err != nil {
				return m, fmt.Errorf("failed to read map value: %s", err)
			}
			m[k] = v
		}
		return m, nil
	}()
	if err != nil {
		return ret, fmt.Errorf("failed to parse authenticate response: %s", err)
	}
	return ret, nil
}

func Fuzz(data []byte) int {
	endpoint, err := net.DialEndPoint(serverURL)
	if err != nil {
		log.Fatalf("failed to contact %s: %s", serverURL, err)
	}

	client0, _ := session.NewClient(endpoint)
	proxy0 := session.NewProxy(client0, object.MetaService0, serviceID, objectID)
	server0 := services.ServerProxy{proxy0}

	data, err0 := server0.CallID(actionID, data)

	// check response
	capability, err := ParseCapabilityMap(data)
	if err == nil {
		statusValue, ok := capability[session.KeyState]
		if ok {
			status, ok := statusValue.(value.IntValue)
			if ok {
				switch uint32(status) {
				case session.StateDone:
					panic("password found")
				case session.StateContinue:
					panic("token renewal")
				}
			}
		}
	}

	// check if everything is still OK.
	endpoint2, err := net.DialEndPoint(serverURL)
	if err != nil {
		panic("gateway has crashed")
	}

	err = session.Authenticate(endpoint2)
	endpoint2.Close()
	if err != nil {
		panic("gateway is broken")
	}

	if err0 == nil {
		return 1
	}
	return 0
}
