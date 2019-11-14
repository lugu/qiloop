package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/aybabtme/rgbterm"
	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/services"
	"github.com/lugu/qiloop/bus/session"
	"github.com/lugu/qiloop/type/basic"
	"github.com/lugu/qiloop/type/object"
)

var (
	infos   = make([]services.ServiceInfo, 0)
	metas   = make([]object.MetaObject, 0)
	colored = true
)

func init() {
	stat, err := os.Stdout.Stat()
	if err != nil {
		colored = false
	}
	mode := stat.Mode()
	if (mode&os.ModeDevice == 0) || (mode&os.ModeCharDevice == 0) {
		colored = false
	}
}

func getObject(sess bus.Session, info services.ServiceInfo,
	objectID uint32) (bus.ObjectProxy, error) {

	proxy, err := sess.Proxy(info.Name, objectID)
	if err != nil {
		return nil, fmt.Errorf("connect service (%s): %s", info.Name, err)
	}
	return bus.MakeObject(proxy), nil
}

func msgType(typ uint8) (color, label string) {
	switch typ {
	case net.Call:
		return "{#cccccc}", "call "
	case net.Reply:
		return "{#cccccc}", "reply"
	case net.Error:
		return "{#ff0000}", "error"
	case net.Post:
		return "{#0088ff}", "post "
	case net.Event:
		return "{#0088ff}", "event"
	case net.Capability:
		return "{#ff8800}", "cap  "
	case net.Cancel:
		return "{#ff8800}", "cancl"
	case net.Cancelled:
		return "{#ff8800}", "cnled"
	default:
		return "{#ff0000}", "unkno"
	}
}

func printEvent(e bus.EventTrace, info *services.ServiceInfo,
	meta *object.MetaObject) {

	action, err := meta.ActionName(e.SlotId)
	if err != nil {
		action = fmt.Sprintf("unknown (%d)", e.SlotId)
	}
	ts := time.Unix(e.Timestamp.Tv_sec,
		e.Timestamp.Tv_usec*1000).Format("2006/01/02 15:04:05.0000000")

	var size = -1
	var sig = "unknown"
	var data = []byte{}
	var buf bytes.Buffer
	err = e.Arguments.Write(&buf)
	if err == nil {
		sig, err = basic.ReadString(&buf)
		if err == nil {
			data = buf.Bytes()
			size = len(data)
		}
	}

	color, typ := msgType(uint8(e.Kind))
	nocolor := "{}"
	out := rgbterm.ColorOut
	if !colored {
		color = ""
		nocolor = ""
		out = os.Stdout
	}

	fmt.Fprintf(out, "%s%s [id %8d] [%s %4d bytes] %s.%s: %s%s\n",
		color, ts, e.Id, typ, size, info.Name, action, sig, nocolor)
}

func trace(serverURL, serviceName string, objectID uint32) {

	stop := make(chan struct{})

	sess, err := session.NewSession(serverURL)
	if err != nil {
		log.Fatalf("%s: %s", serverURL, err)
	}

	proxies := services.Services(sess)

	closer := func(err error) {
		log.Printf("Session terminated (%s)", err.Error())
		close(stop)
	}

	directory, err := proxies.ServiceDirectory(closer)
	if err != nil {
		log.Fatalf("directory: %s", err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		log.Fatalf("services: %s", err)
	}

	serviceID, err := strconv.Atoi(serviceName)
	if err != nil {
		serviceID = -1
	} else {
		serviceName = ""
	}

	var group sync.WaitGroup

	for _, info := range serviceList {

		if serviceID != -1 && uint32(serviceID) != info.ServiceId {
			continue
		} else if serviceName != "" && serviceName != info.Name {
			continue
		}

		obj, err := getObject(sess, info, objectID)
		if err != nil {
			log.Printf("cannot trace %s: %s", info.Name, err)
			continue
		}

		group.Add(1)

		go func(info services.ServiceInfo, obj bus.ObjectProxy) {

			defer group.Done()

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

			meta, err := obj.MetaObject(objectID)
			if err != nil {
				log.Fatalf("%s: MetaObject: %s.", info.Name, err)
			}

			for {
				select {
				case event, ok := <-trace:
					if !ok {
						return
					}
					printEvent(event, &info, &meta)
				case <-stop:
					return
				}

			}
		}(info, obj)
	}

	go func() {
		signalChannel := make(chan os.Signal, 2)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT)
		<-signalChannel
		close(stop)
	}()

	group.Wait()
}
