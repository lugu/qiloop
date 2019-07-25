package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

var (
	infos = make([]services.ServiceInfo, 0)
	metas = make([]object.MetaObject, 0)
)

func getObject(sess bus.Session, info services.ServiceInfo) bus.ObjectProxy {
	proxy, err := sess.Proxy(info.Name, 1)
	if err != nil {
		log.Fatalf("connect service (%s): %s", info.Name, err)
	}
	return bus.MakeObject(proxy)
}

func print(event bus.EventTrace, info *services.ServiceInfo,
	meta *object.MetaObject) {

	var typ string = "unknown"
	switch event.Kind {
	case int32(net.Call):
		typ = "call "
	case int32(net.Reply):
		typ = "reply"
	}
	var action = "unknown"
	action, err := meta.ActionName(event.SlotId)
	if err != nil {
		action = "unknown"
	}
	var size int = -1
	var sig = "unknown"
	var data = []byte{}
	var buf bytes.Buffer
	err = event.Arguments.Write(&buf)
	if err == nil {
		sig, err = basic.ReadString(&buf)
		if err == nil {
			data = buf.Bytes()
			size = len(data)
		}
	}

	fmt.Printf("[%s %4d bytes] %s.%s: %s: %v\n", typ, size, info.Name,
		action, sig, data)
}

func trace(serverURL, serviceName string) {

	sess, err := session.NewSession(serverURL)
	if err != nil {
		panic(err)
	}

	proxies := services.Services(sess)

	directory, err := proxies.ServiceDirectory()
	if err != nil {
		panic(err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})

	for _, info := range serviceList {

		if serviceName != "" && serviceName != info.Name {
			continue
		}

		go func(info services.ServiceInfo) {
			obj := getObject(sess, info)
			err = obj.EnableTrace(true)
			if err != nil {
				log.Fatalf("Failed to start traces: %s.", err)
			}
			defer obj.EnableTrace(false)

			cancel, trace, err := obj.SubscribeTraceObject()
			if err != nil {
				log.Fatalf("Failed to stop stats: %s.", err)
			}
			defer cancel()

			meta, err := obj.MetaObject(1)
			if err != nil {
				log.Fatalf("%s: MetaObject: %s.", info.Name, err)
			}

			for {
				select {
				case event, ok := <-trace:
					if !ok {
						return
					}
					print(event, &info, &meta)
				case <-stop:
					return
				}
			}
		}(info)
	}
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)

	<-signalChannel
	close(stop)
}
