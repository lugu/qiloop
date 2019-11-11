package services

import (
	"fmt"

	"github.com/lugu/qiloop/bus"
	"github.com/lugu/qiloop/bus/net"
	"github.com/lugu/qiloop/bus/util"
)

// NewServer starts a new server which listen to addr and serve
// incomming requests. It uses the given session to register and
// activate the new services.
func NewServer(sess bus.Session, addr string, auth bus.Authenticator) (bus.Server, error) {
	l, err := net.Listen(addr)
	if err != nil {
		return nil, err
	}

	addrs := []string{addr}
	namespace, err := Namespace(sess, addrs)
	if err != nil {
		return nil, err
	}
	return bus.StandAloneServer(l, auth, namespace)
}

// Namespace returns a Namespace which connects to a remote service
// directory.
func Namespace(session bus.Session, endpoints []string) (bus.Namespace, error) {
	services := Services(session)
	directory, err := services.ServiceDirectory(nil)
	if err != nil {
		return nil, fmt.Errorf("connect ServiceDirectory: %s",
			err)
	}
	return &remoteNamespace{
		session:   session,
		Directory: directory,
		EndPoints: endpoints,
		SessionID: "",
		ProcessID: util.ProcessID(),
		MachineID: util.MachineID(),
	}, nil
}

type remoteNamespace struct {
	session   bus.Session
	Directory ServiceDirectoryProxy
	EndPoints []string
	MachineID string
	ProcessID uint32
	SessionID string
}

func (ns *remoteNamespace) serviceInfo(name string) ServiceInfo {
	return ServiceInfo{
		Name:      name,
		ServiceId: 0,
		MachineId: ns.MachineID,
		ProcessId: ns.ProcessID,
		Endpoints: ns.EndPoints,
		SessionId: ns.SessionID,
	}
}

func (ns *remoteNamespace) Reserve(name string) (uint32, error) {
	return ns.Directory.RegisterService(ns.serviceInfo(name))
}

func (ns *remoteNamespace) Remove(serviceID uint32) error {
	return ns.Directory.UnregisterService(serviceID)
}

func (ns *remoteNamespace) Enable(serviceID uint32) error {
	return ns.Directory.ServiceReady(serviceID)
}

func (ns *remoteNamespace) Resolve(name string) (uint32, error) {
	info, err := ns.Directory.Service(name)
	if err != nil {
		return 0, err
	}
	return info.ServiceId, nil
}

func (ns *remoteNamespace) Session(s bus.Server) bus.Session {
	return ns.session
}
