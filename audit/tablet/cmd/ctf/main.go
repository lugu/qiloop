package main

import (
	"flag"
	"fmt"
	"github.com/lugu/qiloop/bus/client"
	tablet "github.com/lugu/qiloop/bus/services/tablet"
	"log"
)

func findTabletService(addr string) (tablet.ALTabletService, error) {

	serviceID, err := findTabletServiceID(addr)
	if err != nil {
		return nil, err
	}

	cache, err := client.NewCachedSession(addr)
	if err != nil {
		log.Fatalf("authentication failed: %s", err)
	}

	err = cache.Lookup("ALTabletService", serviceID)
	if err != nil {
		log.Fatalf("failed to query meta object of ALTabletService: %s", err)
	}

	objectID := uint32(1)
	return tablet.NewALTabletService(cache, objectID)

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

	cache, err := client.NewCachedSession(addr)
	if err != nil {
		log.Fatalf("authentication failed: %s", err)
	}

	err = cache.Lookup("ALTabletService", serviceID)
	if err != nil {
		log.Fatalf("failed to query meta object of ALTabletService: %s", err)
	}

	objectID := uint32(1)
	tablet, err := tablet.NewALTabletService(cache, objectID)
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
