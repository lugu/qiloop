package fuzz

import (
	"bytes"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/object"
	"github.com/lugu/qiloop/type/value"
	"log"
)

var ServerURL string = "tcps://127.0.0.1:9503"

const serviceID = 0
const objectID = 0
const actionID = 8

func Fuzz(data []byte) int {
	endpoint, err := net.DialEndPoint(ServerURL)
	if err != nil {
		log.Fatalf("failed to contact %s: %s", ServerURL, err)
	}

	client0 := client.NewClient(endpoint)
	proxy0 := client.NewProxy(client0, object.MetaService0, serviceID, objectID)
	server0 := services.ServerProxy{proxy0}

	data, err0 := server0.CallID(actionID, data)

	// check response
	buf := bytes.NewBuffer(data)
	capability, err := session.ReadCapabilityMap(buf)
	if err == nil {
		statusValue, ok := capability[client.KeyState]
		if ok {
			status, ok := statusValue.(value.IntValue)
			if ok {
				switch uint32(status) {
				case client.StateDone:
					panic("password found")
				case client.StateContinue:
					panic("token renewal")
				}
			}
		}
	}

	// check if everything is still OK.
	endpoint2, err := net.DialEndPoint(ServerURL)
	if err != nil {
		panic("gateway has crashed")
	}

	err = client.Authenticate(endpoint2)
	endpoint2.Close()
	if err != nil {
		panic("gateway is broken")
	}

	if err0 == nil {
		return 1
	}
	return 0
}
