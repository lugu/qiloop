package main

import (
	"fmt"
	"github.com/lugu/qiloop/meta/proxy"
	"github.com/lugu/qiloop/meta/stage2"
	"github.com/lugu/qiloop/net"
	"github.com/lugu/qiloop/object"
	"github.com/lugu/qiloop/session"
	"github.com/lugu/qiloop/value"
	"io"
	"log"
	"os"
	"strings"
)

func NewSession(conn net.EndPoint, serviceID, objectID, actionID uint32) session.Session {

	sess0 := directorySession{conn, 0, 0, 8}
	service0, err := stage2.NewServer(sess0, 0)
	if err != nil {
		log.Fatalf("failed to create proxy: %s", err)
	}
	permissions := map[string]value.Value{
		"ClientServerSocket":    value.Bool(true),
		"MessageFlags":          value.Bool(true),
		"MetaObjectCache":       value.Bool(true),
		"RemoteCancelableCalls": value.Bool(true),
	}
	_, err = service0.Authenticate(permissions)
	if err != nil {
		log.Fatalf("failed to authenticate: %s", err)
	}
	return directorySession{
		conn,
		serviceID,
		objectID,
		actionID,
	}
}

func NewObject(addr string, serviceID, objectID, actionID uint32) (d *stage2.Object, err error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return d, fmt.Errorf("failed to connect: %s", err)
	}
	sess := NewSession(endpoint, serviceID, objectID, actionID)
	if err != nil {
		return d, fmt.Errorf("failed to create session: %s", err)
	}

	return stage2.NewObject(sess, 1)
}

func NewServiceDirectory(addr string, serviceID, objectID, actionID uint32) (d *stage2.ServiceDirectory, err error) {

	endpoint, err := net.DialEndPoint(addr)
	if err != nil {
		return d, fmt.Errorf("failed to connect: %s", err)
	}
	sess := NewSession(endpoint, serviceID, objectID, actionID)
	if err != nil {
		return d, fmt.Errorf("failed to create session: %s", err)
	}

	return stage2.NewServiceDirectory(sess, 1)
}

func main() {
	var output io.Writer

	if len(os.Args) > 1 {
		filename := os.Args[1]

		file, err := os.Create(filename)
		if err != nil {
			log.Fatalf("failed to open %s: %s", filename, err)
			return
		}
		output = file
		defer file.Close()
	} else {
		output = os.Stdout
	}

	addr := ":9559"
	// directoryServiceID := 1
	// directoryObjectID := 1
	dir, err := NewServiceDirectory(addr, 1, 1, 101)
	if err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}

	serviceInfoList, err := dir.Services()
	if err != nil {
		log.Fatalf("failed to list services: %s", err)
	}

	objects := make([]object.MetaObject, 0)
	objects = append(objects, object.MetaService0)
	objects = append(objects, object.ObjectMetaObject)

	for i, s := range serviceInfoList {

		if i > 1 {
			continue
		}

		addr := strings.TrimPrefix(s.Endpoints[0], "tcp://")
		obj, err := NewObject(addr, s.ServiceId, 1, 2)
		if err != nil {
			log.Printf("failed to create servinceof %s: %s", s.Name, err)
			continue
		}
		meta, err := obj.MetaObject(1)
		if err != nil {
			log.Printf("failed to query MetaObject of %s: %s", s.Name, err)
			continue
		}
		meta.Description = s.Name
		objects = append(objects, meta)
	}
	proxy.GenerateProxys(objects, "services", output)
}
