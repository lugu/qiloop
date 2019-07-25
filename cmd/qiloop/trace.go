package main

import (
	"bytes"
	"fmt"
	"log"
	"reflect"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

var (
	traces = make([]chan bus.EventTrace, 0)
	infos  = make([]services.ServiceInfo, 0)
	metas  = make([]object.MetaObject, 0)
)

func getObject(sess bus.Session, info services.ServiceInfo) bus.ObjectProxy {
	proxy, err := sess.Proxy(info.Name, 1)
	if err != nil {
		log.Fatalf("connect service (%s): %s", info.Name, err)
	}
	return bus.MakeObject(proxy)
}

func print(chosen int, events []bus.EventTrace) {
	info := infos[chosen]
	meta := metas[chosen]

	for _, event := range events {
		var typ string = "unknown"
		switch event.Kind {
		case int32(net.Call):
			typ = "call"
		case int32(net.Reply):
			typ = "reply"
		}
		var action = "unknown"
		action, err := meta.ActionName(event.SlotId)
		if err != nil {
			action = "unknown"
		}
		var size int = -1
		var buf bytes.Buffer
		err = event.Arguments.Write(&buf)
		if err == nil {
			// read the signature back
			_, err := basic.ReadString(&buf)
			if err == nil {
				data := buf.Bytes()
				size = len(data)
			}
		}

		fmt.Printf("[%s] %s.%s (%d bytes): %#v\n", typ, info.Name,
			action, size, event.Arguments)
	}
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

	for _, info := range serviceList {

		if serviceName != "" && serviceName != info.Name {
			continue
		}

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

		metas = append(metas, meta)
		traces = append(traces, trace)
		infos = append(infos, info)
	}

	cases := make([]reflect.SelectCase, len(traces))
	for i, trace := range traces {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(trace)}
	}
	for {
		chosen, _, ok := reflect.Select(cases)
		if !ok {
			return
		}
		ch := traces[chosen]
		events := make([]bus.EventTrace, 0)
	loop:
		for {
			select {

			case event := <-ch:
				events = append(events, event)
			default:
				if len(events) == 0 {
					log.Printf("missing events (%d)", chosen)
					break loop
				}
				go print(chosen, events)
				break loop
			}
		}
	}
}
