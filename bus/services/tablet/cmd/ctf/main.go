package main

import (
	"flag"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	tablet "github.com/lugu/qiloop/bus/services/tablet"
	"github.com/lugu/qiloop/bus/session/basic"
	"log"
)

func findTabletService(addr string) (tablet.ALTabletService, error) {
	objectID := uint32(1)

	// FIXME: broken: shall the the actionID of the call need to
	// querry the MetaObject to resolve the actionID
	actionID := uint32(143)

	serviceID, err := findTabletServiceID(addr)
	if err != nil {
		return nil, err
	}

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		log.Fatalf("failed to contact %s: %s", addr, err)
	}
	session := basic.NewSession(endpoint, serviceID, objectID, actionID)
	return tablet.NewALTabletService(session, objectID)

}

func findTabletServiceID(addr string) (uint32, error) {
	for id := uint32(80); id < uint32(150); id++ {
		if test(addr, id) {
			return id, nil
		}
	}
	return 0, fmt.Errorf("ALTabletService not found")
}

func test(addr string, serviceID uint32) (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			ret = false
		}
	}()

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		log.Fatalf("failed to contact %s: %s", addr, err)
	}

	objectID := uint32(1)
	actionID := uint32(143) // getWifiMacAddress
	session := basic.NewSession(endpoint, serviceID, objectID, actionID)

	tablet, err := tablet.NewALTabletService(session, objectID)
	if err != nil {
		panic(err)
	}

	_, err = tablet.GetWifiMacAddress()
	if err != nil {
		panic(err)
	}
	return true
}

func main() {
	var serverURL = flag.String("qi-url", "tcp://127.0.0.1:9559", "server URL")
	flag.Parse()

	tablet, err := findTabletService(*serverURL)
	if err != nil {
		panic(err)
	}
	ret, err := tablet.RobotIp()
	if err != nil {
		panic(err)
	}
	println(ret)
}
